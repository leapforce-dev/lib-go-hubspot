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
)

type EmailsResponse struct {
	Results []Email `json:"results"`
	Paging  *Paging `json:"paging"`
}

type Email struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeMSString   `json:"createdAt"`
	UpdatedAt    h_types.DateTimeMSString   `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
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

func SetEmailHeaders(properties map[string]string, headers *EmailHeaders) error {
	if headers == nil {
		properties["hs_email_headers"] = ""
		return nil
	}

	b, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	properties["hs_email_headers"] = string(b)

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

type CreateEmailConfig struct {
	Properties map[string]string
}

func (service *Service) CreateEmail(config *CreateEmailConfig) (*Email, *errortools.Error) {
	endpoint := "objects/emails"
	email := Email{}

	properties := struct {
		Properties map[string]string `json:"properties"`
	}{
		config.Properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlV4(endpoint),
		BodyModel:     properties,
		ResponseModel: &email,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &email, nil
}
