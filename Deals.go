package hubspot

import (
	"fmt"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type DealsResponse struct {
	Results []Deal  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Deal stores Deal from Service
//
type Deal struct {
	ID         string         `json:"id"`
	Properties DealProperties `json:"properties"`
	CreatedAt  string         `json:"createdAt"`
	UpdatedAt  string         `json:"updatedAt"`
	Archived   bool           `json:"archived"`
}

type DealProperties struct {
	Amount              *int     `json:"amount"`
	Assist              *int     `json:"assist"`
	Category            *string  `json:"category"`
	CloseDate           *string  `json:"closedate"`
	CreateDate          *string  `json:"createdate"`
	DealName            *string  `json:"dealname"`
	DealStage           *string  `json:"dealstage"`
	ForecastAmount      *float64 `json:"hs_forecast_amount"`
	ForecastProbability *float64 `json:"hs_forecast_probability"`
	LastUpdated         *string  `json:"notes_last_updated"`
	OwnerID             *string  `json:"hubspot_owner_id"`
}

type DealProperty string

const (
	DealPropertyAmountInHomeCurrency                                   DealProperty = "amount_in_home_currency"
	DealPropertyAssist                                                 DealProperty = "assist"
	DealPropertyCategory                                               DealProperty = "category"
	DealPropertyDaysToClose                                            DealProperty = "days_to_close"
	DealPropertyDealCurrencyCode                                       DealProperty = "deal_currency_code"
	DealPropertyHsAcv                                                  DealProperty = "hs_acv"
	DealPropertyHsAllAssignedBusinessUnitIds                           DealProperty = "hs_all_assigned_business_unit_ids"
	DealPropertyHsAnalyticsSource                                      DealProperty = "hs_analytics_source"
	DealPropertyHsAnalyticsSourceData1                                 DealProperty = "hs_analytics_source_data_1"
	DealPropertyHsAnalyticsSourceData2                                 DealProperty = "hs_analytics_source_data_2"
	DealPropertyHsArr                                                  DealProperty = "hs_arr"
	DealPropertyHsCampaign                                             DealProperty = "hs_campaign"
	DealPropertyHsClosedAmount                                         DealProperty = "hs_closed_amount"
	DealPropertyHsClosedAmountInHomeCurrency                           DealProperty = "hs_closed_amount_in_home_currency"
	DealPropertyHsCreatedByUserID                                      DealProperty = "hs_created_by_user_id"
	DealPropertyHsDateEntered9320750                                   DealProperty = "hs_date_entered_9320750"
	DealPropertyHsDateEntered9320751                                   DealProperty = "hs_date_entered_9320751"
	DealPropertyHsDateEntered9320752                                   DealProperty = "hs_date_entered_9320752"
	DealPropertyHsDateEntered9320753                                   DealProperty = "hs_date_entered_9320753"
	DealPropertyHsDateEntered9320755                                   DealProperty = "hs_date_entered_9320755"
	DealPropertyHsDateEntered9321133                                   DealProperty = "hs_date_entered_9321133"
	DealPropertyHsDateEntered9323337                                   DealProperty = "hs_date_entered_9323337"
	DealPropertyHsDateEntered9895007                                   DealProperty = "hs_date_entered_9895007"
	DealPropertyHsDateEntered9895008                                   DealProperty = "hs_date_entered_9895008"
	DealPropertyHsDateEntered9895009                                   DealProperty = "hs_date_entered_9895009"
	DealPropertyHsDateEnteredAppointmentscheduled                      DealProperty = "hs_date_entered_appointmentscheduled"
	DealPropertyHsDateEnteredClosedlost                                DealProperty = "hs_date_entered_closedlost"
	DealPropertyHsDateEnteredClosedwon                                 DealProperty = "hs_date_entered_closedwon"
	DealPropertyHsDateEnteredContractsent                              DealProperty = "hs_date_entered_contractsent"
	DealPropertyHsDateEnteredDecisionmakerboughtin                     DealProperty = "hs_date_entered_decisionmakerboughtin"
	DealPropertyHsDateEnteredPresentationscheduled                     DealProperty = "hs_date_entered_presentationscheduled"
	DealPropertyHsDateEnteredQualifiedtobuy                            DealProperty = "hs_date_entered_qualifiedtobuy"
	DealPropertyHsDateExited9320750                                    DealProperty = "hs_date_exited_9320750"
	DealPropertyHsDateExited9320751                                    DealProperty = "hs_date_exited_9320751"
	DealPropertyHsDateExited9320752                                    DealProperty = "hs_date_exited_9320752"
	DealPropertyHsDateExited9320753                                    DealProperty = "hs_date_exited_9320753"
	DealPropertyHsDateExited9320755                                    DealProperty = "hs_date_exited_9320755"
	DealPropertyHsDateExited9321133                                    DealProperty = "hs_date_exited_9321133"
	DealPropertyHsDateExited9323337                                    DealProperty = "hs_date_exited_9323337"
	DealPropertyHsDateExited9895007                                    DealProperty = "hs_date_exited_9895007"
	DealPropertyHsDateExited9895008                                    DealProperty = "hs_date_exited_9895008"
	DealPropertyHsDateExited9895009                                    DealProperty = "hs_date_exited_9895009"
	DealPropertyHsDateExitedAppointmentscheduled                       DealProperty = "hs_date_exited_appointmentscheduled"
	DealPropertyHsDateExitedClosedlost                                 DealProperty = "hs_date_exited_closedlost"
	DealPropertyHsDateExitedClosedwon                                  DealProperty = "hs_date_exited_closedwon"
	DealPropertyHsDateExitedContractsent                               DealProperty = "hs_date_exited_contractsent"
	DealPropertyHsDateExitedDecisionmakerboughtin                      DealProperty = "hs_date_exited_decisionmakerboughtin"
	DealPropertyHsDateExitedPresentationscheduled                      DealProperty = "hs_date_exited_presentationscheduled"
	DealPropertyHsDateExitedQualifiedtobuy                             DealProperty = "hs_date_exited_qualifiedtobuy"
	DealPropertyHsDealAmountCalculationPreference                      DealProperty = "hs_deal_amount_calculation_preference"
	DealPropertyHsDealStageProbability                                 DealProperty = "hs_deal_stage_probability"
	DealPropertyHsForecastAmount                                       DealProperty = "hs_forecast_amount"
	DealPropertyHsForecastProbability                                  DealProperty = "hs_forecast_probability"
	DealPropertyHsIsClosed                                             DealProperty = "hs_is_closed"
	DealPropertyHsLastmodifieddate                                     DealProperty = "hs_lastmodifieddate"
	DealPropertyHsLikelihoodToClose                                    DealProperty = "hs_likelihood_to_close"
	DealPropertyHsLineItemGlobalTermHsDiscountPercentage               DealProperty = "hs_line_item_global_term_hs_discount_percentage"
	DealPropertyHsLineItemGlobalTermHsDiscountPercentageEnabled        DealProperty = "hs_line_item_global_term_hs_discount_percentage_enabled"
	DealPropertyHsLineItemGlobalTermHsRecurringBillingPeriod           DealProperty = "hs_line_item_global_term_hs_recurring_billing_period"
	DealPropertyHsLineItemGlobalTermHsRecurringBillingPeriodEnabled    DealProperty = "hs_line_item_global_term_hs_recurring_billing_period_enabled"
	DealPropertyHsLineItemGlobalTermHsRecurringBillingStartDate        DealProperty = "hs_line_item_global_term_hs_recurring_billing_start_date"
	DealPropertyHsLineItemGlobalTermHsRecurringBillingStartDateEnabled DealProperty = "hs_line_item_global_term_hs_recurring_billing_start_date_enabled"
	DealPropertyHsLineItemGlobalTermRecurringbillingfrequency          DealProperty = "hs_line_item_global_term_recurringbillingfrequency"
	DealPropertyHsLineItemGlobalTermRecurringbillingfrequencyEnabled   DealProperty = "hs_line_item_global_term_recurringbillingfrequency_enabled"
	DealPropertyHsManualForecastCategory                               DealProperty = "hs_manual_forecast_category"
	DealPropertyHsMergedObjectIds                                      DealProperty = "hs_merged_object_ids"
	DealPropertyHsMrr                                                  DealProperty = "hs_mrr"
	DealPropertyHsNextStep                                             DealProperty = "hs_next_step"
	DealPropertyHsObjectID                                             DealProperty = "hs_object_id"
	DealPropertyHsPredictedAmount                                      DealProperty = "hs_predicted_amount"
	DealPropertyHsPredictedAmountInHomeCurrency                        DealProperty = "hs_predicted_amount_in_home_currency"
	DealPropertyHsProjectedAmount                                      DealProperty = "hs_projected_amount"
	DealPropertyHsProjectedAmountInHomeCurrency                        DealProperty = "hs_projected_amount_in_home_currency"
	DealPropertyHsTcv                                                  DealProperty = "hs_tcv"
	DealPropertyHsTimeIn9320750                                        DealProperty = "hs_time_in_9320750"
	DealPropertyHsTimeIn9320751                                        DealProperty = "hs_time_in_9320751"
	DealPropertyHsTimeIn9320752                                        DealProperty = "hs_time_in_9320752"
	DealPropertyHsTimeIn9320753                                        DealProperty = "hs_time_in_9320753"
	DealPropertyHsTimeIn9320755                                        DealProperty = "hs_time_in_9320755"
	DealPropertyHsTimeIn9321133                                        DealProperty = "hs_time_in_9321133"
	DealPropertyHsTimeIn9323337                                        DealProperty = "hs_time_in_9323337"
	DealPropertyHsTimeIn9895007                                        DealProperty = "hs_time_in_9895007"
	DealPropertyHsTimeIn9895008                                        DealProperty = "hs_time_in_9895008"
	DealPropertyHsTimeIn9895009                                        DealProperty = "hs_time_in_9895009"
	DealPropertyHsTimeInAppointmentscheduled                           DealProperty = "hs_time_in_appointmentscheduled"
	DealPropertyHsTimeInClosedlost                                     DealProperty = "hs_time_in_closedlost"
	DealPropertyHsTimeInClosedwon                                      DealProperty = "hs_time_in_closedwon"
	DealPropertyHsTimeInContractsent                                   DealProperty = "hs_time_in_contractsent"
	DealPropertyHsTimeInDecisionmakerboughtin                          DealProperty = "hs_time_in_decisionmakerboughtin"
	DealPropertyHsTimeInPresentationscheduled                          DealProperty = "hs_time_in_presentationscheduled"
	DealPropertyHsTimeInQualifiedtobuy                                 DealProperty = "hs_time_in_qualifiedtobuy"
	DealPropertyHsUpdatedByUserID                                      DealProperty = "hs_updated_by_user_id"
	DealPropertyHsUserIdsOfAllOwners                                   DealProperty = "hs_user_ids_of_all_owners"
	DealPropertyHubspotOwnerAssigneddate                               DealProperty = "hubspot_owner_assigneddate"
	DealPropertyInsightlyOpportunityID                                 DealProperty = "insightly_opportunity_id"
	DealPropertyLeadSourceBookedMeetings                               DealProperty = "lead_source___booked_meetings"
	DealPropertyLostDealReasons                                        DealProperty = "lost_deal_reasons"
	DealPropertyMeetingBookedBy                                        DealProperty = "meeting_booked_by"
	DealPropertyVidyardVideoSent                                       DealProperty = "vidyard_video_sent"
	DealPropertyDealname                                               DealProperty = "dealname"
	DealPropertyAmount                                                 DealProperty = "amount"
	DealPropertyDealstage                                              DealProperty = "dealstage"
	DealPropertyPipeline                                               DealProperty = "pipeline"
	DealPropertyClosedate                                              DealProperty = "closedate"
	DealPropertyCreatedate                                             DealProperty = "createdate"
	DealPropertyEngagementsLastMeetingBooked                           DealProperty = "engagements_last_meeting_booked"
	DealPropertyEngagementsLastMeetingBookedCampaign                   DealProperty = "engagements_last_meeting_booked_campaign"
	DealPropertyEngagementsLastMeetingBookedMedium                     DealProperty = "engagements_last_meeting_booked_medium"
	DealPropertyEngagementsLastMeetingBookedSource                     DealProperty = "engagements_last_meeting_booked_source"
	DealPropertyHsLatestMeetingActivity                                DealProperty = "hs_latest_meeting_activity"
	DealPropertyHsSalesEmailLastReplied                                DealProperty = "hs_sales_email_last_replied"
	DealPropertyHubspotOwnerID                                         DealProperty = "hubspot_owner_id"
	DealPropertyNotesLastContacted                                     DealProperty = "notes_last_contacted"
	DealPropertyNotesLastUpdated                                       DealProperty = "notes_last_updated"
	DealPropertyNotesNextActivityDate                                  DealProperty = "notes_next_activity_date"
	DealPropertyNumContactedNotes                                      DealProperty = "num_contacted_notes"
	DealPropertyNumNotes                                               DealProperty = "num_notes"
	DealPropertyHsCreatedate                                           DealProperty = "hs_createdate"
	DealPropertyHubspotTeamID                                          DealProperty = "hubspot_team_id"
	DealPropertyDealtype                                               DealProperty = "dealtype"
	DealPropertyHsAllOwnerIds                                          DealProperty = "hs_all_owner_ids"
	DealPropertyDescription                                            DealProperty = "description"
	DealPropertyHsAllTeamIds                                           DealProperty = "hs_all_team_ids"
	DealPropertyHsAllAccessibleTeamIDs                                 DealProperty = "hs_all_accessible_team_ids"
	DealPropertyNumAssociatedContacts                                  DealProperty = "num_associated_contacts"
	DealPropertyClosedLostReason                                       DealProperty = "closed_lost_reason"
	DealPropertyClosedWonReason                                        DealProperty = "closed_won_reason"
)

type GetDealsConfig struct {
	Limit        *uint
	After        *string
	Properties   *[]DealProperty
	Associations *[]ObjectType
	Archived     *bool
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
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				_properties := []string{}
				for _, p := range *config.Properties {
					_properties = append(_properties, string(p))
				}
				values.Set("properties", strings.Join(_properties, ","))
			}
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

	for true {
		dealsResponse := DealsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &dealsResponse,
		}
		fmt.Println(service.url(fmt.Sprintf("%s?%s", endpoint, values.Encode())))

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		deals = append(deals, dealsResponse.Results...)

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
