package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"time"
)

type CustomObjectTypesResponse struct {
	Results []CustomObjectType `json:"results"`
	Paging  *Paging            `json:"paging"`
}

type CustomObjectTypeSchema struct {
	Name                   string                       `json:"name"`
	Labels                 CustomObjectTypeSchemaLabels `json:"labels"`
	PrimaryDisplayProperty string                       `json:"primaryDisplayProperty"`
	RequiredProperties     *[]string                    `json:"requiredProperties,omitempty"`
	Properties             *[]Property                  `json:"properties,omitempty"`
	AssociatedObjects      *[]string                    `json:"associatedObjects,omitempty"`
	MetaType               *string                      `json:"metaType,omitempty"`
}

type CustomObjectTypeSchemaLabels struct {
	Singular string `json:"singular"`
	Plural   string `json:"plural"`
}

type CustomObjectType struct {
	Id           string      `json:"id"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	Properties   *[]Property `json:"properties"`
	Associations []struct {
		Id               string `json:"id"`
		FromObjectTypeId string `json:"fromObjectTypeId"`
		ToObjectTypeId   string `json:"toObjectTypeId"`
		Name             string `json:"name"`
	} `json:"associations"`
	Labels                 CustomObjectTypeSchemaLabels `json:"labels"`
	RequiredProperties     []string                     `json:"requiredProperties"`
	SearchableProperties   []string                     `json:"searchableProperties"`
	PrimaryDisplayProperty string                       `json:"primaryDisplayProperty"`
	MetaType               string                       `json:"metaType"`
	FullyQualifiedName     string                       `json:"fullyQualifiedName"`
	Name                   string                       `json:"name"`
}

func (service *Service) GetCustomObjectTypes() (*[]CustomObjectType, *errortools.Error) {
	var response CustomObjectTypesResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm("schemas"),
		ResponseModel: &response,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Results, nil
}

func (service *Service) CreateCustomObjectType(schema *CustomObjectTypeSchema) (*CustomObjectType, *errortools.Error) {
	var customObjectType CustomObjectType

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm("schemas"),
		BodyModel:     schema,
		ResponseModel: &customObjectType,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObjectType, nil
}

func (service *Service) UpdateCustomObjectType(objectTypeId string, schema *CustomObjectTypeSchema) (*CustomObjectType, *errortools.Error) {
	var customObjectType CustomObjectType

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("schemas/%s", objectTypeId)),
		BodyModel:     schema,
		ResponseModel: &customObjectType,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObjectType, nil
}
