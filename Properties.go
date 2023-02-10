package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type PropertyType string

const (
	PropertyTypeString      PropertyType = "string"
	PropertyTypeNumber      PropertyType = "number"
	PropertyTypeDate        PropertyType = "date"
	PropertyTypeDateTime    PropertyType = "datetime"
	PropertyTypeBool        PropertyType = "bool"
	PropertyTypeEnumeration PropertyType = "enumeration"
)

type PropertyFieldType string

const (
	PropertyFieldTypeTextarea        PropertyFieldType = "textarea"
	PropertyFieldTypeText            PropertyFieldType = "text"
	PropertyFieldTypeDate            PropertyFieldType = "date"
	PropertyFieldTypeFile            PropertyFieldType = "file"
	PropertyFieldTypeNumber          PropertyFieldType = "number"
	PropertyFieldTypeSelect          PropertyFieldType = "select"
	PropertyFieldTypeRadio           PropertyFieldType = "radio"
	PropertyFieldTypeCheckbox        PropertyFieldType = "checkbox"
	PropertyFieldTypeBooleanCheckbox PropertyFieldType = "booleancheckbox"
)

type PropertiesResponse struct {
	Results []Property `json:"results"`
}

// Property stores Property from Service
//
type Property struct {
	Name                 *string            `json:"name"`
	Label                *string            `json:"label"`
	Type                 *PropertyType      `json:"type"`
	FieldType            *PropertyFieldType `json:"fieldType"`
	Description          *string            `json:"description,omitempty"`
	GroupName            *string            `json:"groupName"`
	ReferencedObjectType *string            `json:"referencedObjectType,omitempty"`
	DisplayOrder         *int64             `json:"displayOrder,omitempty"`
	Calculated           *bool              `json:"calculated,omitempty"`
	ExternalOptions      *bool              `json:"externalOption,omitemptys"`
	HasUniqueValue       *bool              `json:"hasUniqueValue,omitempty"`
	Hidden               *bool              `json:"hidden,omitempty"`
	HubspotDefined       *bool              `json:"hubspotDefined,omitempty"`
	FormField            *bool              `json:"formField,omitempty"`
	Options              *[]PropertyOption  `json:"options,omitempty"`
}

type PropertyOption struct {
	Label        string  `json:"label"`
	Value        string  `json:"value"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int64  `json:"displayOrder,omitempty"`
	Hidden       *bool   `json:"hidden,omitempty"`
}

// GetProperties returns all properties
//
func (service *Service) GetProperties(object string) (*[]Property, *errortools.Error) {
	endpoint := "properties"
	propertiesResponse := PropertiesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, object)),
		ResponseModel: &propertiesResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &propertiesResponse.Results, nil
}

// CreateProperty creates a property
//
func (service *Service) CreateProperty(object string, property *Property) (*Property, *errortools.Error) {
	endpoint := "properties"
	newProperty := Property{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, object)),
		BodyModel:     property,
		ResponseModel: &newProperty,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &newProperty, nil
}

// UpdateProperty creates a property
//
func (service *Service) UpdateProperty(object string, property *Property) (*Property, *errortools.Error) {
	endpoint := "properties"
	updatedProperty := Property{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s/%s", endpoint, object, *property.Name)),
		BodyModel:     property,
		ResponseModel: &updatedProperty,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &updatedProperty, nil
}

type PropertyGroupsResponse struct {
	Results []PropertyGroup `json:"results"`
}

type PropertyGroup struct {
	Name         string `json:"name"`
	Label        string `json:"label"`
	DisplayOrder *int64 `json:"displayOrder,omitempty"`
}

// GetPropertyGroups returns all property groups
//
func (service *Service) GetPropertyGroups(object string) (*[]PropertyGroup, *errortools.Error) {
	propertyGroupsResponse := PropertyGroupsResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm(fmt.Sprintf("properties/%s/groups", object)),
		ResponseModel: &propertyGroupsResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &propertyGroupsResponse.Results, nil
}

// CreatePropertyGroup creates a property group
//
func (service *Service) CreatePropertyGroup(object string, propertyGroup *PropertyGroup) (*PropertyGroup, *errortools.Error) {
	newPropertyGroup := PropertyGroup{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(fmt.Sprintf("properties/%s/groups", object)),
		BodyModel:     propertyGroup,
		ResponseModel: &newPropertyGroup,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &newPropertyGroup, nil
}
