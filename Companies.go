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

type CompaniesResponse struct {
	Results []company `json:"results"`
	Paging  *Paging   `json:"paging"`
}

// Company stores Company from Service
//
type company struct {
	ID           string                   `json:"id"`
	Properties   json.RawMessage          `json:"properties"`
	CreatedAt    h_types.DateTimeMSString `json:"createdAt"`
	UpdatedAt    h_types.DateTimeMSString `json:"updatedAt"`
	Archived     bool                     `json:"archived"`
	Associations *Associations            `json:"associations"`
}
type Company struct {
	ID               string
	Properties       CompanyProperties
	CustomProperties map[string]string
	CreatedAt        h_types.DateTimeMSString
	UpdatedAt        h_types.DateTimeMSString
	Archived         bool
	Associations     *Associations
}

type CompanyProperties struct {
	Name                    *string                   `json:"name"`
	Phone                   *string                   `json:"phone"`
	Email                   *string                   `json:"email"`
	DomainName              *string                   `json:"domain"`
	LifecycleStage          *string                   `json:"lifecyclestage"`
	Industry                *string                   `json:"industry"`
	AnnualRevenue           *go_types.Int64String     `json:"annualrevenue"`
	Founded                 *string                   `json:"founded_year"`
	AboutUs                 *string                   `json:"about_us"`
	Street                  *string                   `json:"address"`
	Street2                 *string                   `json:"address2"`
	City                    *string                   `json:"city"`
	State                   *string                   `json:"state"`
	ZIP                     *string                   `json:"zip"`
	Country                 *string                   `json:"country"`
	NumberOfPageViews       *go_types.Int64String     `json:"hs_analytics_num_page_views"`
	LastContacted           *h_types.DateTimeMSString `json:"notes_last_contacted"`
	NumberOfTimesContacted  *go_types.Int64String     `json:"num_contacted_notes"`
	OriginalSourceType      *string                   `json:"hs_analytics_source"`
	NextActivityDate        *h_types.DateTimeString   `json:"notes_next_activity_date"`
	LinkedinCompanyPage     *string                   `json:"linkedin_company_page"`
	FacebookCompanyPage     *string                   `json:"facebook_company_page"`
	NumberOfFormSubmissions *go_types.Int64String     `json:"num_conversion_events"`
	WebsiteURL              *string                   `json:"website"`
	OwnerID                 *string                   `json:"hubspot_owner_id"`
}

type CompanyProperty string

