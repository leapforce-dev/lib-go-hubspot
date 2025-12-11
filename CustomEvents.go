package hubspot

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"time"
)

type SendEventDataConfig struct {
	EventName  string            `json:"eventName"`
	ObjectId   string            `json:"objectId"`
	OccurredAt *time.Time        `json:"occurredAt,omitempty"`
	Properties map[string]string `json:"properties"`
}

func (service *Service) SendEventData(config *SendEventDataConfig) *errortools.Error {
	endpoint := "send"

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm(endpoint),
		BodyModel: config,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
