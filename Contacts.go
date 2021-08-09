package hubspot

import (
	"fmt"
	"net/url"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	go_types "github.com/leapforce-libraries/go_types"
)

type ContactsResponse struct {
	Results []Contact `json:"results"`
	Paging  *Paging   `json:"paging"`
}

// Contact stores Contact from Service
//
type Contact struct {
	ID           string                   `json:"id"`
	Properties   ContactProperties        `json:"properties"`
	CreatedAt    h_types.DateTimeMSString `json:"createdAt"`
	UpdatedAt    h_types.DateTimeMSString `json:"updatedAt"`
	Archived     bool                     `json:"archived"`
	Associations *Associations            `json:"associations"`
}

type ContactProperties struct {
	FirstName                   *string                   `json:"firstname"`
	LastName                    *string                   `json:"lastname"`
	JobTitle                    *string                   `json:"jobtitle"`
	Email                       *string                   `json:"email"`
	MobilePhone                 *string                   `json:"mobilephone"`
	Phone                       *string                   `json:"phone"`
	AveragePageviews            *go_types.Int64String     `json:"hs_analytics_average_page_views"`
	OriginalSource              *string                   `json:"hs_analytics_source"`
	OwnerID                     *go_types.Int64String     `json:"hubspot_owner_id"`
	CampaignOfLastBooking       *string                   `json:"engagements_last_meeting_booked_campaign"`
	CurrentlyInSequence         *go_types.BoolString      `json:"hs_sequences_is_enrolled"`
	DateOfLastMeetingBooked     *h_types.DateTimeMSString `json:"engagements_last_meeting_booked"`
	FirstConversionEventName    *string                   `json:"first_conversion_event_name"`
	FirstConversionDate         *h_types.DateTimeMSString `json:"first_conversion_date"`
	LastActivityDate            *h_types.DateTimeMSString `json:"notes_last_updated"`
	LastContacted               *h_types.DateTimeMSString `json:"notes_last_contacted"`
	LastEngagementDate          *h_types.DateTimeMSString `json:"hs_last_sales_activity_timestamp"`
	MediumOfLastBooking         *string                   `json:"engagements_last_meeting_booked_medium"`
	MembershipNotes             *string                   `json:"hs_content_membership_notes"`
	Message                     *string                   `json:"message"`
	NextActivityDate            *h_types.DateTimeMSString `json:"notes_next_activity_date"`
	NumberOfSalesActivities     *go_types.Int64String     `json:"num_notes"`
	NumberOfTimesContacted      *go_types.Int64String     `json:"num_contacted_notes"`
	RecentSalesEmailClickedDate *h_types.DateTimeMSString `json:"hs_sales_email_last_clicked"`
	RecentSalesEmailOpenedDate  *h_types.DateTimeMSString `json:"hs_sales_email_last_opened"`
	RecentSalesEmailRepliedDate *h_types.DateTimeMSString `json:"hs_sales_email_last_replied"`
	SourceOfLastBooking         *string                   `json:"engagements_last_meeting_booked_source"`
	Status                      *string                   `json:"hs_content_membership_status"`
}

type ContactProperty string