const (
	CompanyPropertyAboutUs                                                     CompanyProperty = "about_us"
	CompanyPropertyAliassen                                                    CompanyProperty = "aliassen"
	CompanyPropertyClientCodeLeapforce                                         CompanyProperty = "client_code_leapforce"
	CompanyPropertyClosedateTimestampEarliestValue                             CompanyProperty = "closedate_timestamp_earliest_value_a2a17e6e"
	CompanyPropertyEmail                                                       CompanyProperty = "email"
	CompanyPropertyFacebookFans                                                CompanyProperty = "facebookfans"
	CompanyPropertyFirstContactCreatedateTimestampEarliestValue                CompanyProperty = "first_contact_createdate_timestamp_earliest_value_78b50eea"
	CompanyPropertyFirstConversionDate                                         CompanyProperty = "first_conversion_date"
	CompanyPropertyFirstConversionDateTimestampEarliestValue                   CompanyProperty = "first_conversion_date_timestamp_earliest_value_61f58f2c"
	CompanyPropertyFirstConversionEventName                                    CompanyProperty = "first_conversion_event_name"
	CompanyPropertyFirstConversionEventNameTimestampEarliestValue              CompanyProperty = "first_conversion_event_name_timestamp_earliest_value_68ddae0a"
	CompanyPropertyFirstDealCreatedDate                                        CompanyProperty = "first_deal_created_date"
	CompanyPropertyFoundedYear                                                 CompanyProperty = "founded_year"
	CompanyPropertyGdsName                                                     CompanyProperty = "gds_name"
	CompanyPropertyAdditionalDomains                                           CompanyProperty = "hs_additional_domains"
	CompanyPropertyAllAssignedBusinessUnitIDs                                  CompanyProperty = "hs_all_assigned_business_unit_ids"
	CompanyPropertyAnalyticsFirstTimestamp                                     CompanyProperty = "hs_analytics_first_timestamp"
	CompanyPropertyAnalyticsFirstTimestampTimestampEarliestValue               CompanyProperty = "hs_analytics_first_timestamp_timestamp_earliest_value_11e3a63a"
	CompanyPropertyAnalyticsFirstTouchConvertingCampaign                       CompanyProperty = "hs_analytics_first_touch_converting_campaign"
	CompanyPropertyAnalyticsFirstTouchConvertingCampaignTimestampEarliestValue CompanyProperty = "hs_analytics_first_touch_converting_campaign_timestamp_earliest_value_4757fe10"
	CompanyPropertyAnalyticsFirstVisitTimestamp                                CompanyProperty = "hs_analytics_first_visit_timestamp"
	CompanyPropertyAnalyticsFirstVisitTimestampTimestampEarliestValueAcc       CompanyProperty = "hs_analytics_first_visit_timestamp_timestamp_earliest_value_accc17ae"
	CompanyPropertyAnalyticsLastTimestamp                                      CompanyProperty = "hs_analytics_last_timestamp"
	CompanyPropertyAnalyticsLastTimestampTimestampLatestValue                  CompanyProperty = "hs_analytics_last_timestamp_timestamp_latest_value_4e16365a"
	CompanyPropertyAnalyticsLastTouchConvertingCampaign                        CompanyProperty = "hs_analytics_last_touch_converting_campaign"
	CompanyPropertyAnalyticsLastTouchConvertingCampaignTimestampLatestValue    CompanyProperty = "hs_analytics_last_touch_converting_campaign_timestamp_latest_value_81a64e30"
	CompanyPropertyAnalyticsLastVisitTimestamp                                 CompanyProperty = "hs_analytics_last_visit_timestamp"
	CompanyPropertyAnalyticsLastVisitTimestampTimestampLatestValue             CompanyProperty = "hs_analytics_last_visit_timestamp_timestamp_latest_value_999a0fce"
	CompanyPropertyAnalyticsNumPageViews                                       CompanyProperty = "hs_analytics_num_page_views"
	CompanyPropertyAnalyticsNumPageViewsCardinalitySum                         CompanyProperty = "hs_analytics_num_page_views_cardinality_sum_e46e85b0"
	CompanyPropertyAnalyticsNumVisits                                          CompanyProperty = "hs_analytics_num_visits"
	CompanyPropertyAnalyticsNumVisitsCardinalitySum                            CompanyProperty = "hs_analytics_num_visits_cardinality_sum_53d952a6"
	CompanyPropertyAnalyticsSource                                             CompanyProperty = "hs_analytics_source"
	CompanyPropertyAnalyticsSourceData1                                        CompanyProperty = "hs_analytics_source_data_1"
	CompanyPropertyAnalyticsSourceData1TimestampEarliestValue                  CompanyProperty = "hs_analytics_source_data_1_timestamp_earliest_value_9b2f1fa1"
	CompanyPropertyAnalyticsSourceData2                                        CompanyProperty = "hs_analytics_source_data_2"
	CompanyPropertyAnalyticsSourceData2TimestampEarliestValue                  CompanyProperty = "hs_analytics_source_data_2_timestamp_earliest_value_9b2f9400"
	CompanyPropertyAnalyticsSourceTimestampEarliestValue                       CompanyProperty = "hs_analytics_source_timestamp_earliest_value_25a3a52c"
	CompanyPropertyAvatarFilemanagerKey                                        CompanyProperty = "hs_avatar_filemanager_key"
	CompanyPropertyCreatedByUserID                                             CompanyProperty = "hs_created_by_user_id"
	CompanyPropertyHubspotCreatedate                                           CompanyProperty = "hs_createdate"
	CompanyPropertyIdealCustomerProfile                                        CompanyProperty = "hs_ideal_customer_profile"
	CompanyPropertyIsTargetAccount                                             CompanyProperty = "hs_is_target_account"
	CompanyPropertyLastBookedMeetingDate                                       CompanyProperty = "hs_last_booked_meeting_date"
	CompanyPropertyLastLoggedCallDate                                          CompanyProperty = "hs_last_logged_call_date"
	CompanyPropertyLastOpenTaskDate                                            CompanyProperty = "hs_last_open_task_date"
	CompanyPropertyLastSalesActivityDate                                       CompanyProperty = "hs_last_sales_activity_date"
	CompanyPropertyLastSalesActivityTimestamp                                  CompanyProperty = "hs_last_sales_activity_timestamp"
	CompanyPropertyLastmodifieddate                                            CompanyProperty = "hs_lastmodifieddate"
	CompanyPropertyMergedObjectIds                                             CompanyProperty = "hs_merged_object_ids"
	CompanyPropertyNumBlockers                                                 CompanyProperty = "hs_num_blockers"
	CompanyPropertyNumContactsWithBuyingRoles                                  CompanyProperty = "hs_num_contacts_with_buying_roles"
	CompanyPropertyNumDecisionMakers                                           CompanyProperty = "hs_num_decision_makers"
	CompanyPropertyNumOpenDeals                                                CompanyProperty = "hs_num_open_deals"
	CompanyPropertyObjectID                                                    CompanyProperty = "hs_object_id"
	CompanyPropertyPredictiveContactscoreV2                                    CompanyProperty = "hs_predictivecontactscore_v2"
	CompanyPropertyPredictiveContactscoreV2NextMaxMax                          CompanyProperty = "hs_predictivecontactscore_v2_next_max_max_d4e58c1e"
	CompanyPropertyTargetAccount                                               CompanyProperty = "hs_target_account"
	CompanyPropertyTargetAccountProbability                                    CompanyProperty = "hs_target_account_probability"
	CompanyPropertyTargetAccountRecommendationSnoozeTime                       CompanyProperty = "hs_target_account_recommendation_snooze_time"
	CompanyPropertyTargetAccountRecommendationState                            CompanyProperty = "hs_target_account_recommendation_state"
	CompanyPropertyTotalDealValue                                              CompanyProperty = "hs_total_deal_value"
	CompanyPropertyUpdatedByUserID                                             CompanyProperty = "hs_updated_by_user_id"
	CompanyPropertyUserIdsOfAllOwners                                          CompanyProperty = "hs_user_ids_of_all_owners"
	CompanyPropertyOwnerAssigneddate                                           CompanyProperty = "hubspot_owner_assigneddate"
	CompanyPropertyInsightlyCompanyID                                          CompanyProperty = "insightly_company_id"
	CompanyPropertyIsPublic                                                    CompanyProperty = "is_public"
	CompanyPropertyNetwerken                                                   CompanyProperty = "netwerken"
	CompanyPropertyNumAssociatedContacts                                       CompanyProperty = "num_associated_contacts"
	CompanyPropertyNumAssociatedDeals                                          CompanyProperty = "num_associated_deals"
	CompanyPropertyNumConversionEvents                                         CompanyProperty = "num_conversion_events"
	CompanyPropertyNumConversionEventsCardinalitySum                           CompanyProperty = "num_conversion_events_cardinality_sum_d095f14b"
	CompanyPropertyRecentConversionDate                                        CompanyProperty = "recent_conversion_date"
	CompanyPropertyRecentConversionDateTimestampLatestValue                    CompanyProperty = "recent_conversion_date_timestamp_latest_value_72856da1"
	CompanyPropertyRecentConversionEventName                                   CompanyProperty = "recent_conversion_event_name"
	CompanyPropertyRecentConversionEventNameTimestampLatestValue               CompanyProperty = "recent_conversion_event_name_timestamp_latest_value_66c820bf"
	CompanyPropertyRecentDealAmount                                            CompanyProperty = "recent_deal_amount"
	CompanyPropertyRecentDealCloseDate                                         CompanyProperty = "recent_deal_close_date"
	CompanyPropertyRelationType                                                CompanyProperty = "relation_type"
	CompanyPropertySourceLeadinfo                                              CompanyProperty = "source_leadinfo"
	CompanyPropertySubscriptionStatus                                          CompanyProperty = "subscription_status"
	CompanyPropertyTestRelationType                                            CompanyProperty = "test_relation_type"
	CompanyPropertyTimezone                                                    CompanyProperty = "timezone"
	CompanyPropertyTotalMoneyRaised                                            CompanyProperty = "total_money_raised"
	CompanyPropertyTotalRevenue                                                CompanyProperty = "total_revenue"
	CompanyPropertyTwitter                                                     CompanyProperty = "twitter"
	CompanyPropertyName                                                        CompanyProperty = "name"
	CompanyPropertyOwnerEmail                                                  CompanyProperty = "owneremail"
	CompanyPropertyTwitterhandle                                               CompanyProperty = "twitterhandle"
	CompanyPropertyOwnername                                                   CompanyProperty = "ownername"
	CompanyPropertyPhone                                                       CompanyProperty = "phone"
	CompanyPropertyTwitterbio                                                  CompanyProperty = "twitterbio"
	CompanyPropertyTwitterfollowers                                            CompanyProperty = "twitterfollowers"
	CompanyPropertyAddress                                                     CompanyProperty = "address"
	CompanyPropertyAddress2                                                    CompanyProperty = "address2"
	CompanyPropertyFacebookCompanyPage                                         CompanyProperty = "facebook_company_page"
	CompanyPropertyCity                                                        CompanyProperty = "city"
	CompanyPropertyLinkedinCompanyPage                                         CompanyProperty = "linkedin_company_page"
	CompanyPropertyLinkedinbio                                                 CompanyProperty = "linkedinbio"
	CompanyPropertyState                                                       CompanyProperty = "state"
	CompanyPropertyGoogleplusPage                                              CompanyProperty = "googleplus_page"
	CompanyPropertyEngagementsLastMeetingBooked                                CompanyProperty = "engagements_last_meeting_booked"
	CompanyPropertyEngagementsLastMeetingBookedCampaign                        CompanyProperty = "engagements_last_meeting_booked_campaign"
	CompanyPropertyEngagementsLastMeetingBookedMedium                          CompanyProperty = "engagements_last_meeting_booked_medium"
	CompanyPropertyEngagementsLastMeetingBookedSource                          CompanyProperty = "engagements_last_meeting_booked_source"
	CompanyPropertyLatestMeetingActivity                                       CompanyProperty = "hs_latest_meeting_activity"
	CompanyPropertySalesEmailLastReplied                                       CompanyProperty = "hs_sales_email_last_replied"
	CompanyPropertyOwnerID                                                     CompanyProperty = "hubspot_owner_id"
	CompanyPropertyNotesLastContacted                                          CompanyProperty = "notes_last_contacted"
	CompanyPropertyNotesLastUpdated                                            CompanyProperty = "notes_last_updated"
	CompanyPropertyNotesNextActivityDate                                       CompanyProperty = "notes_next_activity_date"
	CompanyPropertyNumContactedNotes                                           CompanyProperty = "num_contacted_notes"
	CompanyPropertyNumNotes                                                    CompanyProperty = "num_notes"
	CompanyPropertyZIP                                                         CompanyProperty = "zip"
	CompanyPropertyCountry                                                     CompanyProperty = "country"
	CompanyPropertyTeamID                                                      CompanyProperty = "hubspot_team_id"
	CompanyPropertyAllOwnerIds                                                 CompanyProperty = "hs_all_owner_ids"
	CompanyPropertyWebsite                                                     CompanyProperty = "website"
	CompanyPropertyDomain                                                      CompanyProperty = "domain"
	CompanyPropertyAllTeamIDs                                                  CompanyProperty = "hs_all_team_ids"
	CompanyPropertyAllAccessibleTeamIDs                                        CompanyProperty = "hs_all_accessible_team_ids"
	CompanyPropertyNumberofemployees                                           CompanyProperty = "numberofemployees"
	CompanyPropertyIndustry                                                    CompanyProperty = "industry"
	CompanyPropertyAnnualRevenue                                               CompanyProperty = "annualrevenue"
	CompanyPropertyLifecyclestage                                              CompanyProperty = "lifecyclestage"
	CompanyPropertyLeadStatus                                                  CompanyProperty = "hs_lead_status"
	CompanyPropertyParentCompanyID                                             CompanyProperty = "hs_parent_company_id"
	CompanyPropertyType                                                        CompanyProperty = "type"
	CompanyPropertyDescription                                                 CompanyProperty = "description"
	CompanyPropertyNumChildCompanies                                           CompanyProperty = "hs_num_child_companies"
	CompanyPropertyHubspotScore                                                CompanyProperty = "hubspotscore"
	CompanyPropertyCreatedate                                                  CompanyProperty = "createdate"
	CompanyPropertyClosedate                                                   CompanyProperty = "closedate"
	CompanyPropertyFirstContactCreatedate                                      CompanyProperty = "first_contact_createdate"
	CompanyPropertyDaysToClose                                                 CompanyProperty = "days_to_close"
	CompanyPropertyWebTechnologies                                             CompanyProperty = "web_technologies"
)

