package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type AccessToken struct {
	Token                     string   `json:"token"`
	User                      string   `json:"user"`
	HubDomain                 string   `json:"hub_domain"`
	Scopes                    []string `json:"scopes"`
	ScopeToScopeGroupPks      []int    `json:"scope_to_scope_group_pks"`
	TrialScopes               []int    `json:"trial_scopes"`
	TrialScopeToScopeGroupPks []int    `json:"trial_scope_to_scope_group_pks"`
	HubId                     int      `json:"hub_id"`
	AppId                     int      `json:"app_id"`
	ExpiresIn                 int      `json:"expires_in"`
	UserId                    int      `json:"user_id"`
	TokenType                 string   `json:"token_type"`
}

// InspectAccessToken returns information about an access token
func (service *Service) InspectAccessToken(accessToken string) (*AccessToken, *errortools.Error) {
	var a AccessToken

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlOAuth(fmt.Sprintf("access-tokens/%s", accessToken)),
		ResponseModel: &a,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &a, nil
}
