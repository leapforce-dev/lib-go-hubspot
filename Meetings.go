package hubspot

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
)

type MeetingsResponse struct {
	Results []Meeting `json:"results"`
	Paging  *Paging   `json:"paging"`
}

// Meeting stores Meeting from Service
type Meeting struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type GetMeetingsConfig struct {
	Limit            *uint
	After            *string
	CustomProperties *[]string
	Associations     *[]string
	Archived         *bool
}

// GetMeetings returns all meetings
func (service *Service) GetMeetings(config *GetMeetingsConfig) (*[]Meeting, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/meetings"

	after := ""

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		var _properties []string
		if config.CustomProperties != nil {
			if len(*config.CustomProperties) > 0 {
				_properties = append(_properties, *config.CustomProperties...)
			}
		}
		if len(_properties) > 0 {
			values.Set("properties", strings.Join(_properties, ","))
		}
		if config.Associations != nil {
			if len(*config.Associations) > 0 {
				var _associations []string
				for _, a := range *config.Associations {
					_associations = append(_associations, a)
				}
				values.Set("associations", strings.Join(_associations, ","))
			}
		}
		if config.Archived != nil {
			values.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}

		if config.After != nil {
			after = *config.After
		}
	}

	var meetings []Meeting

	for {
		meetingsResponse := MeetingsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &meetingsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		meetings = append(meetings, meetingsResponse.Results...)

		if config != nil {
			if config.After != nil { // explicit after parameter requested
				break
			}
		}

		if meetingsResponse.Paging == nil {
			break
		}

		if meetingsResponse.Paging.Next.After == "" {
			break
		}

		after = meetingsResponse.Paging.Next.After
	}

	return &meetings, nil
}

type CreateMeetingConfig struct {
	Properties   map[string]string `json:"properties"`
	Associations *[]AssociationTo  `json:"associations,omitempty"`
}

func (service *Service) CreateMeeting(config *CreateMeetingConfig) (*Meeting, *errortools.Error) {
	endpoint := "objects/meetings"
	meeting := Meeting{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &meeting,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &meeting, nil
}

type UpdateMeetingConfig struct {
	MeetingId        string
	CustomProperties map[string]string
}

func (service *Service) UpdateMeeting(config *UpdateMeetingConfig) (*Meeting, *errortools.Error) {
	endpoint := "objects/meetings"
	meeting := Meeting{}

	var properties = make(map[string]string)

	if config.CustomProperties != nil {
		for key, value := range config.CustomProperties {
			properties[key] = value
		}
	}

	var properties_ = struct {
		Properties map[string]string `json:"properties"`
	}{
		properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.MeetingId)),
		BodyModel:     properties_,
		ResponseModel: &meeting,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &meeting, nil
}

func (service *Service) BatchDeleteMeetings(meetingIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(meetingIds) > index {
		if len(meetingIds) > index+maxItemsPerBatch {
			e := service.batchDeleteMeetings(meetingIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteMeetings(meetingIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteMeetings(meetingIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, meetingId := range meetingIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{meetingId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/meetings/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