type GetCompaniesConfig struct {
	Limit            *uint
	After            *string
	Properties       *[]CompanyProperty
	CustomProperties *[]string
	Associations     *[]ObjectType
	Archived         *bool
}

// GetCompanies returns all companies
//
func (service *Service) GetCompanies(config *GetCompaniesConfig) (*[]Company, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/companies"

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

	companies := []Company{}

	for {
		companiesResponse := CompaniesResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &companiesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, c := range companiesResponse.Results {
			company_ := Company{
				ID:               c.ID,
				CreatedAt:        c.CreatedAt,
				UpdatedAt:        c.UpdatedAt,
				Archived:         c.Archived,
				Associations:     c.Associations,
				CustomProperties: make(map[string]string),
			}
			if c.Properties == nil {
				continue
			}

			p := CompanyProperties{}
			err := json.Unmarshal(c.Properties, &p)
			if err != nil {
				return nil, errortools.ErrorMessage(err)
			}
			company_.Properties = p

			if config.CustomProperties != nil {
				p1 := make(map[string]string)
				err := json.Unmarshal(c.Properties, &p1)
				if err != nil {
					return nil, errortools.ErrorMessage(err)
				}

				for _, cp := range *config.CustomProperties {
					value, ok := p1[cp]
					if ok {
						company_.CustomProperties[cp] = value
					}
				}
			}

			companies = append(companies, company_)
		}

		if config.After != nil { // explicit after parameter requested
			break
		}

		if companiesResponse.Paging == nil {
			break
		}

		if companiesResponse.Paging.Next.After == "" {
			break
		}

		after = companiesResponse.Paging.Next.After
	}

	return &companies, nil
}
