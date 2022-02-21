package hubspot

import (
	"fmt"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	go_types "github.com/leapforce-libraries/go_types"
)

type OwnersResponse struct {
	Results []Owner `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Owner stores Owner from Service
//
type Owner struct {
	FirstName string                   `json:"firstName"`
	LastName  string                   `json:"lastName"`
	CreatedAt h_types.DateTimeMSString `json:"createdAt"`
	Archived  bool                     `json:"archived"`
	Teams     []OwnerTeam              `json:"teams"`
	ID        string                   `json:"id"`
	Email     string                   `json:"email"`
	UpdatedAt h_types.DateTimeMSString `json:"updatedAt"`
}

type OwnerTeam struct {
	ID   go_types.Int64String `json:"id"`
	Name string               `json:"name"`
}

type GetOwnersConfig struct {
	Limit *uint
	After *string
	Email *string
}

// GetOwners returns all owners
//
func (service *Service) GetOwners(config *GetOwnersConfig) (*[]Owner, *errortools.Error) {
	values := url.Values{}
	endpoint := "owners"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		if config.Email != nil {
			values.Set("email", *config.Email)
		}
	}

	after := ""
	if config.After != nil {
		after = *config.After
	}

	owners := []Owner{}

	for {
		ownersResponse := OwnersResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.url(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &ownersResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		owners = append(owners, ownersResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if ownersResponse.Paging == nil {
			break
		}

		if ownersResponse.Paging.Next.After == "" {
			break
		}

		after = ownersResponse.Paging.Next.After
	}

	return &owners, nil
}
