package hubspot

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	APIURL string = "https://api.hubapi.com/crm/v3"
	//DateFormat       string = "2006-01-02T15:04:05"
)

// type
//
type Service struct {
	apiKey      string
	httpService *go_http.Service
}

type ServiceConfig struct {
	APIKey string
}

func NewService(config ServiceConfig) (*Service, *errortools.Error) {
	if config.APIKey == "" {
		return nil, errortools.ErrorMessage("APIKey not provided")
	}

	httpServiceConfig := go_http.ServiceConfig{}

	return &Service{
		apiKey:      config.APIKey,
		httpService: go_http.NewService(httpServiceConfig),
	}, nil
}

func (service *Service) httpRequest(httpMethod string, requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add API key
	_url, err := url.Parse(requestConfig.URL)
	if err != nil {
		return nil, nil, errortools.ErrorMessage(err)
	}
	_url.Query().Set("hapikey", service.apiKey)

	(*requestConfig).URL = _url.String()

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HTTPRequest(httpMethod, requestConfig)
	if errorResponse.Message != "" {
		e.SetMessage(errorResponse.Message)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", APIURL, path)
}

func (service *Service) get(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodGet, requestConfig)
}

func (service *Service) post(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPost, requestConfig)
}

func (service *Service) put(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodPut, requestConfig)
}

func (service *Service) delete(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.httpRequest(http.MethodDelete, requestConfig)
}
