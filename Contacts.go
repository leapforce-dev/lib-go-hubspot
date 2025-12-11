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

type ContactsResponse struct {
	Results []Contact `json:"results"`
	Paging  *Paging   `json:"paging"`
}

// Contact stores Contact from Service
type Contact struct {
	Id                    string                       `json:"id"`
	Properties            map[string]string            `json:"properties"`
	CreatedAt             h_types.DateTimeMSString     `json:"createdAt"`
	UpdatedAt             h_types.DateTimeMSString     `json:"updatedAt"`
	Archived              bool                         `json:"archived"`
	ArchivedAt            h_types.DateTimeMSString     `json:"archivedAt"`
	Associations          map[string]AssociationsSet   `json:"associations"`
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory"`
}

type GetContactsConfig struct {
	Limit                 *uint
	After                 *string
	Properties            *[]string
	PropertiesWithHistory *[]string
	Associations          *[]string
	Archived              *bool
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
		if config.PropertiesWithHistory != nil {
			if len(*config.PropertiesWithHistory) > 0 {
				values.Set("propertiesWithHistory", strings.Join(*config.PropertiesWithHistory, ","))
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

func (service *Service) CreateContact(config *CreateObjectConfig) (*Contact, *errortools.Error) {
	endpoint := "objects/contacts"
	contact := Contact{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &contact,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

func (service *Service) BatchCreateContacts(config *BatchObjectsConfig, invalidEmailProperty string) (*[]Contact, *errortools.Error) {
	var contacts []Contact

	for _, batch := range service.batches(len(config.Inputs)) {
		var retrying = false
	retry:
		var r BatchContactsResponse

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

			if response.StatusCode == http.StatusBadRequest && !retrying {
				errorResponse := service.ErrorResponse()
				if errorResponse != nil {
					if strings.HasPrefix(errorResponse.Message, "Property values were not valid: ") {
						stop := checkInvalidEmails(config, invalidEmailProperty, errorResponse.Message, batch)

						if stop {
							goto stop
						}

						retrying = true
						goto retry
					}
				}
			}
		}

	stop:
		if e != nil {
			return nil, e
		}
	ok:
		contacts = append(contacts, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &contacts, nil
}

func checkInvalidEmails(config *BatchObjectsConfig, invalidEmailProperty string, message string, batch batch) (stop bool) {
	m := strings.TrimPrefix(message, "Property values were not valid: ")

	var propertyErrors []PropertyError
	err := json.Unmarshal([]byte(m), &propertyErrors)
	if err != nil {
		fmt.Println(err)
		return true
	}

	hasNonEmailError := false
	for _, propertyError := range propertyErrors {
		fmt.Println(propertyError.Message)

		if propertyError.Error == "INVALID_EMAIL" {
			// English error message
			email := strings.TrimSuffix(strings.TrimPrefix(propertyError.Message, "Email address "), " is invalid")
			// Dutch error message
			email = strings.TrimSuffix(strings.TrimPrefix(propertyError.Message, "E-mailadres "), " is ongeldig")
			for i := batch.startIndex; i < batch.endIndex; i++ {
				if (*config).Inputs[i].Properties["email"] == email {
					(*config).Inputs[i].Properties["email"] = ""
					if invalidEmailProperty != "" {
						(*config).Inputs[i].Properties[invalidEmailProperty] = email
					}
				}
			}
		} else {
			hasNonEmailError = true
		}
	}
	return hasNonEmailError
}

type BatchContactsResponse struct {
	CompletedAt *time.Time        `json:"completedAt"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt *time.Time        `json:"requestedAt"`
	StartedAt   *time.Time        `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []Contact         `json:"results"`
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

func (service *Service) BatchUpdateContacts(config *BatchObjectsConfig, invalidEmailProperty string) (*[]Contact, *errortools.Error) {
	var contacts []Contact

	for _, batch := range service.batches(len(config.Inputs)) {
	retry:
		var r BatchContactsResponse

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
			if response.StatusCode == http.StatusBadRequest {
				errorResponse := service.ErrorResponse()
				if errorResponse != nil {
					if strings.HasPrefix(errorResponse.Message, "Property values were not valid: ") {
						stop := checkInvalidEmails(config, invalidEmailProperty, errorResponse.Message, batch)

						if stop {
							goto stop
						}

						goto retry
					}
				}
			}
		}
	stop:
		if e != nil {
			return nil, e
		}
	ok:
		contacts = append(contacts, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &contacts, nil
}

func (service *Service) UpdateContact(config *UpdateObjectConfig) (*Contact, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/contacts"

	if config == nil {
		return nil, errortools.ErrorMessage("config is nil")
	}

	contact := Contact{}

	if config.IdProperty != nil {
		values.Set("idProperty", *config.IdProperty)
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s?%s", endpoint, config.ObjectId, values.Encode())),
		BodyModel:     config,
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
	IdProperty   *string
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

	if config.IdProperty != nil {
		values.Set("idProperty", *config.IdProperty)
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

	_, response, e := service.httpRequest(&requestConfig)
	if response != nil {
		if response.StatusCode == http.StatusNotFound {
			return nil, nil
		}
	}
	if e != nil {
		return nil, e
	}

	return &contact, nil
}

type SearchObjectsConfig struct {
	Limit        *uint          `json:"limit,omitempty"`
	After        *string        `json:"after,omitempty"`
	FilterGroups *[]FilterGroup `json:"filterGroups,omitempty"`
	Sorts        *[]string      `json:"sorts,omitempty"`
	Query        *string        `json:"query,omitempty"`
	Properties   *[]string      `json:"properties,omitempty"`
}

// SearchContact returns a specific contact
func (service *Service) SearchContact(config *SearchObjectsConfig) (*[]Contact, *errortools.Error) {
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

	properties := []string{}
	if config.FilterGroups != nil {
		for _, filterGroup := range *config.FilterGroups {
			for _, filter := range *filterGroup.Filters {
				if filter.isCustom {
					properties = append(properties, filter.PropertyName)
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

func (service *Service) BatchArchiveContacts(contactIds []string) *errortools.Error {
	var index = 0
	for len(contactIds) > index {
		if len(contactIds) > index+maxItemsPerBatch {
			e := service.batchArchiveContacts(contactIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchArchiveContacts(contactIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchArchiveContacts(contactIds []string) *errortools.Error {
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
