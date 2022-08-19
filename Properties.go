package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type PropertiesResponse struct {
	Results []Property `json:"results"`
}

// Property stores Property from Service
//
type Property struct {
	Name                 string `json:"name"`
	Label                string `json:"label"`
	Type                 string `json:"type"`
	FieldType            string `json:"fieldType"`
	Description          string `json:"description"`
	GroupName            string `json:"groupName"`
	ReferencedObjectType string `json:"referencedObjectType"`
	DisplayOrder         int64  `json:"displayOrder"`
	Calculated           *bool  `json:"calculated"`
	ExternalOptions      *bool  `json:"externalOptions"`
	HasUniqueValue       *bool  `json:"hasUniqueValue"`
	Hidden               *bool  `json:"hidden"`
	HubspotDefined       *bool  `json:"hubspotDefined"`
	FormField            *bool  `json:"formField"`
}

// GetProperties returns all properties
//
func (service *Service) GetProperties(object string) (*[]Property, *errortools.Error) {
	endpoint := "properties"
	propertiesResponse := PropertiesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("%s/%s", endpoint, object)),
		ResponseModel: &propertiesResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &propertiesResponse.Results, nil
}
