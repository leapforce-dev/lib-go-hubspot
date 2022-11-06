package hubspot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	go_types "github.com/leapforce-libraries/go_types"
)

type DealsResponse struct {
	Results []deal  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Deal stores Deal from Service
//
type deal struct {
	Id           string                     `json:"id"`
	Properties   json.RawMessage            `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}
type Deal struct {
	Id               string
	Properties       DealProperties
	CustomProperties map[string]string
	CreatedAt        h_types.DateTimeString
	UpdatedAt        h_types.DateTimeString
	Archived         bool
	Associations     map[string]AssociationsSet
}

type DealProperties struct {
	Amount              *go_types.Float64String   `json:"amount"`
	Assist              *string                   `json:"assist"`
	Category            *string                   `json:"category"`
	CloseDate           *h_types.DateTimeMSString `json:"closedate"`
	CreateDate          *h_types.DateTimeMSString `json:"createdate"`
	DealName            *string                   `json:"dealname"`
	DealStage           *string                   `json:"dealstage"`
	ForecastAmount      *go_types.Float64String   `json:"hs_forecast_amount"`
	ForecastProbability *go_types.Float64String   `json:"hs_forecast_probability"`
	LastUpdated         *h_types.DateTimeMSString `json:"notes_last_updated"`
	OwnerId             *string                   `json:"hubspot_owner_id"`
}

type DealProperty string

const (
	DealPropertyAmountInHomeCurrency                                 DealProperty = "amount_in_home_currency"
	DealPropertyAssist                                               DealProperty = "assist"
	DealPropertyCategory                                             DealProperty = "category"
	DealPropertyDaysToClose                                          DealProperty = "days_to_close"
	DealPropertyDealCurrencyCode                                     DealProperty = "deal_currency_code"
	DealPropertyAcv                                                  DealProperty = "hs_acv"
	DealPropertyAllAssignedBusinessUnitIds                           DealProperty = "hs_all_assigned_business_unit_ids"
	DealPropertyAnalyticsSource                                      DealProperty = "hs_analytics_source"
	DealPropertyAnalyticsSourceData1                                 DealProperty = "hs_analytics_source_data_1"
	DealPropertyAnalyticsSourceData2                                 DealProperty = "hs_analytics_source_data_2"
	DealPropertyArr                                                  DealProperty = "hs_arr"
	DealPropertyCampaign                                             DealProperty = "hs_campaign"
	DealPropertyClosedAmount                                         DealProperty = "hs_closed_amount"
	DealPropertyClosedAmountInHomeCurrency                           DealProperty = "hs_closed_amount_in_home_currency"
	DealPropertyCreatedByUserId                                      DealProperty = "hs_created_by_user_id"
	DealPropertyDateEntered9320750                                   DealProperty = "hs_date_entered_9320750"
	DealPropertyDateEntered9320751                                   DealProperty = "hs_date_entered_9320751"
	DealPropertyDateEntered9320752                                   DealProperty = "hs_date_entered_9320752"
	DealPropertyDateEntered9320753                                   DealProperty = "hs_date_entered_9320753"
	DealPropertyDateEntered9320755                                   DealProperty = "hs_date_entered_9320755"
	DealPropertyDateEntered9321133                                   DealProperty = "hs_date_entered_9321133"
	DealPropertyDateEntered9323337                                   DealProperty = "hs_date_entered_9323337"
	DealPropertyDateEntered9895007                                   DealProperty = "hs_date_entered_9895007"
	DealPropertyDateEntered9895008                                   DealProperty = "hs_date_entered_9895008"
	DealPropertyDateEntered9895009                                   DealProperty = "hs_date_entered_9895009"
	DealPropertyDateEnteredAppointmentScheduled                      DealProperty = "hs_date_entered_appointmentscheduled"
	DealPropertyDateEnteredClosedLost                                DealProperty = "hs_date_entered_closedlost"
	DealPropertyDateEnteredClosedWon                                 DealProperty = "hs_date_entered_closedwon"
	DealPropertyDateEnteredContractSent                              DealProperty = "hs_date_entered_contractsent"
	DealPropertyDateEnteredDecisionMakerBoughtIn                     DealProperty = "hs_date_entered_decisionmakerboughtin"
	DealPropertyDateEnteredPresentationScheduled                     DealProperty = "hs_date_entered_presentationscheduled"
	DealPropertyDateEnteredQualifiedToBuy                            DealProperty = "hs_date_entered_qualifiedtobuy"
	DealPropertyDateExited9320750                                    DealProperty = "hs_date_exited_9320750"
	DealPropertyDateExited9320751                                    DealProperty = "hs_date_exited_9320751"
	DealPropertyDateExited9320752                                    DealProperty = "hs_date_exited_9320752"
	DealPropertyDateExited9320753                                    DealProperty = "hs_date_exited_9320753"
	DealPropertyDateExited9320755                                    DealProperty = "hs_date_exited_9320755"
	DealPropertyDateExited9321133                                    DealProperty = "hs_date_exited_9321133"
	DealPropertyDateExited9323337                                    DealProperty = "hs_date_exited_9323337"
	DealPropertyDateExited9895007                                    DealProperty = "hs_date_exited_9895007"
	DealPropertyDateExited9895008                                    DealProperty = "hs_date_exited_9895008"
	DealPropertyDateExited9895009                                    DealProperty = "hs_date_exited_9895009"
	DealPropertyDateExitedAppointmentScheduled                       DealProperty = "hs_date_exited_appointmentscheduled"
	DealPropertyDateExitedClosedLost                                 DealProperty = "hs_date_exited_closedlost"
	DealPropertyDateExitedClosedWon                                  DealProperty = "hs_date_exited_closedwon"
	DealPropertyDateExitedContractSent                               DealProperty = "hs_date_exited_contractsent"
	DealPropertyDateExitedDecisionMakerBoughtIn                      DealProperty = "hs_date_exited_decisionmakerboughtin"
	DealPropertyDateExitedPresentationScheduled                      DealProperty = "hs_date_exited_presentationscheduled"
	DealPropertyDateExitedQualifiedToBuy                             DealProperty = "hs_date_exited_qualifiedtobuy"
	DealPropertyDealAmountCalculationPreference                      DealProperty = "hs_deal_amount_calculation_preference"
	DealPropertyDealStageProbability                                 DealProperty = "hs_deal_stage_probability"
	DealPropertyForecastAmount                                       DealProperty = "hs_forecast_amount"
	DealPropertyForecastProbability                                  DealProperty = "hs_forecast_probability"
	DealPropertyIsClosed                                             DealProperty = "hs_is_closed"
	DealPropertyLastModifiedDate                                     DealProperty = "hs_lastmodifieddate"
	DealPropertyLikelihoodToClose                                    DealProperty = "hs_likelihood_to_close"
	DealPropertyLineItemGlobalTermHsDiscountPercentage               DealProperty = "hs_line_item_global_term_hs_discount_percentage"
	DealPropertyLineItemGlobalTermHsDiscountPercentageEnabled        DealProperty = "hs_line_item_global_term_hs_discount_percentage_enabled"
	DealPropertyLineItemGlobalTermHsRecurringBillingPeriod           DealProperty = "hs_line_item_global_term_hs_recurring_billing_period"
	DealPropertyLineItemGlobalTermHsRecurringBillingPeriodEnabled    DealProperty = "hs_line_item_global_term_hs_recurring_billing_period_enabled"
	DealPropertyLineItemGlobalTermHsRecurringBillingStartDate        DealProperty = "hs_line_item_global_term_hs_recurring_billing_start_date"
	DealPropertyLineItemGlobalTermHsRecurringBillingStartDateEnabled DealProperty = "hs_line_item_global_term_hs_recurring_billing_start_date_enabled"
	DealPropertyLineItemGlobalTermRecurringBillingfrequency          DealProperty = "hs_line_item_global_term_recurringbillingfrequency"
	DealPropertyLineItemGlobalTermRecurringBillingfrequencyEnabled   DealProperty = "hs_line_item_global_term_recurringbillingfrequency_enabled"
	DealPropertyManualForecastCategory                               DealProperty = "hs_manual_forecast_category"
	DealPropertyMergedObjectIds                                      DealProperty = "hs_merged_object_ids"
	DealPropertyMrr                                                  DealProperty = "hs_mrr"
	DealPropertyNextStep                                             DealProperty = "hs_next_step"
	DealPropertyObjectId                                             DealProperty = "hs_object_id"
	DealPropertyPredictedAmount                                      DealProperty = "hs_predicted_amount"
	DealPropertyPredictedAmountInHomeCurrency                        DealProperty = "hs_predicted_amount_in_home_currency"
	DealPropertyProjectedAmount                                      DealProperty = "hs_projected_amount"
	DealPropertyProjectedAmountInHomeCurrency                        DealProperty = "hs_projected_amount_in_home_currency"
	DealPropertyTcv                                                  DealProperty = "hs_tcv"
	DealPropertyTimeIn9320750                                        DealProperty = "hs_time_in_9320750"
	DealPropertyTimeIn9320751                                        DealProperty = "hs_time_in_9320751"
	DealPropertyTimeIn9320752                                        DealProperty = "hs_time_in_9320752"
	DealPropertyTimeIn9320753                                        DealProperty = "hs_time_in_9320753"
	DealPropertyTimeIn9320755                                        DealProperty = "hs_time_in_9320755"
	DealPropertyTimeIn9321133                                        DealProperty = "hs_time_in_9321133"
	DealPropertyTimeIn9323337                                        DealProperty = "hs_time_in_9323337"
	DealPropertyTimeIn9895007                                        DealProperty = "hs_time_in_9895007"
	DealPropertyTimeIn9895008                                        DealProperty = "hs_time_in_9895008"
	DealPropertyTimeIn9895009                                        DealProperty = "hs_time_in_9895009"
	DealPropertyTimeInAppointmentScheduled                           DealProperty = "hs_time_in_appointmentscheduled"
	DealPropertyTimeInClosedLost                                     DealProperty = "hs_time_in_closedlost"
	DealPropertyTimeInClosedWon                                      DealProperty = "hs_time_in_closedwon"
	DealPropertyTimeInContractSent                                   DealProperty = "hs_time_in_contractsent"
	DealPropertyTimeInDecisionMakerBoughtIn                          DealProperty = "hs_time_in_decisionmakerboughtin"
	DealPropertyTimeInPresentationScheduled                          DealProperty = "hs_time_in_presentationscheduled"
	DealPropertyTimeInQualifiedToBuy                                 DealProperty = "hs_time_in_qualifiedtobuy"
	DealPropertyUpdatedByUserId                                      DealProperty = "hs_updated_by_user_id"
	DealPropertyUserIdsOfAllOwners                                   DealProperty = "hs_user_ids_of_all_owners"
	DealPropertyOwnerAssigneddate                                    DealProperty = "hubspot_owner_assigneddate"
	DealPropertyInsightlyOpportunityId                               DealProperty = "insightly_opportunity_id"
	DealPropertyLeadSourceBookedMeetings                             DealProperty = "lead_source___booked_meetings"
	DealPropertyLostDealReasons                                      DealProperty = "lost_deal_reasons"
	DealPropertyMeetingBookedBy                                      DealProperty = "meeting_booked_by"
	DealPropertyVidyardVideoSent                                     DealProperty = "vidyard_video_sent"
	DealPropertyDealname                                             DealProperty = "dealname"
	DealPropertyAmount                                               DealProperty = "amount"
	DealPropertyDealstage                                            DealProperty = "dealstage"
	DealPropertyPipeline                                             DealProperty = "pipeline"
	DealPropertyCloseDate                                            DealProperty = "closedate"
	DealPropertyCreateDate                                           DealProperty = "createdate"
	DealPropertyEngagementsLastMeetingBooked                         DealProperty = "engagements_last_meeting_booked"
	DealPropertyEngagementsLastMeetingBookedCampaign                 DealProperty = "engagements_last_meeting_booked_campaign"
	DealPropertyEngagementsLastMeetingBookedMedium                   DealProperty = "engagements_last_meeting_booked_medium"
	DealPropertyEngagementsLastMeetingBookedSource                   DealProperty = "engagements_last_meeting_booked_source"
	DealPropertyLatestMeetingActivity                                DealProperty = "hs_latest_meeting_activity"
	DealPropertySalesEmailLastReplied                                DealProperty = "hs_sales_email_last_replied"
	DealPropertyOwnerId                                              DealProperty = "hubspot_owner_id"
	DealPropertyNotesLastContacted                                   DealProperty = "notes_last_contacted"
	DealPropertyNotesLastUpdated                                     DealProperty = "notes_last_updated"
	DealPropertyNotesNextActivityDate                                DealProperty = "notes_next_activity_date"
	DealPropertyNumContactedNotes                                    DealProperty = "num_contacted_notes"
	DealPropertyNumNotes                                             DealProperty = "num_notes"
	DealPropertyHubspotCreatedate                                    DealProperty = "hs_createdate"
	DealPropertyTeamId                                               DealProperty = "hubspot_team_id"
	DealPropertyDealtype                                             DealProperty = "dealtype"
	DealPropertyAllOwnerIds                                          DealProperty = "hs_all_owner_ids"
	DealPropertyDescription                                          DealProperty = "description"
	DealPropertyAllTeamIds                                           DealProperty = "hs_all_team_ids"
	DealPropertyAllAccessibleTeamIds                                 DealProperty = "hs_all_accessible_team_ids"
	DealPropertyNumAssociatedContacts                                DealProperty = "num_associated_contacts"
	DealPropertyClosedLostReason                                     DealProperty = "closed_lost_reason"
	DealPropertyClosedWonReason                                      DealProperty = "closed_won_reason"
)

type GetDealsConfig struct {
	Limit            *uint
	After            *string
	Properties       *[]DealProperty
	CustomProperties *[]string
	Associations     *[]ObjectType
	Archived         *bool
}

// GetDeals returns all deals
//
func (service *Service) GetDeals(config *GetDealsConfig) (*[]Deal, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/deals"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		_properties := []string{}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				for _, p := range *config.Properties {
					_properties = append(_properties, string(p))
				}
			}
		}
		if config.CustomProperties != nil {
			if len(*config.CustomProperties) > 0 {
				_properties = append(_properties, *config.CustomProperties...)
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
		if config.Archived != nil {
			values.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}
	}

	after := ""
	if config.After != nil {
		after = *config.After
	}

	deals := []Deal{}

	for {
		dealsResponse := DealsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &dealsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, d := range dealsResponse.Results {
			deal_, e := getDeal(&d, config.CustomProperties)
			if e != nil {
				return nil, e
			}
			deals = append(deals, *deal_)
		}

		if config.After != nil { // explicit after parameter requested
			break
		}

		if dealsResponse.Paging == nil {
			break
		}

		if dealsResponse.Paging.Next.After == "" {
			break
		}

		after = dealsResponse.Paging.Next.After
	}

	return &deals, nil
}

func getDeal(deal *deal, customProperties *[]string) (*Deal, *errortools.Error) {
	deal_ := Deal{
		Id:               deal.Id,
		CreatedAt:        deal.CreatedAt,
		UpdatedAt:        deal.UpdatedAt,
		Archived:         deal.Archived,
		Associations:     deal.Associations,
		CustomProperties: make(map[string]string),
	}
	if deal.Properties != nil {
		p := DealProperties{}
		err := json.Unmarshal(deal.Properties, &p)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
		deal_.Properties = p
	}

	if customProperties != nil {
		p1 := make(map[string]string)
		err := json.Unmarshal(deal.Properties, &p1)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}

		for _, cp := range *customProperties {
			value, ok := p1[cp]
			if ok {
				deal_.CustomProperties[cp] = value
			}
		}
	}

	return &deal_, nil
}
