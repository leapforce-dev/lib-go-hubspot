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
)

type NotesResponse struct {
	Results []note  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Note stores Note from Service
type note struct {
	Id                    string                       `json:"id"`
	Properties            json.RawMessage              `json:"properties"`
	CreatedAt             h_types.DateTimeString       `json:"createdAt"`
	UpdatedAt             h_types.DateTimeString       `json:"updatedAt"`
	Archived              bool                         `json:"archived"`
	Associations          map[string]AssociationsSet   `json:"associations"`
	PropertiesWithHistory map[string][]PropertyHistory `json:"propertiesWithHistory"`
}

type Note struct {
	Id               string
	CustomProperties map[string]string
	CreatedAt        h_types.DateTimeString
	UpdatedAt        h_types.DateTimeString
	Archived         bool
	Associations     map[string]AssociationsSet
}

type GetNotesConfig struct {
	Limit            *uint
	After            *string
	CustomProperties *[]string
	Associations     *[]string
	Archived         *bool
}

// GetNotes returns all notes
func (service *Service) GetNotes(config *GetNotesConfig) (*[]Note, *errortools.Error) {
	values := url.Values{}
	endpoint := "objects/notes"

	if config != nil {
		if config.Limit != nil {
			values.Set("limit", fmt.Sprintf("%v", *config.Limit))
		}
		var _properties []string
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
				var _associations []string
				for _, a := range *config.Associations {
					_associations = append(_associations, a)
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

	var notes []Note

	for {
		notesResponse := NotesResponse{}

		if after != "" {
			values.Set("after", after)
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlCrm(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &notesResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		for _, d := range notesResponse.Results {
			note_, e := getNote(&d, config.CustomProperties)
			if e != nil {
				return nil, e
			}
			notes = append(notes, *note_)
		}

		if config.After != nil { // explicit after parameter requested
			break
		}

		if notesResponse.Paging == nil {
			break
		}

		if notesResponse.Paging.Next.After == "" {
			break
		}

		after = notesResponse.Paging.Next.After
	}

	return &notes, nil
}

func getNote(note *note, customProperties *[]string) (*Note, *errortools.Error) {
	note_ := Note{
		Id:               note.Id,
		CreatedAt:        note.CreatedAt,
		UpdatedAt:        note.UpdatedAt,
		Archived:         note.Archived,
		Associations:     note.Associations,
		CustomProperties: make(map[string]string),
	}

	if customProperties != nil {
		p1 := make(map[string]string)
		err := json.Unmarshal(note.Properties, &p1)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}

		for _, cp := range *customProperties {
			value, ok := p1[cp]
			if ok {
				note_.CustomProperties[cp] = value
			}
		}
	}

	return &note_, nil
}

type CreateNoteConfig struct {
	Properties   map[string]string `json:"properties"`
	Associations *[]AssociationTo  `json:"associations,omitempty"`
}

func (service *Service) CreateNote(config *CreateNoteConfig) (*Note, *errortools.Error) {
	endpoint := "objects/notes"
	note := Note{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlCrm(endpoint),
		BodyModel:     config,
		ResponseModel: &note,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &note, nil
}

type UpdateNoteConfig struct {
	NoteId           string
	CustomProperties map[string]string
}

func (service *Service) UpdateNote(config *UpdateNoteConfig) (*Note, *errortools.Error) {
	endpoint := "objects/notes"
	note := Note{}

	var properties = make(map[string]string)

	if config.CustomProperties != nil {
		for key, value := range config.CustomProperties {
			properties[key] = value
		}
	}

	var properties_ = struct {
		Properties map[string]string `json:"properties"`
	}{
		properties,
	}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPatch,
		Url:           service.urlCrm(fmt.Sprintf("%s/%s", endpoint, config.NoteId)),
		BodyModel:     properties_,
		ResponseModel: &note,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &note, nil
}

func (service *Service) BatchDeleteNotes(noteIds []string) *errortools.Error {
	var maxItemsPerBatch = 100
	var index = 0
	for len(noteIds) > index {
		if len(noteIds) > index+maxItemsPerBatch {
			e := service.batchDeleteNotes(noteIds[index : index+maxItemsPerBatch])
			if e != nil {
				return e
			}
		} else {
			e := service.batchDeleteNotes(noteIds[index:])
			if e != nil {
				return e
			}
		}

		index += maxItemsPerBatch
	}

	return nil
}

func (service *Service) batchDeleteNotes(noteIds []string) *errortools.Error {
	var body struct {
		Inputs []struct {
			Id string `json:"id"`
		} `json:"inputs"`
	}

	for _, noteId := range noteIds {
		body.Inputs = append(body.Inputs, struct {
			Id string `json:"id"`
		}{noteId})
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.urlCrm("objects/notes/batch/archive"),
		BodyModel: body,
	}

	_, _, e := service.httpRequest(&requestConfig)
	return e
}
