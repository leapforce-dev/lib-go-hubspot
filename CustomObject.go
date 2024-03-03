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

type CustomObjectsResponse struct {
	Results []CustomObject `json:"results"`
	Paging  *Paging        `json:"paging"`
}

type CustomObject struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type GetCustomObjectsConfig struct {
	ObjectType            string
	Limit                 *uint
	After                 *string
	Properties            *[]string
	PropertiesWithHistory *[]string
	Associations          *[]string
	Archived              *bool
}

// GetCustomObjects returns all customObjects
func (service *Service) GetCustomObjects(config *GetCustomObjectsConfig) (*[]CustomObject, *errortools.Error) {
	values := url.Values{}
	endpoint := fmt.Sprintf("objects/%s", config.ObjectType)

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		_properties := []string{}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				_properties = append(_properties, *config.Properties...)
			}
		}
		if config.PropertiesWithHistory != nil {
			if len(*config.PropertiesWithHistory) > 0 {
				values.Set("propertiesWithHistory", strings.Join(*config.PropertiesWithHistory, ","))
			}
		}
		if len(_properties) > 0 {
			values.Set("properties", strings.Join(_properties, ","))
		}
		if config.Associations != nil {
			if len(*config.Associations) > 0 {
				_associations := []string{}
				for _, a := range *config.Associations {
					_associations = append(_associations, string(a))
				}
				values.Set("associations", strings.Join(_associations, ","))
			}
		}
		if config.Archived != nil {
			values.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}
	}

	after := ""
	if config.After != nil {
		after = *config.After
	}

	customObjects := []CustomObject{}

	for {
		customObjectsResponse := CustomObjectsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &customObjectsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		customObjects = append(customObjects, customObjectsResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if customObjectsResponse.Paging == nil {
			break
		}

		if customObjectsResponse.Paging.Next.After == "" {
			break
		}

		after = customObjectsResponse.Paging.Next.After
	}

	return &customObjects, nil
}

func (service *Service) CreateCustomObject(config *CreateObjectConfig) (*CustomObject, *errortools.Error) {
	endpoint := fmt.Sprintf("objects/%s", config.ObjectType)
	customObject := CustomObject{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &customObject,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObject, nil
}

func (service *Service) UpdateCustomObject(config *UpdateObjectConfig) (*CustomObject, *errortools.Error) {
	endpoint := fmt.Sprintf("objects/%s", config.ObjectType)

	customObject := CustomObject{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.ObjectId)),
		BodyModel:     config,
		ResponseModel: &customObject,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObject, nil
}

func (service *Service) BatchDeleteCustomObjects(objectType string, customObjectIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(customObjectIds) > index {
		if len(customObjectIds) > index+maxItemsPerBatch {
			e := service.batchDeleteCustomObjects(objectType, customObjectIds[index:index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteCustomObjects(objectType, customObjectIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteCustomObjects(objectType string, customObjectIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, customObjectId := range customObjectIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{customObjectId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm(fmt.Sprintf("objects/%s/batch/archive", objectType)),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
