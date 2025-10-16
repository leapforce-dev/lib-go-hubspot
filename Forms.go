package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
)

type FormsResponse struct {
	Results []Form  `json:"results"`
	Paging  *Paging `json:"paging"`
}

// Form stores Form from Service
type Form struct {
	PortalId                int           `json:"portalId"`
	Guid                    string        `json:"guid"`
	Name                    string        `json:"name"`
	Action                  string        `json:"action"`
	Method                  string        `json:"method"`
	CssClass                string        `json:"cssClass"`
	Redirect                string        `json:"redirect"`
	SubmitText              string        `json:"submitText"`
	FollowUpId              string        `json:"followUpId"`
	NotifyRecipients        string        `json:"notifyRecipients"`
	LeadNurturingCampaignId string        `json:"leadNurturingCampaignId"`
	FormFieldGroups         interface{}   `json:"formFieldGroups"`
	CreatedAt               int64         `json:"createdAt"`
	UpdatedAt               int64         `json:"updatedAt"`
	PerformableHtml         string        `json:"performableHtml"`
	MigratedFrom            string        `json:"migratedFrom"`
	IgnoreCurrentValues     bool          `json:"ignoreCurrentValues"`
	MetaData                interface{}   `json:"metaData"`
	Deletable               bool          `json:"deletable"`
	InlineMessage           string        `json:"inlineMessage"`
	TmsId                   string        `json:"tmsId"`
	CaptchaEnabled          bool          `json:"captchaEnabled"`
	CampaignGuid            string        `json:"campaignGuid"`
	Cloneable               bool          `json:"cloneable"`
	Editable                bool          `json:"editable"`
	FormType                string        `json:"formType"`
	DeletedAt               int           `json:"deletedAt"`
	ThemeName               string        `json:"themeName"`
	ParentId                int           `json:"parentId"`
	Style                   string        `json:"style"`
	IsPublished             bool          `json:"isPublished"`
	PublishAt               int           `json:"publishAt"`
	UnpublishAt             int           `json:"unpublishAt"`
	PublishedAt             int           `json:"publishedAt"`
	CustomUid               string        `json:"customUid"`
	CreateMarketableForm    bool          `json:"createMarketableForm"`
	EditVersion             int           `json:"editVersion"`
	ThankYouMessageJson     string        `json:"thankYouMessageJson"`
	ThemeColor              string        `json:"themeColor"`
	AlwaysCreateNewCompany  bool          `json:"alwaysCreateNewCompany"`
	InternalUpdatedAt       int64         `json:"internalUpdatedAt"`
	BusinessUnitId          int           `json:"businessUnitId"`
	PortableKey             string        `json:"portableKey"`
	SelectedExternalOptions []interface{} `json:"selectedExternalOptions"`
	EmbedVersion            string        `json:"embedVersion"`
	Enrichable              bool          `json:"enrichable"`
}

// GetForm returns a specific form
func (service *Service) GetForm(formId string) (*Form, *errortools.Error) {
	endpoint := "forms"

	form := Form{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.urlForms(fmt.Sprintf("%s/%s", endpoint, formId)),
		ResponseModel: &form,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &form, nil
}
