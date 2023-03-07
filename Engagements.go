package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"time"
)

type EngagementsResponse struct {
	Results []Engagement `json:"results"`
	HasMore bool         `json:"hasMore"`
	Offset  int64        `json:"offset"`
}

type Engagement struct {
	Engagement   EngagementEngagement   `json:"engagement"`
	Associations EngagementAssociations `json:"associations"`
	Attachments  []Attachment           `json:"attachments"`
	Metadata     EngagementMetadata     `json:"metadata"`
}

type EngagementEngagement struct {
	Id                     int64    `json:"id"`
	PortalId               int64    `json:"portalId"`
	Active                 bool     `json:"active"`
	CreatedAt              *int64   `json:"createdAt,omitempty"`
	LastUpdated            *int64   `json:"lastUpdated,omitempty"`
	Type                   string   `json:"type"`
	Timestamp              *int64   `json:"timestamp,omitempty"`
	AllAccessibleTeamIds   *[]int64 `json:"allAccessibleTeamIds,omitempty"`
	QueueMembershipIds     *[]int64 `json:"queueMembershipIds,omitempty"`
	BodyPreviewIsTruncated *bool    `json:"bodyPreviewIsTruncated,omitempty"`
}

type EngagementAssociations struct {
	ContactIds        []int64 `json:"contactIds"`
	CompanyIds        []int64 `json:"companyIds"`
	DealIds           []int64 `json:"dealIds"`
	OwnerIds          []int64 `json:"ownerIds"`
	WorkflowIds       []int64 `json:"workflowIds"`
	TicketIds         []int64 `json:"ticketIds"`
	ContentIds        []int64 `json:"contentIds"`
	QuoteIds          []int64 `json:"quoteIds"`
	MarketingEventIds []int64 `json:"marketingEventIds"`
}

type Attachment struct {
	Id int64 `json:"id"`
}

type EngagementMetadata struct {
	DurationMilliseconds int64  `json:"durationMilliseconds"`
	Body                 string `json:"body"`
}

func (service *Service) CreateEngagement(engagement *Engagement) (*Engagement, *errortools.Error) {
	var newEngagement Engagement

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlEngagements("engagements"),
		BodyModel:     engagement,
		ResponseModel: &newEngagement,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &newEngagement, nil
}

type GetAllEngagementsConfig struct {
	Limit uint
}

func (service *Service) GetAllEngagements(config *GetAllEngagementsConfig) (*[]Engagement, *errortools.Error) {
	values := url.Values{}
	if config != nil {
		values.Set("limit", fmt.Sprintf("%v", config.Limit))
	}

	var engagements []Engagement

	var offset int64 = 0

	for {
		if offset > 0 {
			values.Set("offset", fmt.Sprintf("%v", offset))
		}

		var engagementsResponse EngagementsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlEngagements(fmt.Sprintf("engagements/paged?%s", values.Encode())),
			ResponseModel: &engagementsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		engagements = append(engagements, engagementsResponse.Results...)

		if !engagementsResponse.HasMore {
			break
		}

		offset = engagementsResponse.Offset
	}

	return &engagements, nil
}

type GetRecentEngagementsConfig struct {
	Count *uint
	Since *time.Time
}

func (service *Service) GetRecentEngagements(config *GetRecentEngagementsConfig) (*[]Engagement, *errortools.Error) {
	values := url.Values{}
	if config != nil {
		if config.Count != nil {
			values.Set("count", fmt.Sprintf("%v", *config.Count))
		}
		if config.Since != nil {
			values.Set("since", fmt.Sprintf("%v", config.Since.UnixMilli()))
		}
	}

	var engagements []Engagement

	var offset int64 = 0

	for {
		if offset > 0 {
			values.Set("offset", fmt.Sprintf("%v", offset))
		}

		var engagementsResponse EngagementsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlEngagements(fmt.Sprintf("engagements/recent/modified?%s", values.Encode())),
			ResponseModel: &engagementsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		engagements = append(engagements, engagementsResponse.Results...)

		if !engagementsResponse.HasMore {
			break
		}

		offset = engagementsResponse.Offset
	}

	return &engagements, nil
}
