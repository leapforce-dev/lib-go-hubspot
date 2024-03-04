package hubspot

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func (service *Service) CreateDeal(config *CreateObjectConfig) (*Deal, *errortools.Error) {
	endpoint := "objects/deals"
	deal := Deal{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &deal,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &deal, nil
}

type BatchDealsResponse struct {
	CompletedAt *time.Time        `json:"completedAt"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt *time.Time        `json:"requestedAt"`
	StartedAt   *time.Time        `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []Deal            `json:"results"`
	Errors      []struct {
		SubCategory json.RawMessage   `json:"subCategory"`
		Context     map[string]string `json:"context"`
		Links       map[string]string `json:"links"`
		Id          string            `json:"id"`
		Category    string            `json:"category"`
		Message     string            `json:"message"`
		Errors      []struct {
			SubCategory string `json:"subCategory"`
			Code        string `json:"code"`
			In          string `json:"in"`
			Context     struct {
				MissingScopes []string `json:"missingScopes"`
			} `json:"context"`
			Message string `json:"message"`
		} `json:"errors"`
		Status string `json:"status"`
	} `json:"errors"`
	Status string `json:"status"`
}

func (service *Service) BatchCreateDeals(config *BatchObjectsConfig) (*[]Deal, *errortools.Error) {
	var deals []Deal

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchDealsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(fmt.Sprintf("objects/%s/batch/create", config.ObjectType)),
			BodyModel:     BatchObjectsConfig{Inputs: config.Inputs[batch.startIndex:batch.endIndex]},
			ResponseModel: &r,
		}

		_, response, e := service.httpRequest(&requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusMultiStatus {
				fmt.Println(r.Errors)
				goto ok
			}
		}
		if e != nil {
			return nil, e
		}
	ok:
		deals = append(deals, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &deals, nil
}

func (service *Service) BatchUpdateDeals(config *BatchObjectsConfig) (*[]Deal, *errortools.Error) {
	var deals []Deal

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchDealsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(fmt.Sprintf("objects/%s/batch/update", config.ObjectType)),
			BodyModel:     BatchObjectsConfig{Inputs: config.Inputs[batch.startIndex:batch.endIndex]},
			ResponseModel: &r,
		}

		_, response, e := service.httpRequest(&requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusMultiStatus {
				fmt.Println(r.Errors)
				goto ok
			}
		}
		if e != nil {
			return nil, e
		}
	ok:
		deals = append(deals, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &deals, nil
}

func (service *Service) UpdateDeal(config *UpdateObjectConfig) (*Deal, *errortools.Error) {
	endpoint := "objects/deals"
	deal := Deal{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.ObjectId)),
		BodyModel:     config,
		ResponseModel: &deal,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &deal, nil
}

func (service *Service) BatchArchiveDeals(dealIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(dealIds) > index {
		if len(dealIds) > index+maxItemsPerBatch {
			e := service.batchArchiveDeals(dealIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchArchiveDeals(dealIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchArchiveDeals(dealIds []string) *errortools.Error {
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
