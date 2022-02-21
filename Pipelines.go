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

type PipelinesResponse struct {
	Results []Pipeline `json:"results"`
}

// Pipeline stores Pipeline from Service
//
type Pipeline struct {
	Label        string                   `json:"label"`
	DisplayOrder int64                    `json:"displayOrder"`
	ID           string                   `json:"id"`
	Stages       []PipelineStage          `json:"stages"`
	CreatedAt    h_types.DateTimeMSString `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString   `json:"updatedAt"`
	Archived     bool                     `json:"archived"`
}

type PipelineStage struct {
	Label        string `json:"label"`
	DisplayOrder int    `json:"displayOrder"`
	MetaData     struct {
		IsClosed    *go_types.BoolString    `json:"isClosed"`
		Probability *go_types.Float64String `json:"probability"`
	} `json:"metadata"`
	ID        string                   `json:"id"`
	CreatedAt h_types.DateTimeMSString `json:"createdAt"`
	UpdatedAt h_types.DateTimeMSString `json:"updatedAt"`
	Archived  bool                     `json:"archived"`
}

type PipelineObjectType string

const (
	PipelineObjectTypeDeals   PipelineObjectType = "deals"
	PipelineObjectTypeTickets PipelineObjectType = "tickets"
)

type GetPipelinesConfig struct {
	ObjectType PipelineObjectType
	Archived   *bool
}

// GetPipelines returns all pipelines
//
func (service *Service) GetPipelines(config *GetPipelinesConfig) (*[]Pipeline, *errortools.Error) {
	values := url.Values{}
	endpoint := "pipelines"

	if config != nil {
		if config.Archived != nil {
			values.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}
	}

	pipelinesResponse := PipelinesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("%s/%s?%s", endpoint, string(config.ObjectType), values.Encode())),
		ResponseModel: &pipelinesResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pipelinesResponse.Results, nil
}
