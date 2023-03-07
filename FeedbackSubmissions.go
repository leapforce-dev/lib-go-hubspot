package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
	"time"
)

type FeedbackSubmissionsResponse struct {
	Results []FeedbackSubmission `json:"results"`
	Paging  Paging               `json:"paging"`
}

type FeedbackSubmission struct {
	Id         string `json:"id"`
	Properties struct {
		AntwoordOpAanbeveling               string    `json:"antwoord_op_aanbeveling"`
		HoeWaarschijnlijk                   string    `json:"hoe_waarschijnlijk_is_het_dat_je_leapforce_zult_aanbevelen_aan_familie__vrienden_of_collega_s_"`
		AgentEmail                          string    `json:"hs_agent_email"`
		AgentResponsibleForTheConversation  string    `json:"hs_agent_id"`
		AgentName                           string    `json:"hs_agent_name"`
		AllAccessibleTeamIds                string    `json:"hs_all_accessible_team_ids"`
		AllAssignedBusinessUnitIds          string    `json:"hs_all_assigned_business_unit_ids"`
		AllOwnerIds                         string    `json:"hs_all_owner_ids"`
		AllTeamIds                          string    `json:"hs_all_team_ids"`
		ChatflowId                          string    `json:"hs_chatflow_id"`
		ChatflowName                        string    `json:"hs_chatflow_name"`
		ChatflowObjectId                    string    `json:"hs_chatflow_object_id"`
		ContactFirstname                    string    `json:"hs_contact_firstname"`
		ContactId                           string    `json:"hs_contact_id"`
		ContactLastname                     string    `json:"hs_contact_lastname"`
		ConversationThreadId                string    `json:"hs_conversation_thread_id"`
		CreatedByUserId                     string    `json:"hs_created_by_user_id"`
		ObjectCreateDate                    time.Time `json:"hs_createdate"`
		ObjectLastModifiedDate              time.Time `json:"hs_lastmodifieddate"`
		MergedRecordIds                     string    `json:"hs_merged_object_ids"`
		RecordId                            string    `json:"hs_object_id"`
		ReadOnlyObject                      bool      `json:"hs_read_only"`
		TagIds                              string    `json:"hs_tag_ids"`
		Tags                                string    `json:"hs_tags"`
		TicketId                            string    `json:"hs_ticket_id"`
		TicketOwnerAvatarUri                string    `json:"hs_ticket_owner_avatar_uri"`
		TicketOwner                         string    `json:"hs_ticket_owner_id"`
		TicketName                          string    `json:"hs_ticket_subject"`
		UniqueCreationKey                   string    `json:"hs_unique_creation_key"`
		UpdatedByUserId                     string    `json:"hs_updated_by_user_id"`
		UserIdsOfAllNotificationFollowers   string    `json:"hs_user_ids_of_all_notification_followers"`
		UserIdsOfAllNotificationUnfollowers string    `json:"hs_user_ids_of_all_notification_unfollowers"`
		UserIdsOfAllOwners                  string    `json:"hs_user_ids_of_all_owners"`
		OwnerAssignedDate                   time.Time `json:"hubspot_owner_assigneddate"`
		Owner                               string    `json:"hubspot_owner_id"`
		HubSpotTeam                         string    `json:"hubspot_team_id"`
		QuoteVoorOpDeWebsite                string    `json:"quote_voor_op_de_website"`
		SupportFeedback                     string    `json:"support_feedback"`
		SupportSurveyBeschrijving           string    `json:"support_survey_beschrijving"`
		IndustryStandardQuestionType        string    `json:"hs_industry_standard_question_type"`
		Sentiment                           string    `json:"hs_sentiment"`
		SurveyId                            string    `json:"hs_survey_id"`
		SurveyType                          string    `json:"hs_survey_type"`
		Source                              string    `json:"hs_survey_channel"`
		Date                                time.Time `json:"hs_submission_timestamp"`
		Rating                              string    `json:"hs_value"`
		FeedbackSentiment                   string    `json:"hs_response_group"`
		Response                            string    `json:"hs_content"`
		IngestionId                         string    `json:"hs_ingestion_id"`
		KnowledgeArticleId                  string    `json:"hs_knowledge_article_id"`
		VisitorId                           string    `json:"hs_visitor_id"`
		EngagementId                        string    `json:"hs_engagement_id"`
		SubmissionUrl                       string    `json:"hs_submission_url"`
		SurveyName                          string    `json:"hs_survey_name"`
		FormGuid                            string    `json:"hs_form_guid"`
		Email                               string    `json:"hs_contact_email_rollup"`
		SubmissionName                      string    `json:"hs_submission_name"`
	} `json:"properties"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Archived  bool      `json:"archived"`
}

type GetFeedbackSubmissionsConfig struct {
	Properties string
}

func (service *Service) GetFeedbackSubmissions(config *GetFeedbackSubmissionsConfig) (*[]FeedbackSubmission, *errortools.Error) {
	values := url.Values{}
	if config != nil {
		values.Set("properties", config.Properties)
	}

	var feedbackSubmissions []FeedbackSubmission

	var after = ""

	for {
		if after != "" {
			values.Set("after", after)
		}

		var feedbackSubmissionsResponse FeedbackSubmissionsResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("objects/feedback_submissions?%s", values.Encode())),
			ResponseModel: &feedbackSubmissionsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		feedbackSubmissions = append(feedbackSubmissions, feedbackSubmissionsResponse.Results...)

		if feedbackSubmissionsResponse.Paging.Next.After == "" {
			break
		}

		after = feedbackSubmissionsResponse.Paging.Next.After
	}

	return &feedbackSubmissions, nil
}
