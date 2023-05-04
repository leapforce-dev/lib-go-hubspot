package hubspot

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"time"
)

type EmailsResponse struct {
	Results []Email `json:"results"`
	Paging  *Paging `json:"paging"`
}

type Email struct {
	Id           string             `json:"id,omitempty"`
	Properties   EmailProperties    `json:"properties"`
	CreatedAt    *time.Time         `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time         `json:"updatedAt,omitempty"`
	Archived     *bool              `json:"archived,omitempty"`
	Associations *EmailAssociations `json:"associations,omitempty"`
}

type EmailUpdate struct {
	Properties   EmailProperties     `json:"properties"`
	Associations *[]EmailAssociation `json:"associations,omitempty"`
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
	OwnerId         string    `json:"hubspot_owner_id"`
	AttachmentIds   string    `json:"hs_attachment_ids"`
	Direction       string    `json:"hs_email_direction"`
	Headers         string    `json:"hs_email_headers,omitempty"`
	SenderEmail     string    `json:"hs_email_sender_email,omitempty"`
	SenderFirstname string    `json:"hs_email_sender_firstname,omitempty"`
	SenderLastname  string    `json:"hs_email_sender_lastname,omitempty"`
	ToEmail         string    `json:"hs_email_to_email,omitempty"`
	ToFirstname     string    `json:"hs_email_to_firstname,omitempty"`
	ToLastname      string    `json:"hs_email_to_lastname,omitempty"`
	Status          string    `json:"hs_email_status"`
	Subject         string    `json:"hs_email_subject"`
	Text            string    `json:"hs_email_text"`
	Timestamp       time.Time `json:"hs_timestamp"`
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
		e.Headers = ""
		return nil
	}

	b, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	e.Headers = string(b)

	return nil
}

func (service *Service) CreateEmail(email *EmailUpdate) (*Email, *errortools.Error) {
	var newEmail Email

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlV4("objects/emails"),
		BodyModel:     email,
		ResponseModel: &newEmail,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &newEmail, nil
}

type ListEmailsConfig struct {
	Limit *uint
	After *string
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

		emails = append(emails, emailsResponse.Results...)

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
