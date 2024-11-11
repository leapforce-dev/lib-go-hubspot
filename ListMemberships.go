package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"time"
)

type ListMembershipsResponse struct {
	Results []ListMembership `json:"results"`
	Paging  *Paging          `json:"paging"`
}

type ListMembership struct {
	RecordId            string    `json:"recordId"`
	MembershipTimestamp time.Time `json:"membershipTimestamp"`
}

type GetListMembershipsConfig struct {
	ListId int64
	Limit  *uint
	After  *string
}

// GetListMemberships returns all listMemberships
func (service *Service) GetListMemberships(config *GetListMembershipsConfig) (*[]ListMembership, *errortools.Error) {
	values := url.Values{}

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
	}

	after := ""
	if config.After != nil {
		after = *config.After
	}

	var listMemberships []ListMembership

	for {
		listMembershipsResponse := ListMembershipsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("lists/%v/memberships?%s", config.ListId, values.Encode())),
			ResponseModel: &listMembershipsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		listMemberships = append(listMemberships, listMembershipsResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if listMembershipsResponse.Paging == nil {
			break
		}

		if listMembershipsResponse.Paging.Next.After == "" {
			break
		}

		after = listMembershipsResponse.Paging.Next.After
	}

	return &listMemberships, nil
}
