package hubspot

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type OwnersResponse struct {
	Results []Owner `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Owner stores Owner from Service
//
type Owner struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	CreatedAt string      `json:"createdAt"`
	Archived  bool        `json:"archived"`
	Teams     []OwnerTeam `json:"teams"`
	ID        string      `json:"id"`
	Email     string      `json:"email"`
	UpdatedAt string      `json:"updatedAt"`
}

type OwnerTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &ownersResponse,
		}

		_, _, e := service.get(&requestConfig)
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
