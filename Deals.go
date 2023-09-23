package hubspot

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
)

type DealsResponse struct {
	Results []Deal  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Deal stores Deal from Service
type Deal struct {
	Id                    string                       `json:"id"`
	Properties            map[string]string            `json:"properties"`
	CreatedAt             h_types.DateTimeString       `json:"createdAt"`
	UpdatedAt             h_types.DateTimeString       `json:"updatedAt"`
	Archived              bool                         `json:"archived"`
	Associations          map[string]AssociationsSet   `json:"associations"`
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory"`
}

type PropertyHistory struct {
	Value           string    `json:"value"`
	Timestamp       time.Time `json:"timestamp"`
	SourceType      string    `json:"sourceType"`
	SourceId        string    `json:"sourceId"`
	UpdatedByUserId int       `json:"updatedByUserId"`
}

type GetDealsConfig struct {
	Limit                 *uint
	After                 *string
	Properties            *[]string
	PropertiesWithHistory *[]string
	Associations          *[]string
	Archived              *bool
}

// GetDeals returns all deals
func (service *Service) GetDeals(config *GetDealsConfig) (*[]Deal, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/deals"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		_properties := []string{}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				_properties = append(_properties, *config.Properties...)
			}
		}
		if config.PropertiesWithHistory != nil {
			if len(*config.PropertiesWithHistory) > 0 {
				values.Set("propertiesWithHistory", strings.Join(*config.PropertiesWithHistory, ","))
			}
		}
		if len(_properties) > 0 {
			values.Set("properties", strings.Join(_properties, ","))
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

	after := ""
	if config.After != nil {
		after = *config.After
	}

	deals := []Deal{}

	for {
		dealsResponse := DealsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &dealsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		deals = append(deals, dealsResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if dealsResponse.Paging == nil {
			break
		}

		if dealsResponse.Paging.Next.After == "" {
			break
		}

		after = dealsResponse.Paging.Next.After
	}

	return &deals, nil
}

type CreateDealConfig struct {
	Properties map[string]string
}

func (service *Service) CreateDeal(config *CreateDealConfig) (*Deal, *errortools.Error) {
	endpoint := "objects/deals"
	deal := Deal{}

	var properties_ = struct {
		Properties map[string]string `json:"properties"`
	}{
		config.Properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     properties_,
		ResponseModel: &deal,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &deal, nil
}

type UpdateDealConfig struct {
	DealId     string
	Properties map[string]string
}

func (service *Service) UpdateDeal(config *UpdateDealConfig) (*Deal, *errortools.Error) {
	endpoint := "objects/deals"
	deal := Deal{}

	var properties_ = struct {
		Properties map[string]string `json:"properties"`
	}{
		config.Properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.DealId)),
		BodyModel:     properties_,
		ResponseModel: &deal,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &deal, nil
}

func (service *Service) BatchDeleteDeals(dealIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(dealIds) > index {
		if len(dealIds) > index+maxItemsPerBatch {
			e := service.batchDeleteDeals(dealIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteDeals(dealIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteDeals(dealIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, dealId := range dealIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{dealId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/deals/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
