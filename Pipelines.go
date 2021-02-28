package hubspot

import (
	"fmt"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type PipelinesResponse struct {
	Results []Pipeline `json:"results"`
}

// Pipeline stores Pipeline from Service
//
type Pipeline struct {
	Label        string          `json:"label"`
	DisplayOrder int             `json:"displayOrder"`
	ID           string          `json:"id"`
	Stages       []PipelineStage `json:"stages"`
	CreatedAt    string          `json:"createdAt"`
	UpdatedAt    string          `json:"updatedAt"`
	Archived     bool            `json:"archived"`
}

type PipelineStage struct {
	Label        string `json:"label"`
	DisplayOrder int    `json:"displayOrder"`
	MetaData     struct {
		IsClosed    *string `json:"isClosed"`
		Probability *string `json:"probability"`
	} `json:"metadata"`
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	Archived  bool   `json:"archived"`
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
		URL:           service.url(fmt.Sprintf("%s/%s?%s", endpoint, string(config.ObjectType), values.Encode())),
		ResponseModel: &pipelinesResponse,
	}

	_, _, e := service.get(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pipelinesResponse.Results, nil
}
