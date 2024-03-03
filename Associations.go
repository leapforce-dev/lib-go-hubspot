package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"time"
)

type AssociationsSet struct {
	Results []Association `json:"results"`
}

type Association struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}

type AssociationsV4Set struct {
	Results []AssociationV4 `json:"results"`
}

type AssociationId struct {
	Id string `json:"id"`
}

type AssociationV4 struct {
	From AssociationId   `json:"from"`
	To   []AssociationTo `json:"to"`
}

type AssociationTo struct {
	ToObjectId       int64              `json:"toObjectId"`
	AssociationTypes []AssociationLabel `json:"associationTypes"`
}

func (a *AssociationTo) ToV4() *AssociationToV4 {
	if a == nil {
		return nil
	}

	a4 := AssociationToV4{
		To: AssociationId{
			Id: fmt.Sprintf("%v", a.ToObjectId),
		},
	}

	for _, t := range a.AssociationTypes {
		a4.Types = append(a4.Types, AssociationTypeV4{
			AssociationCategory: t.Category,
			AssociationTypeId:   t.TypeId,
		})
	}

	return &a4
}

type AssociationToV4 struct {
	To    AssociationId       `json:"to"`
	Types []AssociationTypeV4 `json:"types"`
}

type AssociationType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AssociationLabel struct {
	Category string  `json:"category"`
	TypeId   int64   `json:"typeId"`
	Label    *string `json:"label"`
}

type AssociationTypeV4 struct {
	AssociationCategory string `json:"associationCategory"`
	AssociationTypeId   int64  `json:"associationTypeId"`
}

type BatchGetAssociationsInput struct {
	Id string `json:"id"`
}

type BatchGetAssociationsConfig struct {
	FromObjectType string                      `json:"-"`
	ToObjectType   string                      `json:"-"`
	Inputs         []BatchGetAssociationsInput `json:"inputs"`
}

