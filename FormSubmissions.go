package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type FormSubmissionsResponse struct {
	Results []FormSubmission `json:"results"`
	Paging  *Paging          `json:"paging"`
}

type FormSubmission struct {
	ConversionId string                `json:"conversionId"`
	SubmittedAt  int64                 `json:"submittedAt"`
	Values       []FormSubmissionValue `json:"values"`
	PageUrl      string                `json:"pageUrl"`
}

type FormSubmissionValue struct {
	Name         string `json:"name"`
	Value        string `json:"value"`
	ObjectTypeId string `json:"objectTypeId"`
}

type GetFormSubmissionsConfig struct {
	Limit *uint
	After *string
}

// GetFormSubmissions returns all formSubmissions
func (service *Service) GetFormSubmissions(formId string, config *GetFormSubmissionsConfig) (*[]FormSubmission, *errortools.Error) {
	values := url.Values{}
	endpoint := "submissions/forms"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
	}

	after := ""
	if config.After != nil {
		after = *config.After
	}

	formSubmissions := []FormSubmission{}

	for {
		formSubmissionsResponse := FormSubmissionsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlFormIntegrations(fmt.Sprintf("%s/%s?%s", endpoint, formId, values.Encode())),
			ResponseModel: &formSubmissionsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		formSubmissions = append(formSubmissions, formSubmissionsResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if formSubmissionsResponse.Paging == nil {
			break
		}

		if formSubmissionsResponse.Paging.Next.After == "" {
			break
		}

		after = formSubmissionsResponse.Paging.Next.After
	}

	return &formSubmissions, nil
}
