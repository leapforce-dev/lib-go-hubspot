package hubspot

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type AccountInfoDetails struct {
	PortalId int64 `json:"portalId"`
}

func (service *Service) GetAccountInfoDetails() (*AccountInfoDetails, *errortools.Error) {
	var accountInfoDetails AccountInfoDetails

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlAccountInfo("details"),
		ResponseModel: &accountInfoDetails,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &accountInfoDetails, nil
}