const (
	ContactPropertyCompanySize                                    ContactProperty = "company_size"
	ContactPropertyDateOfBirth                                    ContactProperty = "date_of_birth"
	ContactPropertyDatesToRemember                                ContactProperty = "dates_to_remember"
	ContactPropertyDaysToClose                                    ContactProperty = "days_to_close"
	ContactPropertyDegree                                         ContactProperty = "degree"
	ContactPropertyDescription                                    ContactProperty = "description"
	ContactPropertyFacebook                                       ContactProperty = "facebook"
	ContactPropertyFieldOfStudy                                   ContactProperty = "field_of_study"
	ContactPropertyFirstConversionDate                            ContactProperty = "first_conversion_date"
	ContactPropertyFirstConversionEventName                       ContactProperty = "first_conversion_event_name"
	ContactPropertyFirstDealCreatedDate                           ContactProperty = "first_deal_created_date"
	ContactPropertyGender                                         ContactProperty = "gender"
	ContactPropertyGraduationDate                                 ContactProperty = "graduation_date"
	ContactPropertyAdditionalEmails                               ContactProperty = "hs_additional_emails"
	ContactPropertyAllAssignedBusinessUnitIDs                     ContactProperty = "hs_all_assigned_business_unit_ids"
	ContactPropertyAllContactVids                                 ContactProperty = "hs_all_contact_vids"
	ContactPropertyAnalyticsFirstTouchConvertingCampaign          ContactProperty = "hs_analytics_first_touch_converting_campaign"
	ContactPropertyAnalyticsLastTouchConvertingCampaign           ContactProperty = "hs_analytics_last_touch_converting_campaign"
	ContactPropertyAvatarFilemanagerKey                           ContactProperty = "hs_avatar_filemanager_key"
	ContactPropertyBuyingRole                                     ContactProperty = "hs_buying_role"
	ContactPropertyCalculatedFormSubmissions                      ContactProperty = "hs_calculated_form_submissions"
	ContactPropertyCalculatedMergedVids                           ContactProperty = "hs_calculated_merged_vids"
	ContactPropertyCalculatedMobileNumber                         ContactProperty = "hs_calculated_mobile_number"
	ContactPropertyCalculatedPhoneNumber                          ContactProperty = "hs_calculated_phone_number"
	ContactPropertyCalculatedPhoneNumberAreaCode                  ContactProperty = "hs_calculated_phone_number_area_code"
	ContactPropertyCalculatedPhoneNumberCountryCode               ContactProperty = "hs_calculated_phone_number_country_code"
	ContactPropertyCalculatedPhoneNumberRegionCode                ContactProperty = "hs_calculated_phone_number_region_code"
	ContactPropertyContentMembershipEmailConfirmed                ContactProperty = "hs_content_membership_email_confirmed"
	ContactPropertyContentMembershipNotes                         ContactProperty = "hs_content_membership_notes"
	ContactPropertyContentMembershipRegisteredAt                  ContactProperty = "hs_content_membership_registered_at"
	ContactPropertyContentMembershipRegistrationDomainSentTo      ContactProperty = "hs_content_membership_registration_domain_sent_to"
	ContactPropertyContentMembershipRegistrationEmailSentAt       ContactProperty = "hs_content_membership_registration_email_sent_at"
	ContactPropertyContentMembershipStatus                        ContactProperty = "hs_content_membership_status"
	ContactPropertyConversationsVisitorEmail                      ContactProperty = "hs_conversations_visitor_email"
	ContactPropertyCountIsUnworked                                ContactProperty = "hs_count_is_unworked"
	ContactPropertyCountIsWorked                                  ContactProperty = "hs_count_is_worked"
	ContactPropertyCreatedByConversations                         ContactProperty = "hs_created_by_conversations"
	ContactPropertyCreatedByUserID                                ContactProperty = "hs_created_by_user_id"
	ContactPropertyHubspotCreatedate                              ContactProperty = "hs_createdate"
	ContactPropertyDocumentLastRevisited                          ContactProperty = "hs_document_last_revisited"
	ContactPropertyEmailBadAddress                                ContactProperty = "hs_email_bad_address"
	ContactPropertyEmailCustomerQuarantinedReason                 ContactProperty = "hs_email_customer_quarantined_reason"
	ContactPropertyEmailDomain                                    ContactProperty = "hs_email_domain"
	ContactPropertyEmailHardBounceReason                          ContactProperty = "hs_email_hard_bounce_reason"
	ContactPropertyEmailHardBounceReasonEnum                      ContactProperty = "hs_email_hard_bounce_reason_enum"
	ContactPropertyEmailQuarantined                               ContactProperty = "hs_email_quarantined"
	ContactPropertyEmailQuarantinedReason                         ContactProperty = "hs_email_quarantined_reason"
	ContactPropertyEmailRecipientFatigueRecoveryTime              ContactProperty = "hs_email_recipient_fatigue_recovery_time"
	ContactPropertyEmailSendsSinceLastEngagement                  ContactProperty = "hs_email_sends_since_last_engagement"
	ContactPropertyEmailconfirmationstatus                        ContactProperty = "hs_emailconfirmationstatus"
	ContactPropertyFacebookAdClicked                              ContactProperty = "hs_facebook_ad_clicked"
	ContactPropertyFacebookClickID                                ContactProperty = "hs_facebook_click_id"
	ContactPropertyFacebookID                                     ContactProperty = "hs_facebookid"
	ContactPropertyFeedbackLastNpsFollowUp                        ContactProperty = "hs_feedback_last_nps_follow_up"
	ContactPropertyFeedbackLastNpsRating                          ContactProperty = "hs_feedback_last_nps_rating"
	ContactPropertyFeedbackLastSurveyDate                         ContactProperty = "hs_feedback_last_survey_date"
	ContactPropertyFeedbackShowNpsWebSurvey                       ContactProperty = "hs_feedback_show_nps_web_survey"
	ContactPropertyFirstEngagementObjectID                        ContactProperty = "hs_first_engagement_object_id"
	ContactPropertyGoogleClickID                                  ContactProperty = "hs_google_click_id"
	ContactPropertyGoogleplusid                                   ContactProperty = "hs_googleplusid"
	ContactPropertyIPTimezone                                     ContactProperty = "hs_ip_timezone"
	ContactPropertyIsContact                                      ContactProperty = "hs_is_contact"
	ContactPropertyIsUnworked                                     ContactProperty = "hs_is_unworked"
	ContactPropertyLastSalesActivityDate                          ContactProperty = "hs_last_sales_activity_date"
	ContactPropertyLastSalesActivityTimestamp                     ContactProperty = "hs_last_sales_activity_timestamp"
	ContactPropertyHubspotLastModifiedDate                        ContactProperty = "hs_lastmodifieddate"
	ContactPropertyLatestSequenceEndedDate                        ContactProperty = "hs_latest_sequence_ended_date"
	ContactPropertyLatestSequenceEnrolled                         ContactProperty = "hs_latest_sequence_enrolled"
	ContactPropertyLatestSequenceEnrolledDate                     ContactProperty = "hs_latest_sequence_enrolled_date"
	ContactPropertyLatestSequenceFinishedDate                     ContactProperty = "hs_latest_sequence_finished_date"
	ContactPropertyLatestSequenceUnenrolledDate                   ContactProperty = "hs_latest_sequence_unenrolled_date"
	ContactPropertyLeadStatus                                     ContactProperty = "hs_lead_status"
	ContactPropertyLegalBasis                                     ContactProperty = "hs_legal_basis"
	ContactPropertyLinkedinID                                     ContactProperty = "hs_linkedinid"
	ContactPropertyMarketableReasonID                             ContactProperty = "hs_marketable_reason_id"
	ContactPropertyMarketableReasonType                           ContactProperty = "hs_marketable_reason_type"
	ContactPropertyMarketableStatus                               ContactProperty = "hs_marketable_status"
	ContactPropertyMarketableUntilRenewal                         ContactProperty = "hs_marketable_until_renewal"
	ContactPropertyMergedObjectIDs                                ContactProperty = "hs_merged_object_ids"
	ContactPropertyObjectID                                       ContactProperty = "hs_object_id"
	ContactPropertyPredictivecontactscoreV2                       ContactProperty = "hs_predictivecontactscore_v2"
	ContactPropertyPredictiveScoringTier                          ContactProperty = "hs_predictivescoringtier"
	ContactPropertySaFirstEngagementDate                          ContactProperty = "hs_sa_first_engagement_date"
	ContactPropertySaFirstEngagementDescr                         ContactProperty = "hs_sa_first_engagement_descr"
	ContactPropertySaFirstEngagementObjectType                    ContactProperty = "hs_sa_first_engagement_object_type"
	ContactPropertySalesEmailLastClicked                          ContactProperty = "hs_sales_email_last_clicked"
	ContactPropertySalesEmailLastOpened                           ContactProperty = "hs_sales_email_last_opened"
	ContactPropertySearchableCalculatedInternationalMobileNumber  ContactProperty = "hs_searchable_calculated_international_mobile_number"
	ContactPropertySearchableCalculatedInternationalPhoneNumber   ContactProperty = "hs_searchable_calculated_international_phone_number"
	ContactPropertySearchableCalculatedMobileNumber               ContactProperty = "hs_searchable_calculated_mobile_number"
	ContactPropertySearchableCalculatedPhoneNumber                ContactProperty = "hs_searchable_calculated_phone_number"
	ContactPropertySequencesEnrolledCount                         ContactProperty = "hs_sequences_enrolled_count"
	ContactPropertySequencesIsEnrolled                            ContactProperty = "hs_sequences_is_enrolled"
	ContactPropertyTestpurge                                      ContactProperty = "hs_testpurge"
	ContactPropertyTestrollback                                   ContactProperty = "hs_testrollback"
	ContactPropertyTimeBetweenContactCreationAndDealClose         ContactProperty = "hs_time_between_contact_creation_and_deal_close"
	ContactPropertyTimeBetweenContactCreationAndDealCreation      ContactProperty = "hs_time_between_contact_creation_and_deal_creation"
	ContactPropertyTimeToFirstEngagement                          ContactProperty = "hs_time_to_first_engagement"
	ContactPropertyTimeToMoveFromLeadToCustomer                   ContactProperty = "hs_time_to_move_from_lead_to_customer"
	ContactPropertyTimeToMoveFromMarketingqualifiedleadToCustomer ContactProperty = "hs_time_to_move_from_marketingqualifiedlead_to_customer"
	ContactPropertyTimeToMoveFromOpportunityToCustomer            ContactProperty = "hs_time_to_move_from_opportunity_to_customer"
	ContactPropertyTimeToMoveFromSalesqualifiedleadToCustomer     ContactProperty = "hs_time_to_move_from_salesqualifiedlead_to_customer"
	ContactPropertyTimeToMoveFromSubscriberToCustomer             ContactProperty = "hs_time_to_move_from_subscriber_to_customer"
	ContactPropertyTwitterid                                      ContactProperty = "hs_twitterid"
	ContactPropertyUpdatedByUserID                                ContactProperty = "hs_updated_by_user_id"
	ContactPropertyUserIDsOfAllOwners                             ContactProperty = "hs_user_ids_of_all_owners"
	ContactPropertyOwnerAssigneddate                              ContactProperty = "hubspot_owner_assigneddate"
	ContactPropertyInsightlyContactID                             ContactProperty = "insightly_contact_id"
	ContactPropertyIPCity                                         ContactProperty = "ip_city"
	ContactPropertyIPCountry                                      ContactProperty = "ip_country"
	ContactPropertyIPCountryCode                                  ContactProperty = "ip_country_code"
	ContactPropertyIPLatLon                                       ContactProperty = "ip_latlon"
	ContactPropertyIPState                                        ContactProperty = "ip_state"
	ContactPropertyIPStateCode                                    ContactProperty = "ip_state_code"
	ContactPropertyIPZIPCode                                      ContactProperty = "ip_zipcode"
	ContactPropertyJobFunction                                    ContactProperty = "job_function"
	ContactPropertyLastmodifieddate                               ContactProperty = "lastmodifieddate"
	ContactPropertyLinkedin                                       ContactProperty = "linkedin"
	ContactPropertyMaritalStatus                                  ContactProperty = "marital_status"
	ContactPropertyMilitaryStatus                                 ContactProperty = "military_status"
	ContactPropertyNeverbouncevalidationresult                    ContactProperty = "neverbouncevalidationresult"
	ContactPropertyNumAssociatedDeals                             ContactProperty = "num_associated_deals"
	ContactPropertyNumConversionEvents                            ContactProperty = "num_conversion_events"
	ContactPropertyNumUniqueConversionEvents                      ContactProperty = "num_unique_conversion_events"
	ContactPropertyRecentConversionDate                           ContactProperty = "recent_conversion_date"
	ContactPropertyRecentConversionEventName                      ContactProperty = "recent_conversion_event_name"
	ContactPropertyRecentDealAmount                               ContactProperty = "recent_deal_amount"
	ContactPropertyRecentDealCloseDate                            ContactProperty = "recent_deal_close_date"
	ContactPropertyRecordID                                       ContactProperty = "record_id"
	ContactPropertyRelationshipStatus                             ContactProperty = "relationship_status"
	ContactPropertySchool                                         ContactProperty = "school"
	ContactPropertySeniority                                      ContactProperty = "seniority"
	ContactPropertyStartDate                                      ContactProperty = "start_date"
	ContactPropertyTest                                           ContactProperty = "test"
	ContactPropertyTestAmber                                      ContactProperty = "test_amber"
	ContactPropertyTotalRevenue                                   ContactProperty = "total_revenue"
	ContactPropertyTwitter                                        ContactProperty = "twitter"
	ContactPropertyTyepformAnswers                                ContactProperty = "tyepform_answers"
	ContactPropertyTypeformNameResult                             ContactProperty = "typeform_name_result"
	ContactPropertyWorkEmail                                      ContactProperty = "work_email"
	ContactPropertyZoomWebinarAttendanceAverageDuration           ContactProperty = "zoom_webinar_attendance_average_duration"
	ContactPropertyZoomWebinarAttendanceCount                     ContactProperty = "zoom_webinar_attendance_count"
	ContactPropertyZoomWebinarJoinlink                            ContactProperty = "zoom_webinar_joinlink"
	ContactPropertyZoomWebinarRegistrationCount                   ContactProperty = "zoom_webinar_registration_count"
	ContactPropertyFirstName                                      ContactProperty = "firstname"
	ContactPropertyAnalyticsFirstURL                              ContactProperty = "hs_analytics_first_url"
	ContactPropertyEmailDelivered                                 ContactProperty = "hs_email_delivered"
	ContactPropertyEmailOptout8575803                             ContactProperty = "hs_email_optout_8575803"
	ContactPropertyEmailOptout8780041                             ContactProperty = "hs_email_optout_8780041"
	ContactPropertyTwitterHandle                                  ContactProperty = "twitterhandle"
	ContactPropertyCurrentlyInWorkflow                            ContactProperty = "currentlyinworkflow"
	ContactPropertyFollowercount                                  ContactProperty = "followercount"
	ContactPropertyAnalyticsLastURL                               ContactProperty = "hs_analytics_last_url"
	ContactPropertyEmailOpen                                      ContactProperty = "hs_email_open"
	ContactPropertyLastName                                       ContactProperty = "lastname"
	ContactPropertyPartnerRegistrationDomain                      ContactProperty = "partner_registration_domain"
	ContactPropertyPartnerRegistrationMessage                     ContactProperty = "partner_registration_message"
	ContactPropertyPartnerRegistrationStatus                      ContactProperty = "partner_registration_status"
	ContactPropertyAnalyticsNumPageViews                          ContactProperty = "hs_analytics_num_page_views"
	ContactPropertyEmailClick                                     ContactProperty = "hs_email_click"
	ContactPropertyPartnerRegistrationAction                      ContactProperty = "partner_registration_action"
	ContactPropertyPartnerRegistrationExpiryDate                  ContactProperty = "partner_registration_expiry_date"
	ContactPropertySalutation                                     ContactProperty = "salutation"
	ContactPropertyTwitterProfilePhoto                            ContactProperty = "twitterprofilephoto"
	ContactPropertyEmail                                          ContactProperty = "email"
	ContactPropertyAnalyticsNumVisits                             ContactProperty = "hs_analytics_num_visits"
	ContactPropertyEmailBounce                                    ContactProperty = "hs_email_bounce"
	ContactPropertyPersona                                        ContactProperty = "hs_persona"
	ContactPropertySocialLastEngagement                           ContactProperty = "hs_social_last_engagement"
	ContactPropertyAnalyticsNumEventCompletions                   ContactProperty = "hs_analytics_num_event_completions"
	ContactPropertyEmailOptOut                                    ContactProperty = "hs_email_optout"
	ContactPropertySocialTwitterClicks                            ContactProperty = "hs_social_twitter_clicks"
	ContactPropertyMobilephone                                    ContactProperty = "mobilephone"
	ContactPropertyPhone                                          ContactProperty = "phone"
	ContactPropertyFax                                            ContactProperty = "fax"
	ContactPropertyAnalyticsFirstTimestamp                        ContactProperty = "hs_analytics_first_timestamp"
	ContactPropertyEmailLastEmailName                             ContactProperty = "hs_email_last_email_name"
	ContactPropertyEmailLastSendDate                              ContactProperty = "hs_email_last_send_date"
	ContactPropertySocialFacebookClicks                           ContactProperty = "hs_social_facebook_clicks"
	ContactPropertyAddress                                        ContactProperty = "address"
	ContactPropertyEngagementsLastMeetingBooked                   ContactProperty = "engagements_last_meeting_booked"
	ContactPropertyEngagementsLastMeetingBookedCampaign           ContactProperty = "engagements_last_meeting_booked_campaign"
	ContactPropertyEngagementsLastMeetingBookedMedium             ContactProperty = "engagements_last_meeting_booked_medium"
	ContactPropertyEngagementsLastMeetingBookedSource             ContactProperty = "engagements_last_meeting_booked_source"
	ContactPropertyAnalyticsFirstVisitTimestamp                   ContactProperty = "hs_analytics_first_visit_timestamp"
	ContactPropertyEmailLastOpenDate                              ContactProperty = "hs_email_last_open_date"
	ContactPropertyLatestMeetingActivity                          ContactProperty = "hs_latest_meeting_activity"
	ContactPropertySalesEmailLastReplied                          ContactProperty = "hs_sales_email_last_replied"
	ContactPropertySocialLinkedinClicks                           ContactProperty = "hs_social_linkedin_clicks"
	ContactPropertyOwnerID                                        ContactProperty = "hubspot_owner_id"
	ContactPropertyNotesLastContacted                             ContactProperty = "notes_last_contacted"
	ContactPropertyNotesLastUpdated                               ContactProperty = "notes_last_updated"
	ContactPropertyNotesNextActivityDate                          ContactProperty = "notes_next_activity_date"
	ContactPropertyNumContactedNotes                              ContactProperty = "num_contacted_notes"
	ContactPropertyNumNotes                                       ContactProperty = "num_notes"
	ContactPropertyOwneremail                                     ContactProperty = "owneremail"
	ContactPropertyOwnername                                      ContactProperty = "ownername"
	ContactPropertySurveymonkeyeventlastupdated                   ContactProperty = "surveymonkeyeventlastupdated"
	ContactPropertyWebinareventlastupdated                        ContactProperty = "webinareventlastupdated"
	ContactPropertyCity                                           ContactProperty = "city"
	ContactPropertyAnalyticsLastTimestamp                         ContactProperty = "hs_analytics_last_timestamp"
	ContactPropertyEmailLastClickDate                             ContactProperty = "hs_email_last_click_date"
	ContactPropertySocialGooglePlusClicks                         ContactProperty = "hs_social_google_plus_clicks"
	ContactPropertyTeamID                                         ContactProperty = "hubspot_team_id"
	ContactPropertyLinkedinBio                                    ContactProperty = "linkedinbio"
	ContactPropertyTwitterBio                                     ContactProperty = "twitterbio"
	ContactPropertyAllOwnerIDs                                    ContactProperty = "hs_all_owner_ids"
	ContactPropertyAnalyticsLastVisitTimestamp                    ContactProperty = "hs_analytics_last_visit_timestamp"
	ContactPropertyEmailFirstSendDate                             ContactProperty = "hs_email_first_send_date"
	ContactPropertySocialNumBroadcastClicks                       ContactProperty = "hs_social_num_broadcast_clicks"
	ContactPropertyState                                          ContactProperty = "state"
	ContactPropertyAllTeamIDs                                     ContactProperty = "hs_all_team_ids"
	ContactPropertyAnalyticsSource                                ContactProperty = "hs_analytics_source"
	ContactPropertyEmailFirstOpenDate                             ContactProperty = "hs_email_first_open_date"
	ContactPropertyZIP                                            ContactProperty = "zip"
	ContactPropertyCountry                                        ContactProperty = "country"
	ContactPropertyAllAccessibleTeamIDs                           ContactProperty = "hs_all_accessible_team_ids"
	ContactPropertyAnalyticsSourceData1                           ContactProperty = "hs_analytics_source_data_1"
	ContactPropertyEmailFirstClickDate                            ContactProperty = "hs_email_first_click_date"
	ContactPropertyLinkedinconnections                            ContactProperty = "linkedinconnections"
	ContactPropertyAnalyticsSourceData2                           ContactProperty = "hs_analytics_source_data_2"
	ContactPropertyEmailIsIneligible                              ContactProperty = "hs_email_is_ineligible"
	ContactPropertyLanguage                                       ContactProperty = "hs_language"
	ContactPropertyKloutscoregeneral                              ContactProperty = "kloutscoregeneral"
	ContactPropertyAnalyticsFirstReferrer                         ContactProperty = "hs_analytics_first_referrer"
	ContactPropertyEmailFirstReplyDate                            ContactProperty = "hs_email_first_reply_date"
	ContactPropertyJobTitle                                       ContactProperty = "jobtitle"
	ContactPropertyPhoto                                          ContactProperty = "photo"
	ContactPropertyAnalyticsLastReferrer                          ContactProperty = "hs_analytics_last_referrer"
	ContactPropertyEmailLastReplyDate                             ContactProperty = "hs_email_last_reply_date"
	ContactPropertyMessage                                        ContactProperty = "message"
	ContactPropertyClosedate                                      ContactProperty = "closedate"
	ContactPropertyAnalyticsAveragePageViews                      ContactProperty = "hs_analytics_average_page_views"
	ContactPropertyEmailReplied                                   ContactProperty = "hs_email_replied"
	ContactPropertyAnalyticsRevenue                               ContactProperty = "hs_analytics_revenue"
	ContactPropertyLifecyclestageLeadDate                         ContactProperty = "hs_lifecyclestage_lead_date"
	ContactPropertyLifecyclestageMarketingqualifiedleadDate       ContactProperty = "hs_lifecyclestage_marketingqualifiedlead_date"
	ContactPropertyLifecyclestageOpportunityDate                  ContactProperty = "hs_lifecyclestage_opportunity_date"
	ContactPropertyLifecyclestage                                 ContactProperty = "lifecyclestage"
	ContactPropertyLifecyclestageSalesqualifiedleadDate           ContactProperty = "hs_lifecyclestage_salesqualifiedlead_date"
	ContactPropertyCreatedate                                     ContactProperty = "createdate"
	ContactPropertyLifecyclestageEvangelistDate                   ContactProperty = "hs_lifecyclestage_evangelist_date"
	ContactPropertyLifecyclestageCustomerDate                     ContactProperty = "hs_lifecyclestage_customer_date"
	ContactPropertyHubspotscore                                   ContactProperty = "hubspotscore"
	ContactPropertyCompany                                        ContactProperty = "company"
	ContactPropertyLifecyclestageSubscriberDate                   ContactProperty = "hs_lifecyclestage_subscriber_date"
	ContactPropertyLifecyclestageOtherDate                        ContactProperty = "hs_lifecyclestage_other_date"
	ContactPropertyWebsite                                        ContactProperty = "website"
	ContactPropertyNumemployees                                   ContactProperty = "numemployees"
	ContactPropertyAnnualrevenue                                  ContactProperty = "annualrevenue"
	ContactPropertyIndustry                                       ContactProperty = "industry"
	ContactPropertyAssociatedcompanyid                            ContactProperty = "associatedcompanyid"
	ContactPropertyAssociatedcompanylastupdated                   ContactProperty = "associatedcompanylastupdated"
	ContactPropertyPredictivecontactscorebucket                   ContactProperty = "hs_predictivecontactscorebucket"
	ContactPropertyPredictivecontactscore                         ContactProperty = "hs_predictivecontactscore"
)

type GetContactsConfig struct {
	Limit        *uint
	After        *string
	Properties   *[]ContactProperty
	Associations *[]ObjectType
	Archived     *bool
}

// GetContacts returns all contacts
//
func (service *Service) GetContacts(config *GetContactsConfig) (*[]Contact, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/contacts"

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

	contacts := []Contact{}

	for {
		contactsResponse := ContactsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			URL:           service.url(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &contactsResponse,
		}

		_, _, e := service.get(&requestConfig)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, contactsResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if contactsResponse.Paging == nil {
			break
		}

		if contactsResponse.Paging.Next.After == "" {
			break
		}

		after = contactsResponse.Paging.Next.After
	}

	return &contacts, nil
}
