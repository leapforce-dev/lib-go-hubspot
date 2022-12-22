package hubspot

import (
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

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
