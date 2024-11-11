package hubspot

import (
	"fmt"
	"net/http"
	"time"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Object struct {
	CreatedAt             time.Time                    `json:"createdAt"`
	Archived              bool                         `json:"archived"`
	ArchivedAt            time.Time                    `json:"archivedAt"`
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory"`
	Id                    string                       `json:"id"`
	Properties            map[string]string            `json:"properties"`
	UpdatedAt             time.Time                    `json:"updatedAt"`
}

type BatchGetObjectsInput struct {
	Id string `json:"id"`
}

type BatchGetObjectsConfig struct {
	ObjectType            string
	PropertiesWithHistory []string               `json:"propertiesWithHistory,omitempty"`
	IdProperty            string                 `json:"idProperty,omitempty"`
	Inputs                []BatchGetObjectsInput `json:"inputs"`
	Properties            []string               `json:"properties,omitempty"`
}

type BatchGetObjectsResponse struct {
	CompletedAt time.Time         `json:"completedAt"`
	RequestedAt time.Time         `json:"requestedAt"`
	StartedAt   time.Time         `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []Object          `json:"results"`
	Status      string            `json:"status"`
}

func (service *Service) BatchGetObjects(config *BatchGetObjectsConfig) (*[]Object, *errortools.Error) {
	if config == nil {
		return nil, nil
	}
	if len(config.Inputs) == 0 {
		return nil, nil
	}

	endpoint := fmt.Sprintf("objects/%s/batch/read", config.ObjectType)

	var objects []Object

	for _, batch := range service.batches(len(config.Inputs)) {
		var batchGetObjectsResponse BatchGetObjectsResponse

		requestConfig := go_http.RequestConfig{
			Method: http.MethodPost,
			Url:    service.urlCrm(endpoint),
			BodyModel: BatchGetObjectsConfig{
				PropertiesWithHistory: config.PropertiesWithHistory,
				IdProperty:            config.IdProperty,
				Inputs:                config.Inputs[batch.startIndex:batch.endIndex],
				Properties:            config.Properties,
			},
			ResponseModel: &batchGetObjectsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		objects = append(objects, batchGetObjectsResponse.Results...)
	}

	return &objects, nil
}
