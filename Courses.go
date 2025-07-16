package hubspot

import (
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CoursesResponse struct {
	Results []Course `json:"results"`
	Paging  *Paging  `json:"paging"`
}

// Course stores Course from Service
type Course struct {
	Id                    string                       `json:"id"`
	Properties            map[string]string            `json:"properties"`
	CreatedAt             h_types.DateTimeString       `json:"createdAt"`
	UpdatedAt             h_types.DateTimeString       `json:"updatedAt"`
	Archived              bool                         `json:"archived"`
	Associations          map[string]AssociationsSet   `json:"associations"`
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory"`
}

type GetCoursesConfig struct {
	Limit                 *uint
	After                 *string
	Properties            *[]string
	PropertiesWithHistory *[]string
	Associations          *[]string
	Archived              *bool
}

// GetCourses returns all courses
func (service *Service) GetCourses(config *GetCoursesConfig) (*[]Course, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/0-410"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		_properties := []string{}
		if config.Properties != nil {
			if len(*config.Properties) > 0 {
				_properties = append(_properties, *config.Properties...)
			}
		}
		if config.PropertiesWithHistory != nil {
			if len(*config.PropertiesWithHistory) > 0 {
				values.Set("propertiesWithHistory", strings.Join(*config.PropertiesWithHistory, ","))
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

	courses := []Course{}

	for {
		coursesResponse := CoursesResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &coursesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		courses = append(courses, coursesResponse.Results...)

		if config.After != nil { // explicit after parameter requested
			break
		}

		if coursesResponse.Paging == nil {
			break
		}

		if coursesResponse.Paging.Next.After == "" {
			break
		}

		after = coursesResponse.Paging.Next.After
	}

	return &courses, nil
}

func (service *Service) CreateCourse(config *CreateObjectConfig) (*Course, *errortools.Error) {
	endpoint := "objects/0-410"
	course := Course{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &course,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &course, nil
}

type BatchCoursesResponse struct {
	CompletedAt *time.Time        `json:"completedAt"`
	NumErrors   int               `json:"numErrors"`
	RequestedAt *time.Time        `json:"requestedAt"`
	StartedAt   *time.Time        `json:"startedAt"`
	Links       map[string]string `json:"links"`
	Results     []Course          `json:"results"`
	Errors      []struct {
		SubCategory json.RawMessage   `json:"subCategory"`
		Context     map[string]string `json:"context"`
		Links       map[string]string `json:"links"`
		Id          string            `json:"id"`
		Category    string            `json:"category"`
		Message     string            `json:"message"`
		Errors      []struct {
			SubCategory string `json:"subCategory"`
			Code        string `json:"code"`
			In          string `json:"in"`
			Context     struct {
				MissingScopes []string `json:"missingScopes"`
			} `json:"context"`
			Message string `json:"message"`
		} `json:"errors"`
		Status string `json:"status"`
	} `json:"errors"`
	Status string `json:"status"`
}

func (service *Service) BatchCreateCourses(config *BatchObjectsConfig) (*[]Course, *errortools.Error) {
	var courses []Course

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchCoursesResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(fmt.Sprintf("objects/%s/batch/create", config.ObjectType)),
			BodyModel:     BatchObjectsConfig{Inputs: config.Inputs[batch.startIndex:batch.endIndex]},
			ResponseModel: &r,
		}

		_, response, e := service.httpRequest(&requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusMultiStatus {
				fmt.Println(r.Errors)
				goto ok
			}
		}
		if e != nil {
			return nil, e
		}
	ok:
		courses = append(courses, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &courses, nil
}

func (service *Service) BatchUpdateCourses(config *BatchObjectsConfig) (*[]Course, *errortools.Error) {
	var courses []Course

	for _, batch := range service.batches(len(config.Inputs)) {
		var r BatchCoursesResponse

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlCrm(fmt.Sprintf("objects/%s/batch/update", config.ObjectType)),
			BodyModel:     BatchObjectsConfig{Inputs: config.Inputs[batch.startIndex:batch.endIndex]},
			ResponseModel: &r,
		}

		_, response, e := service.httpRequest(&requestConfig)
		if response != nil {
			if response.StatusCode == http.StatusMultiStatus {
				fmt.Println(r.Errors)
				goto ok
			}
		}
		if e != nil {
			return nil, e
		}
	ok:
		courses = append(courses, r.Results...)

		fmt.Println("batch", batch.startIndex)
	}

	return &courses, nil
}

func (service *Service) UpdateCourse(config *UpdateObjectConfig) (*Course, *errortools.Error) {
	endpoint := "objects/0-410"
	course := Course{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.ObjectId)),
		BodyModel:     config,
		ResponseModel: &course,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &course, nil
}

func (service *Service) BatchArchiveCourses(courseIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(courseIds) > index {
		if len(courseIds) > index+maxItemsPerBatch {
			e := service.batchArchiveCourses(courseIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchArchiveCourses(courseIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchArchiveCourses(courseIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, courseId := range courseIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{courseId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/0-410/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
