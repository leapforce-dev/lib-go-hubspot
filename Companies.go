package hubspot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
)

type CompaniesResponse struct {
	Results []Company `json:"results"`
	Paging  *Paging   `json:"paging"`
}

// Company stores Company from Service
type Company struct {
	Id                    string                       `json:"id"`
	Properties            map[string]string            `json:"properties"`
	CreatedAt             h_types.DateTimeMSString     `json:"createdAt"`
	UpdatedAt             h_types.DateTimeMSString     `json:"updatedAt"`
	Archived              bool                         `json:"archived"`
	Associations          map[string]AssociationsSet   `json:"associations"`
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory"`
}

type GetCompaniesConfig struct {
	Limit                 *uint
	After                 *string
	Properties            *[]string
	PropertiesWithHistory *[]string
	Associations          *[]string
	Archived              *bool
}

// GetCompanies returns all companies
func (service *Service) GetCompanies(config *GetCompaniesConfig) (*[]Company, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/companies"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				values.Set("properties", strings.Join(*config.Properties, ","))
			}
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

	companies := []Company{}

	for {
		companiesResponse := CompaniesResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &companiesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		companies = append(companies, companiesResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if companiesResponse.Paging == nil {
			break
		}

		if companiesResponse.Paging.Next.After == "" {
			break
		}

		after = companiesResponse.Paging.Next.After
	}

	return &companies, nil
}

func (service *Service) CreateCompany(config *CreateObjectConfig) (*Company, *errortools.Error) {
	endpoint := "objects/companies"
	company := Company{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &company,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &company, nil
}

func (service *Service) BatchCreateCompanies(config *BatchObjectsConfig, invalidEmailProperty string) (*[]Company, *errortools.Error) {
	var companies []Company
	var retrying = false

	for _, batch := range service.batches(len(config.Inputs)) {
	retry:
		var r BatchCompaniesResponse

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

			if response.StatusCode == http.StatusBadRequest && !retrying {
				errorResponse := service.ErrorResponse()
				if errorResponse != nil {
					if strings.HasPrefix(errorResponse.Message, "Property values were not valid: ") {
						stop := checkInvalidEmails(config, invalidEmailProperty, errorResponse.Message, batch)

						if stop {
							goto stop
						}

						retrying = true
						goto retry
					}
				}
			}
		}

	stop:
		if e != nil {
			return nil, e
		}
	ok:
		companies = append(companies, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &companies, nil
}

type BatchCompaniesResponse struct {
	CompletedAt *time.Time        `json:"completedAt"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt *time.Time        `json:"requestedAt"`
	StartedAt   *time.Time        `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []Company         `json:"results"`
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

func (service *Service) BatchUpdateCompanies(config *BatchObjectsConfig, invalidEmailProperty string) (*[]Company, *errortools.Error) {
	var companies []Company

	for _, batch := range service.batches(len(config.Inputs)) {
	retry:
		var r BatchCompaniesResponse

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
			if response.StatusCode == http.StatusBadRequest {
				errorResponse := service.ErrorResponse()
				if errorResponse != nil {
					if strings.HasPrefix(errorResponse.Message, "Property values were not valid: ") {
						stop := checkInvalidEmails(config, invalidEmailProperty, errorResponse.Message, batch)

						if stop {
							goto stop
						}

						goto retry
					}
				}
			}
		}
	stop:
		if e != nil {
			return nil, e
		}
	ok:
		companies = append(companies, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &companies, nil
}

func (service *Service) UpdateCompany(config *UpdateObjectConfig) (*Company, *errortools.Error) {
	endpoint := "objects/companies"
	company := Company{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.ObjectId)),
		BodyModel:     config,
		ResponseModel: &company,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &company, nil
}

type GetCompanyConfig struct {
	CompanyId    string
	Properties   *[]string
	Associations *[]string
}

// GetCompany returns a specific company
func (service *Service) GetCompany(config *GetCompanyConfig) (*Company, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/companies"

	if config == nil {
		return nil, errortools.ErrorMessage("config is nil")
	}

	_properties := []string{}
	if config.Properties != nil {
		if len(*config.Properties) > 0 {
			_properties = append(_properties, *config.Properties...)
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

	company := Company{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s?%s", endpoint, config.CompanyId, values.Encode())),
		ResponseModel: &company,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &company, nil
}

type FilterGroup struct {
	Filters *[]filter `json:"filters"`
}

func (fg *FilterGroup) AddPropertyFilter(operator string, property string, value string, highValue string) {
	if fg.Filters == nil {
		fg.Filters = &[]filter{}
	}

	*fg.Filters = append(*fg.Filters, filter{
		Operator:     operator,
		PropertyName: property,
		Value:        value,
		HighValue:    highValue,
		isCustom:     false,
	})
}

func (fg *FilterGroup) AddCustomPropertyFilter(operator string, propertyName string, value string, highValue string) {
	if fg.Filters == nil {
		fg.Filters = &[]filter{}
	}

	*fg.Filters = append(*fg.Filters, filter{
		Operator:     operator,
		PropertyName: propertyName,
		Value:        value,
		HighValue:    highValue,
		isCustom:     true,
	})
}

type filter struct {
	Operator     string `json:"operator"`
	PropertyName string `json:"propertyName,omitempty"`
	Value        string `json:"value"`
	HighValue    string `json:"highValue,omitempty"`
	isCustom     bool   `json:"-"`
}

// SearchCompanies returns a specific company
func (service *Service) SearchCompanies(config *SearchObjectsConfig) (*[]Company, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config is nil")
	}

	endpoint := "objects/companies/search"

	companiesResponse := CompaniesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &companiesResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	after := config.After

	companies := []Company{}

	for {
		companiesResponse := CompaniesResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm("objects/companies/search"),
			BodyModel:     config,
			ResponseModel: &companiesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		companies = append(companies, companiesResponse.Results...)

		if after != nil { // explicit after parameter requested
			break
		}

		if companiesResponse.Paging == nil {
			break
		}

		if companiesResponse.Paging.Next.After == "" {
			break
		}

		config.After = &companiesResponse.Paging.Next.After
	}

	return &companies, nil
}

func (service *Service) ArchiveCompany(companyId string) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.urlCrm(fmt.Sprintf("objects/companies/%s", companyId)),
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}

func (service *Service) BatchArchiveCompanies(companyIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(companyIds) > index {
		if len(companyIds) > index+maxItemsPerBatch {
			e := service.batchArchiveCompanies(companyIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchArchiveCompanies(companyIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchArchiveCompanies(companyIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, companyId := range companyIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{companyId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/companies/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
