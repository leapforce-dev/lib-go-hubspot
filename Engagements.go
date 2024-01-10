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

type EngagementType string

const (
	EngagementTypeCall          EngagementType = "call"
	EngagementTypeCommunication EngagementType = "communication"
	EngagementTypeEmail         EngagementType = "email"
	EngagementTypeMeeting       EngagementType = "meeting"
	EngagementTypeNote          EngagementType = "note"
	EngagementTypePostalMail    EngagementType = "postal_mail"
	EngagementTypeTask          EngagementType = "task"
)

type EngagementsResponse struct {
	Results []Engagement `json:"results"`
	Paging  *Paging      `json:"paging"`
}

// Engagement stores Engagement from Service
type Engagement struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type ListEngagementsConfig struct {
	Type         EngagementType
	Limit        *uint
	After        *string
	Properties   *[]string
	Associations *[]string
	Archived     *bool
}

// ListEngagements returns all engagements
func (service *Service) ListEngagements(config *ListEngagementsConfig) (*[]Engagement, *errortools.Error) {
	values := url.Values{}
	endpoint := fmt.Sprintf("objects/%v", config.Type)

	after := ""

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		var _properties []string
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				_properties = append(_properties, *config.Properties...)
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

	var engagements []Engagement

	for {
		engagementsResponse := EngagementsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &engagementsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		engagements = append(engagements, engagementsResponse.Results...)

		if config != nil {
			if config.After != nil { // explicit after parameter requested
				break
			}
		}

		if engagementsResponse.Paging == nil {
			break
		}

		if engagementsResponse.Paging.Next.After == "" {
			break
		}

		after = engagementsResponse.Paging.Next.After
	}

	return &engagements, nil
}

type CreateEngagementConfig struct {
	Type         EngagementType
	Properties   map[string]string `json:"properties"`
	Associations *[]AssociationTo  `json:"associations,omitempty"`
}

func (service *Service) CreateEngagement(config *CreateEngagementConfig) (*Engagement, *errortools.Error) {
	endpoint := fmt.Sprintf("objects/%v", config.Type)
	engagement := Engagement{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &engagement,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &engagement, nil
}

type UpdateEngagementConfig struct {
	Type         EngagementType
	EngagementId string
	Properties   map[string]string
}

func (service *Service) UpdateEngagement(config *UpdateEngagementConfig) (*Engagement, *errortools.Error) {
	endpoint := fmt.Sprintf("objects/%v", config.Type)
	engagement := Engagement{}

	var properties = make(map[string]string)

	if config.Properties != nil {
		for key, value := range config.Properties {
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
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.EngagementId)),
		BodyModel:     properties_,
		ResponseModel: &engagement,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &engagement, nil
}

func (service *Service) BatchDeleteEngagements(engagementType EngagementType, engagementIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(engagementIds) > index {
		if len(engagementIds) > index+maxItemsPerBatch {
			e := service.batchDeleteEngagements(engagementType, engagementIds[index:index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteEngagements(engagementType, engagementIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteEngagements(engagementType EngagementType, engagementIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, engagementId := range engagementIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{engagementId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm(fmt.Sprintf("objects/%v/batch/archive", engagementType)),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
