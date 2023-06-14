package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	oauth2 "github.com/leapforce-libraries/go_oauth2"
	"github.com/leapforce-libraries/go_oauth2/tokensource"
	"net/http"
	"net/url"
	"time"
)

const (
	apiName            string = "Hubspot"
	apiUrlEngagements  string = "https://api.hubapi.com/engagements"
	apiUrlFiles        string = "https://api.hubapi.com/files"
	apiUrlCrm          string = "https://api.hubapi.com/crm"
	apiUrlOAuth        string = "https://api.hubapi.com/oauth"
	apiUrlAccountInfo  string = "https://api.hubapi.com/account-info"
	defaultRedirectUrl string = "http://localhost:8080/oauth/redirect"
	authUrl            string = "https://app-eu1.hubspot.com/oauth/authorize"
	//tokenUrl           string = "https://api.pinterest.com/v5/oauth/token"
	tokenHttpMethod string = http.MethodPost
)

type authorizationMode string

const (
	authorizationModeOAuth2      authorizationMode = "oauth2"
	authorizationModeApiKey      authorizationMode = "apikey"
	authorizationModeAccessToken authorizationMode = "accesstoken"
)

type Service struct {
	authorizationMode authorizationMode
	clientId          string
	apiKey            string
	accessToken       string
	httpService       *go_http.Service
	oAuth2Service     *oauth2.Service
	redirectUrl       *string
	errorResponse     *ErrorResponse
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
		authorizationMode: authorizationModeAccessToken,
		accessToken:       config.BearerToken,
		httpService:       httpService,
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
		authorizationMode: authorizationModeApiKey,
		apiKey:            apiKey,
		httpService:       httpService,
	}, nil
}

type ServiceWithOAuth2Config struct {
	ClientId     string
	ClientSecret string
	TokenSource  tokensource.TokenSource
	RedirectUrl  *string
}

func NewServiceWithOAuth2(cfg *ServiceWithOAuth2Config) (*Service, *errortools.Error) {
	if cfg == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if cfg.ClientId == "" {
		return nil, errortools.ErrorMessage("ClientId not provided")
	}

	redirectUrl := defaultRedirectUrl
	if cfg.RedirectUrl != nil {
		redirectUrl = *cfg.RedirectUrl
	}

	oauth2ServiceConfig := oauth2.ServiceConfig{
		ClientId:        cfg.ClientId,
		ClientSecret:    cfg.ClientSecret,
		RedirectUrl:     redirectUrl,
		AuthUrl:         authUrl,
		TokenUrl:        fmt.Sprintf("%s/v1/token", apiUrlOAuth),
		TokenHttpMethod: tokenHttpMethod,
		TokenSource:     cfg.TokenSource,
	}
	oauth2Service, e := oauth2.NewService(&oauth2ServiceConfig)
	if e != nil {
		return nil, e
	}

	return &Service{
		authorizationMode: authorizationModeOAuth2,
		clientId:          cfg.ClientId,
		oAuth2Service:     oauth2Service,
		redirectUrl:       cfg.RedirectUrl,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	return service.tryHttpRequest(requestConfig, false)
}

func (service *Service) tryHttpRequest(requestConfig *go_http.RequestConfig, isRetry bool) (*http.Request, *http.Response, *errortools.Error) {
	var request *http.Request
	var response *http.Response
	var e *errortools.Error

	// add error model
	service.errorResponse = &ErrorResponse{}
	requestConfig.ErrorModel = service.errorResponse

	if service.authorizationMode == authorizationModeOAuth2 {
		request, response, e = service.oAuth2Service.HttpRequest(requestConfig)
	} else {
		if service.authorizationMode == authorizationModeAccessToken {
			// add authentication header
			header := http.Header{}
			header.Set("Authorization", fmt.Sprintf("Bearer %s", service.accessToken))
			(*requestConfig).NonDefaultHeaders = &header
		} else if service.authorizationMode == authorizationModeApiKey {
			// add Api key
			_url, err := url.Parse(requestConfig.Url)
			if err != nil {
				return nil, nil, errortools.ErrorMessage(err)
			}
			query := _url.Query()
			query.Set("hapikey", service.apiKey)

			(*requestConfig).Url = fmt.Sprintf("%s://%s%s?%s", _url.Scheme, _url.Host, _url.Path, query.Encode())
		}

		request, response, e = service.httpService.HttpRequest(requestConfig)
	}

	if e != nil {
		if response.StatusCode == http.StatusTooManyRequests && !isRetry {
			remaining := response.Header["X-Hubspot-Ratelimit-Daily-Remaining"]

			if len(remaining) == 1 {
				if remaining[0] != "0" {
					// try to catch the per second rate limit, but try this only once (isRetry)
					fmt.Println("waiting 1 second...")
					time.Sleep(time.Second)

					return service.tryHttpRequest(requestConfig, true)
				}
			}
		}
		if service.errorResponse.Message != "" {
			e.SetMessage(service.errorResponse.Message)
		}
	}

	if e != nil {
		return request, response, e
	}

	return request, response, nil
}

func (service *Service) AuthorizeUrl(scope string) string {
	if service.redirectUrl == nil {
		return ""
	}
	return fmt.Sprintf("https://app-eu1.hubspot.com/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s", service.clientId, *service.redirectUrl, scope)
}

func (service *Service) GetTokenFromCode(r *http.Request) *errortools.Error {
	return service.oAuth2Service.GetTokenFromCode(r, nil)
}

func (service *Service) urlEngagements(path string) string {
	return fmt.Sprintf("%s/v1/%s", apiUrlEngagements, path)
}

func (service *Service) urlFiles(path string) string {
	return fmt.Sprintf("%s/v3/%s", apiUrlFiles, path)
}

func (service *Service) urlCrm(path string) string {
	return fmt.Sprintf("%s/v3/%s", apiUrlCrm, path)
}

func (service *Service) urlOAuth(path string) string {
	return fmt.Sprintf("%s/v1/%s", apiUrlOAuth, path)
}

func (service *Service) urlAccountInfo(path string) string {
	return fmt.Sprintf("%s/v3/%s", apiUrlAccountInfo, path)
}

func (service *Service) urlV4(path string) string {
	return fmt.Sprintf("%s/v4/%s", apiUrlCrm, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.accessToken
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
