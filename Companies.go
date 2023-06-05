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
type company struct {
	Id           string                     `json:"id"`
	Properties   json.RawMessage            `json:"properties"`
	CreatedAt    h_types.DateTimeMSString   `json:"createdAt"`
	UpdatedAt    h_types.DateTimeMSString   `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}

type Company struct {
	Id               string
	Properties       CompanyProperties
	CustomProperties map[string]string
	CreatedAt        h_types.DateTimeMSString
	UpdatedAt        h_types.DateTimeMSString
	Archived         bool
	Associations     map[string]AssociationsSet
}

type CompanyProperties struct {
	Name                    *string                   `json:"name,omitempty"`
	Phone                   *string                   `json:"phone,omitempty"`
	Email                   *string                   `json:"email,omitempty"`
	DomainName              *string                   `json:"domain,omitempty"`
	LifecycleStage          *string                   `json:"lifecyclestage,omitempty"`
	Industry                *string                   `json:"industry,omitempty"`
	AnnualRevenue           *go_types.Int64String     `json:"annualrevenue,omitempty"`
	Founded                 *string                   `json:"founded_year,omitempty"`
	AboutUs                 *string                   `json:"about_us,omitempty"`
	Street                  *string                   `json:"address,omitempty"`
	Street2                 *string                   `json:"address2,omitempty"`
	City                    *string                   `json:"city,omitempty"`
	State                   *string                   `json:"state,omitempty"`
	Zip                     *string                   `json:"zip,omitempty"`
	Country                 *string                   `json:"country,omitempty"`
	NumberOfPageViews       *go_types.Int64String     `json:"hs_analytics_num_page_views,omitempty"`
	LastContacted           *h_types.DateTimeMSString `json:"notes_last_contacted,omitempty"`
	NumberOfTimesContacted  *go_types.Int64String     `json:"num_contacted_notes,omitempty"`
	OriginalSourceType      *string                   `json:"hs_analytics_source,omitempty"`
	NextActivityDate        *h_types.DateTimeString   `json:"notes_next_activity_date,omitempty"`
	LinkedinCompanyPage     *string                   `json:"linkedin_company_page,omitempty"`
	FacebookCompanyPage     *string                   `json:"facebook_company_page,omitempty"`
	TwitterHandle           *string                   `json:"twitterhandle,omitempty"`
	NumberOfFormSubmissions *go_types.Int64String     `json:"num_conversion_events,omitempty"`
	WebsiteUrl              *string                   `json:"website,omitempty"`
	OwnerId                 *string                   `json:"hubspot_owner_id,omitempty"`
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
	CompanyPropertyAllAssignedBusinessUnitIds                                  CompanyProperty = "hs_all_assigned_business_unit_ids"
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
	CompanyPropertyCreatedByUserId                                             CompanyProperty = "hs_created_by_user_id"
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
	CompanyPropertyObjectId                                                    CompanyProperty = "hs_object_id"
	CompanyPropertyPredictiveContactscoreV2                                    CompanyProperty = "hs_predictivecontactscore_v2"
	CompanyPropertyPredictiveContactscoreV2NextMaxMax                          CompanyProperty = "hs_predictivecontactscore_v2_next_max_max_d4e58c1e"
	CompanyPropertyTargetAccount                                               CompanyProperty = "hs_target_account"
	CompanyPropertyTargetAccountProbability                                    CompanyProperty = "hs_target_account_probability"
	CompanyPropertyTargetAccountRecommendationSnoozeTime                       CompanyProperty = "hs_target_account_recommendation_snooze_time"
	CompanyPropertyTargetAccountRecommendationState                            CompanyProperty = "hs_target_account_recommendation_state"
	CompanyPropertyTotalDealValue                                              CompanyProperty = "hs_total_deal_value"
	CompanyPropertyUpdatedByUserId                                             CompanyProperty = "hs_updated_by_user_id"
	CompanyPropertyUserIdsOfAllOwners                                          CompanyProperty = "hs_user_ids_of_all_owners"
	CompanyPropertyOwnerAssigneddate                                           CompanyProperty = "hubspot_owner_assigneddate"
	CompanyPropertyInsightlyCompanyId                                          CompanyProperty = "insightly_company_id"
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
	CompanyPropertyOwnerId                                                     CompanyProperty = "hubspot_owner_id"
	CompanyPropertyNotesLastContacted                                          CompanyProperty = "notes_last_contacted"
	CompanyPropertyNotesLastUpdated                                            CompanyProperty = "notes_last_updated"
	CompanyPropertyNotesNextActivityDate                                       CompanyProperty = "notes_next_activity_date"
	CompanyPropertyNumContactedNotes                                           CompanyProperty = "num_contacted_notes"
	CompanyPropertyNumNotes                                                    CompanyProperty = "num_notes"
	CompanyPropertyZIP                                                         CompanyProperty = "zip"
	CompanyPropertyCountry                                                     CompanyProperty = "country"
	CompanyPropertyTeamId                                                      CompanyProperty = "hubspot_team_id"
	CompanyPropertyAllOwnerIds                                                 CompanyProperty = "hs_all_owner_ids"
	CompanyPropertyWebsite                                                     CompanyProperty = "website"
	CompanyPropertyDomain                                                      CompanyProperty = "domain"
	CompanyPropertyAllTeamIds                                                  CompanyProperty = "hs_all_team_ids"
	CompanyPropertyAllAccessibleTeamIds                                        CompanyProperty = "hs_all_accessible_team_ids"
	CompanyPropertyNumberofemployees                                           CompanyProperty = "numberofemployees"
	CompanyPropertyIndustry                                                    CompanyProperty = "industry"
	CompanyPropertyAnnualRevenue                                               CompanyProperty = "annualrevenue"
	CompanyPropertyLifecyclestage                                              CompanyProperty = "lifecyclestage"
	CompanyPropertyLeadStatus                                                  CompanyProperty = "hs_lead_status"
	CompanyPropertyParentCompanyId                                             CompanyProperty = "hs_parent_company_id"
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
	Limit        *uint
	After        *string
	Properties   *[]string
	Associations *[]string
	Archived     *bool
}

// GetCompanies returns all companies
func (service *Service) GetCompanies(config *GetCompaniesConfig) (*[]Company, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/companies"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				values.Set("properties", strings.Join(*config.Properties, ","))
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

	companies := []Company{}

	for {
		companiesResponse := CompaniesResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &companiesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, c := range companiesResponse.Results {
			company_, e := getCompany(&c, config.Properties)
			if e != nil {
				return nil, e
			}
			/*company_ := Company{
				Id:               c.Id,
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
			}*/

			companies = append(companies, *company_)
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

type CreateCompanyConfig struct {
	Properties       CompanyProperties
	CustomProperties map[string]string
}

func (service *Service) CreateCompany(config *CreateCompanyConfig) (*Company, *errortools.Error) {
	endpoint := "objects/companies"
	company := Company{}

	body, e := companyPropertiesBody(config.Properties, config.CustomProperties)
	if e != nil {
		return nil, e
	}

	properties := struct {
		Properties map[string]string `json:"properties"`
	}{
		body,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     properties,
		ResponseModel: &company,
	}

	_, _, e = service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &company, nil
}

func companyPropertiesBody(properties CompanyProperties, customProperties map[string]string) (map[string]string, *errortools.Error) {
	// marshal
	b, err := json.Marshal(properties)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	// unmarshal to map
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	if customProperties == nil {
		return m, nil
	}
	if len(customProperties) == 0 {
		return m, nil
	}

	// append custom properties to map
	for key, value := range customProperties {
		if _, ok := m[key]; !ok {
			m[key] = value
		}
	}

	return m, nil
}

type UpdateCompanyConfig struct {
	CompanyId        string
	Properties       CompanyProperties
	CustomProperties map[string]string
}

func (service *Service) UpdateCompany(config *UpdateCompanyConfig) (*Company, *errortools.Error) {
	endpoint := "objects/companies"
	company := Company{}

	body, e := companyPropertiesBody(config.Properties, config.CustomProperties)
	if e != nil {
		return nil, e
	}

	properties := struct {
		Properties map[string]string `json:"properties"`
	}{
		body,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.CompanyId)),
		BodyModel:     properties,
		ResponseModel: &company,
	}

	_, _, e = service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &company, nil
}

type GetCompanyConfig struct {
	CompanyId        string
	Properties       *[]CompanyProperty
	CustomProperties *[]string
	Associations     *[]string
}

// GetCompany returns a specific company
func (service *Service) GetCompany(config *GetCompanyConfig) (*Company, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/companies"

	if config == nil {
		return nil, errortools.ErrorMessage("config is nil")
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

	company := company{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s?%s", endpoint, config.CompanyId, values.Encode())),
		ResponseModel: &company,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return getCompany(&company, config.CustomProperties)
}

func getCompany(company *company, customProperties *[]string) (*Company, *errortools.Error) {
	company_ := Company{
		Id:               company.Id,
		CreatedAt:        company.CreatedAt,
		UpdatedAt:        company.UpdatedAt,
		Archived:         company.Archived,
		Associations:     company.Associations,
		CustomProperties: make(map[string]string),
	}
	if company.Properties != nil {
		p := CompanyProperties{}
		err := json.Unmarshal(company.Properties, &p)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
		company_.Properties = p
	}

	if customProperties != nil {
		p1 := make(map[string]string)
		err := json.Unmarshal(company.Properties, &p1)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}

		for _, cp := range *customProperties {
			value, ok := p1[cp]
			if ok {
				company_.CustomProperties[cp] = value
			}
		}
	}

	return &company_, nil
}

type SearchCompanyConfig struct {
	Limit        *uint          `json:"limit,omitempty"`
	After        *string        `json:"after,omitempty"`
	FilterGroups *[]FilterGroup `json:"filterGroups,omitempty"`
	Sorts        *[]string      `json:"sorts,omitempty"`
	Query        *string        `json:"query,omitempty"`
	Properties   *[]string      `json:"properties,omitempty"`
}

type FilterGroup struct {
	Filters *[]filter `json:"filters"`
}

func (fg *FilterGroup) AddPropertyFilter(operator string, property string, value string, highValue string) {
	if fg.Filters == nil {
		fg.Filters = &[]filter{}
	}

	*fg.Filters = append(*fg.Filters, filter{
		Operator:     operator,
		PropertyName: property,
		Value:        value,
		HighValue:    highValue,
		isCustom:     false,
	})
}

func (fg *FilterGroup) AddCustomPropertyFilter(operator string, propertyName string, value string, highValue string) {
	if fg.Filters == nil {
		fg.Filters = &[]filter{}
	}

	*fg.Filters = append(*fg.Filters, filter{
		Operator:     operator,
		PropertyName: propertyName,
		Value:        value,
		HighValue:    highValue,
		isCustom:     true,
	})
}

type filter struct {
	Operator     string `json:"operator"`
	PropertyName string `json:"propertyName,omitempty"`
	Value        string `json:"value"`
	HighValue    string `json:"highValue,omitempty"`
	isCustom     bool   `json:"-"`
}

// SearchCompany returns a specific company
func (service *Service) SearchCompany(config *SearchCompanyConfig) (*[]Company, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("Config is nil")
	}

	endpoint := "objects/companies/search"

	companiesResponse := CompaniesResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &companiesResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	customProperties := []string{}
	if config.FilterGroups != nil {
		for _, filterGroup := range *config.FilterGroups {
			for _, filter := range *filterGroup.Filters {
				if filter.isCustom {
					customProperties = append(customProperties, filter.PropertyName)
				}
			}
		}
	}

	after := config.After

	companies := []Company{}

	for {
		companiesResponse := CompaniesResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm("objects/companies/search"),
			BodyModel:     config,
			ResponseModel: &companiesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, c := range companiesResponse.Results {
			company_, e := getCompany(&c, config.Properties)
			if e != nil {
				return nil, e
			}

			companies = append(companies, *company_)
		}

		if after != nil { // explicit after parameter requested
			break
		}

		if companiesResponse.Paging == nil {
			break
		}

		if companiesResponse.Paging.Next.After == "" {
			break
		}

		config.After = &companiesResponse.Paging.Next.After
	}

	return &companies, nil
}

func (service *Service) DeleteCompany(companyId string) *errortools.Error {
	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.urlCrm(fmt.Sprintf("objects/companies/%s", companyId)),
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}

func (service *Service) BatchDeleteCompanies(companyIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(companyIds) > index {
		if len(companyIds) > index+maxItemsPerBatch {
			e := service.batchDeleteCompanies(companyIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteCompanies(companyIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteCompanies(companyIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, companyId := range companyIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{companyId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/companies/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
