package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

const maxBatchSize int = 10000

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

type AssociationV4 struct {
	From struct {
		Id string `json:"id"`
	} `json:"from"`
	To []struct {
		ToObjectId       int64             `json:"toObjectId"`
		AssociationTypes []AssociationType `json:"associationTypes"`
	} `json:"to"`
}

type AssociationType struct {
	Category string  `json:"category"`
	TypeId   int64   `json:"typeId"`
	Label    *string `json:"label"`
}

type AssociationTypeV4 struct {
	AssociationCategory string `json:"associationCategory"`
	AssociationTypeId   int64  `json:"associationTypeId"`
}

type BatchGetAssociationsConfig struct {
	FromObjectType string
	ToObjectType   string
	Ids            []string
}

func (service *Service) BatchGetAssociations(config *BatchGetAssociationsConfig) (*AssociationsV4Set, *errortools.Error) {
	if len(config.Ids) == 0 {
		return nil, nil
	}

	endpoint := fmt.Sprintf("associations/%v/%v/batch/read", config.FromObjectType, config.ToObjectType)

	ids := config.Ids
	var associationsV4Set AssociationsV4Set

	for len(ids) > 0 {

		var body struct {
			Inputs []struct {
				Id string `json:"id"`
			} `json:"inputs"`
		}

		for i, id := range ids {
			if i == maxBatchSize {
				break
			}
			idStruct := struct {
				Id string `json:"id"`
			}{id}
			body.Inputs = append(body.Inputs, idStruct)
		}

		var associationsV4Set_ AssociationsV4Set

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodPost,
			Url:           service.urlV4(endpoint),
			BodyModel:     body,
			ResponseModel: &associationsV4Set_,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		associationsV4Set.Results = append(associationsV4Set.Results, associationsV4Set_.Results...)

		if len(ids) > maxBatchSize {
			ids = ids[maxBatchSize:]
		} else {
			break
		}
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

type GetAssociationsConfig struct {
	FromObjectType string
	FromObjectId   string
	ToObjectType   string
}

type GetAssociationsResponse struct {
	Results []struct {
		ToObjectId       int64             `json:"toObjectId"`
		AssociationTypes []AssociationType `json:"associationTypes"`
	} `json:"results"`
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
	Ids            []string
}

func (service *Service) GetAssociationTypes(config *GetAssociationTypesConfig) (*[]AssociationType, *errortools.Error) {
	var response struct {
		Results []AssociationType `json:"results"`
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
