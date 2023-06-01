package hubspot

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"time"
)

type CustomObjectSchema struct {
	Name   string `json:"name"`
	Labels struct {
		Singular string `json:"singular"`
		Plural   string `json:"plural"`
	} `json:"labels"`
	PrimaryDisplayProperty string     `json:"primaryDisplayProperty"`
	RequiredProperties     []string   `json:"requiredProperties"`
	Properties             []Property `json:"properties"`
	AssociatedObjects      []string   `json:"associatedObjects"`
	MetaType               string     `json:"metaType"`
}

type CustomObject struct {
	Id           string     `json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	Properties   []Property `json:"properties"`
	Associations []struct {
		Id               string `json:"id"`
		FromObjectTypeId string `json:"fromObjectTypeId"`
		ToObjectTypeId   string `json:"toObjectTypeId"`
		Name             string `json:"name"`
	} `json:"associations"`
	Labels struct {
		Singular string `json:"singular"`
		Plural   string `json:"plural"`
	} `json:"labels"`
	RequiredProperties     []string `json:"requiredProperties"`
	SearchableProperties   []string `json:"searchableProperties"`
	PrimaryDisplayProperty string   `json:"primaryDisplayProperty"`
	MetaType               string   `json:"metaType"`
	FullyQualifiedName     string   `json:"fullyQualifiedName"`
	Name                   string   `json:"name"`
}

func (service *Service) CreateCustomObject(schema *CustomObjectSchema) (*CustomObject, *errortools.Error) {
	var customObject CustomObject

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm("schemas"),
		BodyModel:     schema,
		ResponseModel: &customObject,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &customObject, nil
}
