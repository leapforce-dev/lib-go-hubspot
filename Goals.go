package hubspot

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	go_types "github.com/leapforce-libraries/go_types"
	"net/http"
	"net/url"
	"strings"
)

type GoalsResponse struct {
	Results []goal  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Goal stores Goal from Service
type goal struct {
	Id           string                     `json:"id"`
	Properties   json.RawMessage            `json:"properties"`
	CreatedAt    h_types.DateTimeString     `json:"createdAt"`
	UpdatedAt    h_types.DateTimeString     `json:"updatedAt"`
	Archived     bool                       `json:"archived"`
	Associations map[string]AssociationsSet `json:"associations"`
}
type Goal struct {
	Id         string
	Properties GoalProperties
	CreatedAt  h_types.DateTimeString
	UpdatedAt  h_types.DateTimeString
	Archived   bool
	//Associations     map[string]AssociationsSet
}

type GoalProperties struct {
	CreatedByUserId *go_types.Int64String `json:"hs_created_by_user_id,omitempty"`
	//CreatedDate      *h_types.DateTimeMSString `json:"hs_createdate,omitempty"`
	EndDateTime *h_types.DateTimeMSString `json:"hs_end_datetime,omitempty"`
	GoalName    *string                   `json:"hs_goal_name,omitempty"`
	//LastModifiedDate *h_types.DateTimeMSString `json:"hs_lastmodifieddate,omitempty"`
	StartDateTime *h_types.DateTimeMSString `json:"hs_start_datetime,omitempty"`
	ObjectId      *go_types.Int64String     `json:"hs_object_id,omitempty"`
	TargetAmount  *go_types.Float64String   `json:"hs_target_amount,omitempty"`
}

type GetGoalsConfig struct {
	Limit      *uint
	After      *string
	Properties *[]string
	//Associations *[]ObjectType
	Archived *bool
}

// GetGoals returns all goals
func (service *Service) GetGoals(config *GetGoalsConfig) (*[]Goal, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/goal_targets"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}

		var _properties []string
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				for _, p := range *config.Properties {
					_properties = append(_properties, string(p))
				}
			}
		}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				_properties = append(_properties, *config.Properties...)
			}
		}
		if len(_properties) > 0 {
			values.Set("properties", strings.Join(_properties, ","))
		}
		/*if config.Associations != nil {
			if len(*config.Associations) > 0 {
				_associations := []string{}
				for _, a := range *config.Associations {
					_associations = append(_associations, string(a))
				}
				values.Set("associations", strings.Join(_associations, ","))
			}
		}*/
		if config.Archived != nil {
			values.Set("archived", fmt.Sprintf("%v", *config.Archived))
		}
	}

	after := ""
	if config.After != nil {
		after = *config.After
	}

	var goals []Goal

	for {
		goalsResponse := GoalsResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &goalsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, d := range goalsResponse.Results {
			goal_, e := getGoal(&d, config.Properties)
			if e != nil {
				return nil, e
			}
			goals = append(goals, *goal_)
		}

		if config.After != nil { // explicit after parameter requested
			break
		}

		if goalsResponse.Paging == nil {
			break
		}

		if goalsResponse.Paging.Next.After == "" {
			break
		}

		after = goalsResponse.Paging.Next.After
	}

	return &goals, nil
}

func getGoal(goal *goal, properties *[]string) (*Goal, *errortools.Error) {
	goal_ := Goal{
		Id:         goal.Id,
		CreatedAt:  goal.CreatedAt,
		UpdatedAt:  goal.UpdatedAt,
		Archived:   goal.Archived,
		Properties: GoalProperties{},
	}

	if goal.Properties != nil {
		p := GoalProperties{}
		err := json.Unmarshal(goal.Properties, &p)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
		goal_.Properties = p
	}

	if properties != nil {
		p1 := make(map[string]string)
		err := json.Unmarshal(goal.Properties, &p1)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
	}

	return &goal_, nil
}
