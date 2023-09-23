package hubspot

import (
	"bytes"
	"encoding/json"
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	h_types "github.com/leapforce-libraries/go_hubspot/types"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
)

type FilesResponse struct {
	Results []File  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// File stores File from Service
type File struct {
	Id                string                 `json:"id"`
	CreatedAt         h_types.DateTimeString `json:"createdAt"`
	ArchivedAt        h_types.DateTimeString `json:"archivedAt"`
	UpdatedAt         h_types.DateTimeString `json:"updatedAt"`
	Archived          bool                   `json:"archived"`
	ParentFolderId    string                 `json:"parentFolderId"`
	Name              string                 `json:"name"`
	Path              string                 `json:"path"`
	Size              int                    `json:"size"`
	Height            int                    `json:"height"`
	Width             int                    `json:"width"`
	Encoding          string                 `json:"encoding"`
	Type              string                 `json:"type"`
	Extension         string                 `json:"extension"`
	DefaultHostingUrl string                 `json:"defaultHostingUrl"`
	Url               string                 `json:"url"`
	IsUsableInContent bool                   `json:"isUsableInContent"`
	Access            string                 `json:"access"`
}

type UploadFileConfig struct {
	File         []byte            `json:"file"`
	FolderId     *string           `json:"folderId,omitempty"`
	FolderPath   *string           `json:"folderPath,omitempty"`
	FileName     string            `json:"fileName"`
	CharsetHunch *string           `json:"charsetHunch,omitempty"`
	Options      UploadFileOptions `json:"options,omitempty"`
}

type UploadFileOptions struct {
	Access                      string  `json:"access"`
	TTL                         *string `json:"ttl,omitempty"`
	Overwrite                   *bool   `json:"overwrite,omitempty"`
	DuplicateValidationStrategy *string `json:"duplicateValidationStrategy,omitempty"`
	DuplicateValidationScope    *string `json:"duplicateValidationScope,omitempty"`
}

// UploadFile uploads a file to Hubspot
func (service *Service) UploadFile(config *UploadFileConfig) (*File, *errortools.Error) {
	endpoint := "files"

	var file File

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if config.FolderPath != nil {
		err := writer.WriteField("folderPath", *config.FolderPath)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
	}

	if config.FolderId != nil {
		err := writer.WriteField("folderId", *config.FolderId)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}
	}

	err := writer.WriteField("fileName", config.FileName)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	b, err := json.Marshal(config.Options)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	err = writer.WriteField("options", string(b))
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	part, err := writer.CreateFormFile("file", config.FileName)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	_, err = io.Copy(part, bytes.NewBuffer(config.File))
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	err = writer.Close()
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	re, err := http.NewRequest(http.MethodPost, service.urlFiles(endpoint), body)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}
	re.Header.Add("Content-Type", writer.FormDataContentType())
	if service.authorizationMode == authorizationModeAccessToken {
		re.Header.Add("Authorization", fmt.Sprintf("Bearer %s", service.accessToken))
	} else {
		return nil, errortools.ErrorMessage("Method only allowed with access token")
	}

	client := &http.Client{}
	res, err := client.Do(re)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	br, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	err = json.Unmarshal(br, &file)
	if err != nil {
		return nil, errortools.ErrorMessage(err)
	}

	if file.Id == "" {
		err = json.Unmarshal(br, &service.errorResponse)
		if err != nil {
			return nil, errortools.ErrorMessage(err)
		}

		if service.errorResponse.Status == "error" {
			return nil, errortools.ErrorMessagef("error: %s", service.errorResponse.Message)
		}
	}

	return &file, nil
}

type GetFileConfig struct {
	FileId     string
	Properties *[]string
}

// GetFile retrieves a file from Hubspot
func (service *Service) GetFile(config *GetFileConfig) (*File, *errortools.Error) {
	if config == nil {
		return nil, errortools.ErrorMessage("GetFileConfig must not be a nil pointer")
	}

	values := url.Values{}
	endpoint := "files"

	if config.Properties != nil {
		if len(*config.Properties) > 0 {
			values.Set("properties", strings.Join(*config.Properties, ","))
		}
	}

	var file File

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlFiles(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
		ResponseModel: &file,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &file, nil
}
