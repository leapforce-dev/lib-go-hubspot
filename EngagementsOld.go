package hubspot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
)

type EngagementOldsResponse struct {
	Results []EngagementOld `json:"results"`
	Paging  *Paging         `json:"paging"`
}

// EngagementOld stores EngagementOld from Service
type EngagementOld struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type ListEngagementOldsConfig struct {
	Type         EngagementType
	Limit        *uint
	After        *string
	Properties   *[]string
	Associations *[]string
	Archived     *bool
}

// ListEngagementOlds returns all engagements
func (service *Service) ListEngagementOlds(config *ListEngagementOldsConfig) (*[]EngagementOld, *errortools.Error) {
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

	var engagements []EngagementOld

	for {
		engagementsResponse := EngagementOldsResponse{}

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

type CreateEngagementOldConfig struct {
	Type         EngagementType
	Properties   map[string]string  `json:"properties"`
	Associations *[]AssociationToV4 `json:"associations,omitempty"`
}

func (service *Service) CreateEngagementOld(config *CreateEngagementOldConfig) (*EngagementOld, *errortools.Error) {
	endpoint := fmt.Sprintf("objects/%v", config.Type)
	engagement := EngagementOld{}

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

type UpdateEngagementOldConfig struct {
	Type            EngagementType
	EngagementOldId string
	Properties      map[string]string
}

func (service *Service) UpdateEngagementOld(config *UpdateEngagementOldConfig) (*EngagementOld, *errortools.Error) {
	endpoint := fmt.Sprintf("objects/%v", config.Type)
	engagement := EngagementOld{}

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
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.EngagementOldId)),
		BodyModel:     properties_,
		ResponseModel: &engagement,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &engagement, nil
}

func (service *Service) BatchArchiveEngagementOlds(engagementType EngagementType, engagementIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(engagementIds) > index {
		if len(engagementIds) > index+maxItemsPerBatch {
			e := service.batchArchiveEngagementOlds(engagementType, engagementIds[index:index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchArchiveEngagementOlds(engagementType, engagementIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchArchiveEngagementOlds(engagementType EngagementType, engagementIds []string) *errortools.Error {
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

// SearchEngagementOlds returns a specific engagement
func (service *Service) SearchEngagementOlds(objectType ObjectType, config *SearchObjectsConfig) (*[]EngagementOld, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config is nil")
	}

	endpoint := fmt.Sprintf("objects/%s/search", objectType)

	engagementsResponse := EngagementOldsResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &engagementsResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	after := config.After

	var engagements []EngagementOld

	for {
		engagementsResponse := EngagementOldsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(endpoint),
			BodyModel:     config,
			ResponseModel: &engagementsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		engagements = append(engagements, engagementsResponse.Results...)

		if after != nil { // explicit after parameter requested
			break
		}

		if engagementsResponse.Paging == nil {
			break
		}

		if engagementsResponse.Paging.Next.After == "" {
			break
		}

		config.After = &engagementsResponse.Paging.Next.After
	}

	return &engagements, nil
}

type BatchEngagementOldsResponse struct {
	CompletedAt *time.Time        `json:"completedAt"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt *time.Time        `json:"requestedAt"`
	StartedAt   *time.Time        `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []EngagementOld   `json:"results"`
	Errors      []struct {
		SubCategory json.RawMessage   `json:"subCategory"`
		Context     map[string]string `json:"context"`
		Links       map[string]string `json:"links"`
		Id          string            `json:"id"`
		Category    string            `json:"category"`
		Message     string            `json:"message"`
		Errors      []struct {
			SubCategory string `json:"subCategory"`
			Code        string `json:"code"`
			In          string `json:"in"`
			Context     struct {
				MissingScopes []string `json:"missingScopes"`
			} `json:"context"`
			Message string `json:"message"`
		} `json:"errors"`
		Status string `json:"status"`
	} `json:"errors"`
	Status string `json:"status"`
}

func (service *Service) BatchCreateEngagementOlds(config *BatchObjectsConfig) (*[]EngagementOld, *errortools.Error) {
	var engagements []EngagementOld

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchEngagementOldsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(fmt.Sprintf("objects/%s/batch/create", config.ObjectType)),
			BodyModel:     BatchObjectsConfig{Inputs: config.Inputs[batch.startIndex:batch.endIndex]},
			ResponseModel: &r,
		}

		_, response, e := service.httpRequest(&requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusMultiStatus {
				fmt.Println(r.Errors)
				goto ok
			}
		}
		if e != nil {
			return nil, e
		}
	ok:
		engagements = append(engagements, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &engagements, nil
}

func (service *Service) BatchUpdateEngagementOlds(config *BatchObjectsConfig) (*[]EngagementOld, *errortools.Error) {
	var engagements []EngagementOld

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchEngagementOldsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(fmt.Sprintf("objects/%s/batch/update", config.ObjectType)),
			BodyModel:     BatchObjectsConfig{Inputs: config.Inputs[batch.startIndex:batch.endIndex]},
			ResponseModel: &r,
		}

		_, response, e := service.httpRequest(&requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusMultiStatus {
				fmt.Println(r.Errors)
				goto ok
			}
		}
		if e != nil {
			return nil, e
		}
	ok:
		engagements = append(engagements, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &engagements, nil
}
