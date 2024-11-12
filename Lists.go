package hubspot

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"time"
)

type SearchListsResponse struct {
	Offset  uint32 `json:"offset"`
	HasMore bool   `json:"hasMore"`
	Lists   []List `json:"lists"`
	Total   int32  `json:"total"`
}

// List stores List from Service
type List struct {
	ListId               string    `json:"listId"`
	ListVersion          int       `json:"listVersion"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	FiltersUpdatedAt     time.Time `json:"filtersUpdatedAt"`
	ProcessingStatus     string    `json:"processingStatus"`
	CreatedById          string    `json:"createdById"`
	UpdatedById          string    `json:"updatedById"`
	ProcessingType       string    `json:"processingType"`
	ObjectTypeId         string    `json:"objectTypeId"`
	Name                 string    `json:"name"`
	AdditionalProperties struct {
		HsLastRecordAddedAt  string `json:"hs_last_record_added_at"`
		HsListSize           string `json:"hs_list_size"`
		HsListReferenceCount string `json:"hs_list_reference_count"`
	} `json:"additionalProperties"`
}

type SearchListsConfig struct {
	Offset          *uint32   `json:"offset"`
	Query           *string   `json:"query,omitempty"`
	ProcessingTypes *[]string `json:"processingTypes,omitempty"`
}

// SearchLists returns all lists
func (service *Service) SearchLists(config *SearchListsConfig) (*[]List, *errortools.Error) {
	endpoint := "lists/search"

	var config_ SearchListsConfig
	if config != nil {
		config_ = *config
	}

	var lists []List

	for {
		listsResponse := SearchListsResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(endpoint),
			BodyModel:     config_,
			ResponseModel: &listsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		lists = append(lists, listsResponse.Lists...)

		if !listsResponse.HasMore {
			break
		}

		config_.Offset = &listsResponse.Offset
	}

	return &lists, nil
}
