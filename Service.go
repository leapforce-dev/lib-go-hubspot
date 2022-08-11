package hubspot

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

const (
	apiName string = "Hubspot"
	apiUrl  string = "https://api.hubapi.com/crm"
)

// type
//
type Service struct {
	apiKey        string
	bearerToken   string
	httpService   *go_http.Service
	errorResponse *ErrorResponse
}

type ServiceConfig struct {
	BearerToken string
}

func NewService(config *ServiceConfig) (*Service, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if config.BearerToken == "" {
		return nil, errortools.ErrorMessage("BearerToken not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		bearerToken: config.BearerToken,
		httpService: httpService,
	}, nil
}
func NewServiceWithApiKey(apiKey string) (*Service, *errortools.Error) {
	if apiKey == "" {
		return nil, errortools.ErrorMessage("apiKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiKey:      apiKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	if service.bearerToken != "" {
		// add authentication header
		header := http.Header{}
		header.Set("Authorization", fmt.Sprintf("Bearer %s", service.bearerToken))
		(*requestConfig).NonDefaultHeaders = &header
	}

	if service.apiKey != "" {
		// add Api key
		_url, err := url.Parse(requestConfig.Url)
		if err != nil {
			return nil, nil, errortools.ErrorMessage(err)
		}
		query := _url.Query()
		query.Set("hapikey", service.apiKey)

		(*requestConfig).Url = fmt.Sprintf("%s://%s%s?%s", _url.Scheme, _url.Host, _url.Path, query.Encode())
	}

	// add error model
	service.errorResponse = &ErrorResponse{}
	(*requestConfig).ErrorModel = &service.errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if service.errorResponse.Message != "" {
		e.SetMessage(service.errorResponse.Message)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/v3/%s", apiUrl, path)
}

func (service *Service) urlV4(path string) string {
	return fmt.Sprintf("%s/v4/%s", apiUrl, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.bearerToken
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}

func (service *Service) ErrorResponse() *ErrorResponse {
	return service.errorResponse
}