func (service *Service) BatchGetAssociations(config *BatchGetAssociationsConfig) (*AssociationsV4Set, *errortools.Error) {
	if config == nil {
		return nil, nil
	}
	if len(config.Inputs) == 0 {
		return nil, nil
	}

	endpoint := fmt.Sprintf("associations/%v/%v/batch/read", config.FromObjectType, config.ToObjectType)

	var associationsV4Set AssociationsV4Set

	for _, batch := range service.batches(len(config.Inputs)) {
		var associationsV4Set_ AssociationsV4Set

		requestConfig := go_http.RequestConfig{
			Method: http.MethodPost,
			Url:    service.urlV4(endpoint),
			BodyModel: BatchGetAssociationsConfig{
				Inputs: config.Inputs[batch.startIndex:batch.endIndex],
			},
			ResponseModel: &associationsV4Set_,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		associationsV4Set.Results = append(associationsV4Set.Results, associationsV4Set_.Results...)
	}

	return &associationsV4Set, nil
}

type CreateAssociationConfig struct {
	FromObjectType   string
	FromObjectId     string
	ToObjectType     string
	ToObjectId       string
	AssociationTypes []AssociationTypeV4
}

type CreateAssociationResponse struct {
	FromObjectTypeId string   `json:"fromObjectTypeId"`
	FromObjectId     int64    `json:"fromObjectId"`
	ToObjectTypeId   string   `json:"toObjectTypeId"`
	ToObjectId       int64    `json:"toObjectId"`
	Labels           []string `json:"labels"`
}

func (service *Service) CreateAssociation(config *CreateAssociationConfig) (*CreateAssociationResponse, *errortools.Error) {
	if config == nil {
		return nil, nil
	}

	endpoint := fmt.Sprintf("objects/%s/%s/associations/%s/%s", config.FromObjectType, config.FromObjectId, config.ToObjectType, config.ToObjectId)

	var createAssociationResponse CreateAssociationResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPut,
		Url:           service.urlV4(endpoint),
		BodyModel:     config.AssociationTypes,
		ResponseModel: &createAssociationResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &createAssociationResponse, nil
}

type BatchCreateAssociationsInput struct {
	Types []AssociationTypeV4 `json:"types"`
	From  AssociationId       `json:"from"`
	To    AssociationId       `json:"to"`
}

type BatchCreateAssociationsConfig struct {
	FromObjectType string                         `json:"-"`
	ToObjectType   string                         `json:"-"`
	Inputs         []BatchCreateAssociationsInput `json:"inputs"`
}

type BatchCreateAssociationsResponse struct {
	CompletedAt *time.Time `json:"completedAt"`
	RequestedAt *time.Time `json:"requestedAt"`
	StartedAt   *time.Time `json:"startedAt"`
	Links       struct {
		AdditionalProp1 string `json:"additionalProp1"`
		AdditionalProp2 string `json:"additionalProp2"`
		AdditionalProp3 string `json:"additionalProp3"`
	} `json:"links"`
	Results []CreateAssociationResponse `json:"results"`
	Status  string                      `json:"status"`
}

func (service *Service) BatchCreateAssociations(config *BatchCreateAssociationsConfig) (*[]CreateAssociationResponse, *errortools.Error) {
	if config == nil {
		return nil, nil
	}
	if len(config.Inputs) == 0 {
		return nil, nil
	}

	endpoint := fmt.Sprintf("associations/%s/%s/batch/create", config.FromObjectType, config.ToObjectType)

	var r []CreateAssociationResponse

	for _, batch := range service.batches(len(config.Inputs)) {
		var batchCreateAssociationsResponse BatchCreateAssociationsResponse

		requestConfig := go_http.RequestConfig{
			Method: http.MethodPost,
			Url:    service.urlV4(endpoint),
			BodyModel: BatchCreateAssociationsConfig{
				Inputs: config.Inputs[batch.startIndex:batch.endIndex],
			},
			ResponseModel: &batchCreateAssociationsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		r = append(r, batchCreateAssociationsResponse.Results...)
	}

	return &r, nil
}

type BatchArchiveAssociationsInput struct {
	From AssociationId   `json:"from"`
	To   []AssociationId `json:"to"`
}

type BatchArchiveAssociationsConfig struct {
	FromObjectType string                          `json:"-"`
	ToObjectType   string                          `json:"-"`
	Inputs         []BatchArchiveAssociationsInput `json:"inputs"`
}

func (service *Service) BatchArchiveAssociations(config *BatchArchiveAssociationsConfig) *errortools.Error {
	if config == nil {
		return nil
	}
	if len(config.Inputs) == 0 {
		return nil
	}

	endpoint := fmt.Sprintf("associations/%s/%s/batch/archive", config.FromObjectType, config.ToObjectType)

	for _, batch := range service.batches(len(config.Inputs)) {
		requestConfig := go_http.RequestConfig{
			Method: http.MethodPost,
			Url:    service.urlV4(endpoint),
			BodyModel: BatchArchiveAssociationsConfig{
				Inputs: config.Inputs[batch.startIndex:batch.endIndex],
			},
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return e
		}
	}

	return nil
}

type GetAssociationsConfig struct {
	FromObjectType string
	FromObjectId   string
	ToObjectType   string
}

type GetAssociationsResponse struct {
	Results []AssociationTo `json:"results"`
}

func (service *Service) GetAssociations(config *GetAssociationsConfig) (*GetAssociationsResponse, *errortools.Error) {
	if config == nil {
		return nil, nil
	}

	endpoint := fmt.Sprintf("objects/%s/%s/associations/%s", config.FromObjectType, config.FromObjectId, config.ToObjectType)

	var getAssociationsResponse GetAssociationsResponse

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlV4(endpoint),
		ResponseModel: &getAssociationsResponse,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &getAssociationsResponse, nil
}

type DeleteAssociationConfig struct {
	FromObjectType string
	FromObjectId   string
	ToObjectType   string
	ToObjectId     string
}

func (service *Service) DeleteAssociation(config *DeleteAssociationConfig) *errortools.Error {
	if config == nil {
		return nil
	}

	endpoint := fmt.Sprintf("objects/%s/%s/associations/%s/%s", config.FromObjectType, config.FromObjectId, config.ToObjectType, config.ToObjectId)

	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.urlV4(endpoint),
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

type GetAssociationTypesConfig struct {
	FromObjectType string
	ToObjectType   string
}

func (service *Service) GetAssociationTypes(config *GetAssociationTypesConfig) (*[]AssociationType, *errortools.Error) {
	var response struct {
		Results []AssociationType `json:"results"`
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlCrm(fmt.Sprintf("associations/%s/%s/types", config.FromObjectType, config.ToObjectType)),
		ResponseModel: &response,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Results, nil
}

type GetAssociationLabelsConfig struct {
	FromObjectType string
	ToObjectType   string
}

func (service *Service) GetAssociationLabels(config *GetAssociationLabelsConfig) (*[]AssociationLabel, *errortools.Error) {
	var response struct {
		Results []AssociationLabel `json:"results"`
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlV4(fmt.Sprintf("associations/%s/%s/labels", config.FromObjectType, config.ToObjectType)),
		ResponseModel: &response,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &response.Results, nil
}
