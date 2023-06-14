package hubspot

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type EmailsResponse struct {
	Results []email `json:"results"`
	Paging  *Paging `json:"paging"`
}

type email struct {
	Id           string                     `json:"id"`
	Properties   json.RawMessage            `json:"properties"`
	CreatedAt    h_types.DateTimeMSString   `json:"createdAt"`
	UpdatedAt    h_types.DateTimeMSString   `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type Email struct {
	Id               string
	Properties       EmailProperties
	CustomProperties map[string]string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Archived         bool
	Associations     map[string]AssociationsSet
}

type EmailAssociations struct {
	Contacts AssociationsSet `json:"contacts"`
}

type EmailAssociation struct {
	To struct {
		Id string `json:"id"`
	} `json:"to"`
	Types []AssociationTypeV4 `json:"types"`
}

func NewEmailAssociation(toId string, category string, typeId int64) EmailAssociation {
	return EmailAssociation{
		To: struct {
			Id string `json:"id"`
		}{toId},
		Types: []AssociationTypeV4{
			{
				AssociationCategory: category,
				AssociationTypeId:   typeId,
			},
		},
	}
}

type EmailProperties struct {
	OwnerId         *string    `json:"hubspot_owner_id"`
	AttachmentIds   *string    `json:"hs_attachment_ids"`
	Direction       *string    `json:"hs_email_direction"`
	Headers         *string    `json:"hs_email_headers,omitempty"`
	SenderEmail     *string    `json:"hs_email_sender_email,omitempty"`
	SenderFirstname *string    `json:"hs_email_sender_firstname,omitempty"`
	SenderLastname  *string    `json:"hs_email_sender_lastname,omitempty"`
	ToEmail         *string    `json:"hs_email_to_email,omitempty"`
	ToFirstname     *string    `json:"hs_email_to_firstname,omitempty"`
	ToLastname      *string    `json:"hs_email_to_lastname,omitempty"`
	Status          *string    `json:"hs_email_status"`
	Subject         *string    `json:"hs_email_subject"`
	Text            *string    `json:"hs_email_text"`
	Timestamp       *time.Time `json:"hs_timestamp"`
}

type EmailHeaderItem struct {
	Email     string  `json:"email"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
}

type EmailHeaders struct {
	From EmailHeaderItem    `json:"from"`
	To   *[]EmailHeaderItem `json:"to,omitempty"`
	Cc   *[]EmailHeaderItem `json:"cc,omitempty"`
	Bcc  *[]EmailHeaderItem `json:"bcc,omitempty"`
}

func (e *EmailProperties) SetHeaders(headers *EmailHeaders) error {
	if headers == nil {
		h := ""
		e.Headers = &h
		return nil
	}

	b, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	h := string(b)

	e.Headers = &h

	return nil
}

type ListEmailsConfig struct {
	Limit        *uint
	After        *string
	Properties   *[]string
	Associations *[]string
	Archived     *bool
}

func (service *Service) ListEmails(config *ListEmailsConfig) (*[]Email, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/emails"

	after := ""

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		if config.After != nil {
			after = *config.After
		}
		if config.Properties != nil {
			values.Set("properties", strings.Join(*config.Properties, ","))
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

	var emails []Email

	for {
		if after != "" {
			values.Set("after", after)
		}

		var emailsResponse EmailsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlV4(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &emailsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, em := range emailsResponse.Results {
			email_, e := getEmail(&em, config.Properties)
			if e != nil {
				return nil, e
			}
			emails = append(emails, *email_)
		}

		if config != nil {
			if config.After != nil { // explicit after parameter requested
				break
			}
		}

		if emailsResponse.Paging == nil {
			break
		}

		if emailsResponse.Paging.Next.After == "" {
			break
		}

		after = emailsResponse.Paging.Next.After
	}

	return &emails, nil
}

func getEmail(email *email, customProperties *[]string) (*Email, *errortools.Error) {
	email_ := Email{
		Id:               email.Id,
		CreatedAt:        email.CreatedAt.Value(),
		UpdatedAt:        email.UpdatedAt.Value(),
		Archived:         email.Archived,
		Associations:     email.Associations,
		CustomProperties: make(map[string]string),
	}
	if email.Properties != nil {
		p := EmailProperties{}
		err := json.Unmarshal(email.Properties, &p)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
		email_.Properties = p
	}

	if customProperties != nil {
		p1 := make(map[string]string)
		err := json.Unmarshal(email.Properties, &p1)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}

		for _, cp := range *customProperties {
			value, ok := p1[cp]
			if ok {
				email_.CustomProperties[cp] = value
			}
		}
	}

	return &email_, nil
}

func (service *Service) BatchDeleteEmails(emailIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(emailIds) > index {
		if len(emailIds) > index+maxItemsPerBatch {
			e := service.batchDeleteEmails(emailIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteEmails(emailIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteEmails(emailIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, emailId := range emailIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{emailId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/emails/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}

type CreateEmailConfig struct {
	Properties       EmailProperties
	CustomProperties map[string]string
}

func (service *Service) CreateEmail(config *CreateEmailConfig) (*Email, *errortools.Error) {
	endpoint := "objects/emails"
	email := Email{}

	body, e := emailPropertiesBody(config.Properties, config.CustomProperties)
	if e != nil {
		return nil, e
	}

	properties := struct {
		Properties map[string]string `json:"properties"`
	}{
		body,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlV4(endpoint),
		BodyModel:     properties,
		ResponseModel: &email,
	}

	_, _, e = service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &email, nil
}

func emailPropertiesBody(properties EmailProperties, customProperties map[string]string) (map[string]string, *errortools.Error) {
	// marshal
	b, err := json.Marshal(properties)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	// unmarshal to map
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	if customProperties == nil {
		return m, nil
	}
	if len(customProperties) == 0 {
		return m, nil
	}

	// append custom properties to map
	for key, value := range customProperties {
		if _, ok := m[key]; !ok {
			m[key] = value
		}
	}

	return m, nil
}
