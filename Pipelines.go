package hubspot

import (
	"fmt"
	go_types "github.com/leapforce-libraries/go_types"
	"net/http"
	"net/url"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
)

type PipelinesResponse struct {
	Results []Pipeline `json:"results"`
}

// Pipeline stores Pipeline from Service
type Pipeline struct {
	Label        string                    `json:"label"`
	DisplayOrder int64                     `json:"displayOrder"`
	Id           string                    `json:"id,omitempty"`
	Stages       []PipelineStage           `json:"stages"`
	CreatedAt    *h_types.DateTimeMSString `json:"createdAt,omitempty"`
	UpdatedAt    *h_types.DateTimeString   `json:"updatedAt,omitempty"`
	Archived     *bool                     `json:"archived,omitempty"`
}

type PipelineStage struct {
	Label        string                    `json:"label"`
	DisplayOrder int64                     `json:"displayOrder"`
	MetaData     PipelineStageMetaData     `json:"metadata"`
	Id           string                    `json:"id,omitempty"`
	CreatedAt    *h_types.DateTimeMSString `json:"createdAt,omitempty"`
	UpdatedAt    *h_types.DateTimeMSString `json:"updatedAt,omitempty"`
	Archived     *bool                     `json:"archived,omitempty"`
}

type PipelineStageMetaData struct {
	IsClosed    *go_types.BoolString    `json:"isClosed"`
	Probability *go_types.Float64String `json:"probability"`
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
		Url:           service.urlCrm(fmt.Sprintf("%s/%s?%s", endpoint, string(config.ObjectType), values.Encode())),
		ResponseModel: &pipelinesResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pipelinesResponse.Results, nil
}

type CreatePipelineConfig struct {
	ObjectType PipelineObjectType
	Pipeline   *Pipeline
}

// CreatePipeline creates a new pipeline
func (service *Service) CreatePipeline(config *CreatePipelineConfig) (*Pipeline, *errortools.Error) {
	var pipelineNew Pipeline

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(fmt.Sprintf("pipelines/%s", config.ObjectType)),
		BodyModel:     config.Pipeline,
		ResponseModel: &pipelineNew,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pipelineNew, nil
}

type UpdatePipelineConfig struct {
	ObjectType PipelineObjectType
	PipelineId string
	Pipeline   *Pipeline
}

// UpdatePipeline updates an existing pipeline
func (service *Service) UpdatePipeline(config *UpdatePipelineConfig) (*Pipeline, *errortools.Error) {
	var pipelineNew Pipeline

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPut,
		Url:           service.urlCrm(fmt.Sprintf("pipelines/%s/%s", config.ObjectType, config.PipelineId)),
		BodyModel:     config.Pipeline,
		ResponseModel: &pipelineNew,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &pipelineNew, nil
}
