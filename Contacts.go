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

type ContactsResponse struct {
	Results []Contact `json:"results"`
	Paging  *Paging   `json:"paging"`
}

// Contact stores Contact from Service
type Contact struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeMSString   `json:"createdAt"`
	UpdatedAt    h_types.DateTimeMSString   `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type GetContactsConfig struct {
	Limit        *uint
	After        *string
	Properties   *[]string
	Associations *[]string
	Archived     *bool
}

// GetContacts returns all contacts
func (service *Service) GetContacts(config *GetContactsConfig) (*[]Contact, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/contacts"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				values.Set("properties", strings.Join(*config.Properties, ","))
			}
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

	contacts := []Contact{}

	for {
		contactsResponse := ContactsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &contactsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, contactsResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if contactsResponse.Paging == nil {
			break
		}

		if contactsResponse.Paging.Next.After == "" {
			break
		}

		after = contactsResponse.Paging.Next.After
	}

	return &contacts, nil
}

type CreateContactConfig struct {
	Properties map[string]string
}

func (service *Service) CreateContact(config *CreateContactConfig) (*Contact, *errortools.Error) {
	endpoint := "objects/contacts"
	contact := Contact{}

	properties := struct {
		Properties map[string]string `json:"properties"`
	}{
		config.Properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     properties,
		ResponseModel: &contact,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type UpdateContactConfig struct {
	ContactId  string
	Properties map[string]string
}

func (service *Service) UpdateContact(config *UpdateContactConfig) (*Contact, *errortools.Error) {
	endpoint := "objects/contacts"
	contact := Contact{}

	properties := struct {
		Properties map[string]string `json:"properties"`
	}{
		config.Properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.ContactId)),
		BodyModel:     properties,
		ResponseModel: &contact,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type GetContactConfig struct {
	ContactId    string
	Properties   *[]string
	Associations *[]string
}

// GetContact returns a specific contact
func (service *Service) GetContact(config *GetContactConfig) (*Contact, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/contacts"

	if config == nil {
		return nil, errortools.ErrorMessage("config is nil")
	}

	_properties := []string{}
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
			_associations := []string{}
			for _, a := range *config.Associations {
				_associations = append(_associations, string(a))
			}
			values.Set("associations", strings.Join(_associations, ","))
		}
	}

	contact := Contact{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s?%s", endpoint, config.ContactId, values.Encode())),
		ResponseModel: &contact,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type SearchContactConfig struct {
	Limit        *uint          `json:"limit,omitempty"`
	After        *string        `json:"after,omitempty"`
	FilterGroups *[]FilterGroup `json:"filterGroups,omitempty"`
	Sorts        *[]string      `json:"sorts,omitempty"`
	Query        *string        `json:"query,omitempty"`
	Properties   *[]string      `json:"properties,omitempty"`
}

// SearchContact returns a specific contact
func (service *Service) SearchContact(config *SearchContactConfig) (*[]Contact, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config is nil")
	}

	endpoint := "objects/contacts/search"

	contactsResponse := ContactsResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &contactsResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	customProperties := []string{}
	if config.FilterGroups != nil {
		for _, filterGroup := range *config.FilterGroups {
			for _, filter := range *filterGroup.Filters {
				if filter.isCustom {
					customProperties = append(customProperties, filter.PropertyName)
				}
			}
		}
	}

	after := config.After

	contacts := []Contact{}

	for {
		contactsResponse := ContactsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm("objects/contacts/search"),
			BodyModel:     config,
			ResponseModel: &contactsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, contactsResponse.Results...)

		if after != nil { // explicit after parameter requested
			break
		}

		if contactsResponse.Paging == nil {
			break
		}

		if contactsResponse.Paging.Next.After == "" {
			break
		}

		config.After = &contactsResponse.Paging.Next.After
	}

	return &contacts, nil
}

func (service *Service) DeleteContact(contactId string) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.urlCrm(fmt.Sprintf("objects/contacts/%s", contactId)),
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}

func (service *Service) BatchDeleteContacts(contactIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(contactIds) > index {
		if len(contactIds) > index+maxItemsPerBatch {
			e := service.batchDeleteContacts(contactIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteContacts(contactIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteContacts(contactIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, contactId := range contactIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{contactId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/contacts/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
