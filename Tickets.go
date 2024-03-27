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

type TicketsResponse struct {
	Results []Ticket `json:"results"`
	Paging  *Paging  `json:"paging"`
}

// Ticket stores Ticket from Service
type Ticket struct {
	Id           string                     `json:"id"`
	Properties   map[string]string          `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type ListTicketsConfig struct {
	Limit        *uint
	After        *string
	Properties   *[]string
	Associations *[]string
	Archived     *bool
}

// ListTickets returns all tickets
func (service *Service) ListTickets(config *ListTicketsConfig) (*[]Ticket, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/tickets"
	after := ""

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		var _properties []string
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
				var _associations []string
				for _, a := range *config.Associations {
					_associations = append(_associations, a)
				}
				values.Set("associations", strings.Join(_associations, ","))
			}
		}
		if config.Archived != nil {
			values.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}

		if config.After != nil {
			after = *config.After
		}
	}

	var tickets []Ticket

	for {
		ticketsResponse := TicketsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &ticketsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		tickets = append(tickets, ticketsResponse.Results...)

		if config != nil {
			if config.After != nil { // explicit after parameter requested
				break
			}
		}

		if ticketsResponse.Paging == nil {
			break
		}

		if ticketsResponse.Paging.Next.After == "" {
			break
		}

		after = ticketsResponse.Paging.Next.After
	}

	return &tickets, nil
}

type CreateTicketConfig struct {
	Properties   map[string]string  `json:"properties"`
	Associations *[]AssociationToV4 `json:"associations,omitempty"`
}

func (service *Service) CreateTicket(config *CreateTicketConfig) (*Ticket, *errortools.Error) {
	endpoint := "objects/tickets"
	ticket := Ticket{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &ticket,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &ticket, nil
}

type UpdateTicketConfig struct {
	TicketId   string
	Properties map[string]string
}

func (service *Service) UpdateTicket(config *UpdateTicketConfig) (*Ticket, *errortools.Error) {
	endpoint := "objects/tickets"
	ticket := Ticket{}

	var properties = make(map[string]string)

	if config.Properties != nil {
		for key, value := range config.Properties {
			properties[key] = value
		}
	}

	var properties_ = struct {
		Properties map[string]string `json:"properties"`
	}{
		properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.TicketId)),
		BodyModel:     properties_,
		ResponseModel: &ticket,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &ticket, nil
}

func (service *Service) BatchArchiveTickets(ticketIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(ticketIds) > index {
		if len(ticketIds) > index+maxItemsPerBatch {
			e := service.batchArchiveTickets(ticketIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchArchiveTickets(ticketIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchArchiveTickets(ticketIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, ticketId := range ticketIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{ticketId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/tickets/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}

// SearchTickets returns a specific ticket
func (service *Service) SearchTickets(config *SearchObjectsConfig) (*[]Ticket, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config is nil")
	}

	endpoint := "objects/tickets/search"

	ticketsResponse := TicketsResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &ticketsResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	after := config.After

	var tickets []Ticket

	for {
		ticketsResponse := TicketsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(endpoint),
			BodyModel:     config,
			ResponseModel: &ticketsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		tickets = append(tickets, ticketsResponse.Results...)

		if after != nil { // explicit after parameter requested
			break
		}

		if ticketsResponse.Paging == nil {
			break
		}

		if ticketsResponse.Paging.Next.After == "" {
			break
		}

		config.After = &ticketsResponse.Paging.Next.After
	}

	return &tickets, nil
}

type BatchTicketsResponse struct {
	CompletedAt *time.Time        `json:"completedAt"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt *time.Time        `json:"requestedAt"`
	StartedAt   *time.Time        `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []Ticket          `json:"results"`
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

func (service *Service) BatchCreateTickets(config *BatchObjectsConfig) (*[]Ticket, *errortools.Error) {
	var tickets []Ticket

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchTicketsResponse

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
		tickets = append(tickets, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &tickets, nil
}

func (service *Service) BatchUpdateTickets(config *BatchObjectsConfig) (*[]Ticket, *errortools.Error) {
	var tickets []Ticket

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchTicketsResponse

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
		tickets = append(tickets, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &tickets, nil
}
