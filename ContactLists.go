package hubspot

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

type ContactListsResponse struct {
	Lists   []ContactList `json:"lists"`
	HasMore bool          `json:"has-more"`
	Offset  int64         `json:"offset"`
}

// ContactList stores ContactList from Service
type ContactList struct {
	MetaData        *ContactListMetaData  `json:"metaData,omitempty"`
	Name            string                `json:"name"`
	Filters         [][]ContactListFilter `json:"filters"`
	PortalId        int64                 `json:"portalId"`
	CreatedAt       *int64                `json:"createdAt,omitempty"`
	ListId          *int64                `json:"listId,omitempty"`
	AuthorId        *int64                `json:"authorId,omitempty"`
	UpdatedAt       *int64                `json:"updatedAt,omitempty"`
	ListType        *string               `json:"listType,omitempty"`
	InternalListId  *int64                `json:"internalListId,omitempty"`
	Deleteable      *bool                 `json:"deleteable,omitempty"`
	Archived        *bool                 `json:"archived,omitempty"`
	IlsFilterBranch *string               `json:"ilsFilterBranch,omitempty"`
	ReadOnly        *bool                 `json:"readOnly,omitempty"`
	Internal        *bool                 `json:"internal,omitempty"`
	LimitExempt     *bool                 `json:"limitExempt,omitempty"`
	Dynamic         bool                  `json:"dynamic"`
}

type ContactListMetaData struct {
	Processing                  string `json:"processing"`
	LastProcessingStateChangeAt int64  `json:"lastProcessingStateChangeAt"`
	Size                        int    `json:"size"`
	LastSizeChangeAt            int64  `json:"lastSizeChangeAt"`
	Error                       string `json:"error"`
	ListReferencesCount         *int64 `json:"listReferencesCount"`
	ParentFolderId              *int64 `json:"parentFolderId"`
}

type ContactListFilter struct {
	FilterFamily      *string `json:"filterFamily,omitempty"`
	WithinTimeMode    *string `json:"withinTimeMode,omitempty"`
	CheckPastVersions *bool   `json:"checkPastVersions,omitempty"`
	Type              string  `json:"type"`
	Property          string  `json:"property"`
	Value             string  `json:"value"`
	Operator          string  `json:"operator"`
}

type GetContactListsConfig struct {
	Offset *int64
	Count  *int64
}

// GetContactLists returns all contactLists
func (service *Service) GetContactLists(config *GetContactListsConfig) (*[]ContactList, *errortools.Error) {
	values := url.Values{}
	endpoint := "lists"

	var offset int64 = 0

	if config != nil {
		if config.Count != nil {
			values.Set("count", fmt.Sprintf("%v", *config.Count))
		}
		if config.Offset != nil {
			offset = *config.Offset
		}
	}

	contactLists := []ContactList{}

	for {
		contactListsResponse := ContactListsResponse{}

		if offset > 0 {
			values.Set("offset", fmt.Sprintf("%v", offset))
		}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlContacts(fmt.Sprintf("%s?%s", endpoint, values.Encode())),
			ResponseModel: &contactListsResponse,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		contactLists = append(contactLists, contactListsResponse.Lists...)

		if config != nil {
			if config.Offset != nil { // explicit after parameter requested
				break
			}
		}

		if !contactListsResponse.HasMore {
			break
		}

		offset = contactListsResponse.Offset
	}

	return &contactLists, nil
}

func (service *Service) CreateContactList(contactList *ContactList) (*ContactList, *errortools.Error) {
	endpoint := "lists"
	contactListNew := ContactList{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlContacts(endpoint),
		BodyModel:     contactList,
		ResponseModel: &contactListNew,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contactListNew, nil
}

func (service *Service) UpdateContactList(contactList *ContactList) (*ContactList, *errortools.Error) {
	if contactList.ListId == nil {
		return nil, errortools.ErrorMessage("ListId is required")
	}

	endpoint := "lists"
	contactListNew := ContactList{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlContacts(fmt.Sprintf("%s/%v", endpoint, *contactList.ListId)),
		BodyModel:     contactList,
		ResponseModel: &contactListNew,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &contactListNew, nil
}

func (service *Service) DeleteContactList(contactListId string) *errortools.Error {
	endpoint := "lists"

	requestConfig := go_http.RequestConfig{
		Method: http.MethodDelete,
		Url:    service.urlContacts(fmt.Sprintf("%s/%s", endpoint, contactListId)),
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return e
	}

	return nil
}

type AddContactsToContactListConfig struct {
	ListId int64    `json:"-"`
	Vids   []int    `json:"vids"`
	Emails []string `json:"emails"`
}

type AddContactsToContactListResponse struct {
	Updated       []int    `json:"updated"`
	Discarded     []int    `json:"discarded"`
	InvalidVids   []int    `json:"invalidVids"`
	InvalidEmails []string `json:"invalidEmails"`
}

func (service *Service) AddContactsToContactList(config *AddContactsToContactListConfig) (*AddContactsToContactListResponse, *errortools.Error) {
	endpoint := "lists"
	res := AddContactsToContactListResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.urlContacts(fmt.Sprintf("%s/%v/add", endpoint, config.ListId)),
		BodyModel:     config,
		ResponseModel: &res,
	}

	_, _, e := service.httpRequest(&requestConfig)
	if e != nil {
		return nil, e
	}

	return &res, nil
}

type GetContactsInContactListConfig struct {
	ListId int64 `json:"-"`
}

type GetContactsInContactListResponse struct {
	Contacts  []ContactInContactList `json:"contacts"`
	HasMore   bool                   `json:"has-more"`
	VidOffset int                    `json:"vid-offset"`
}

type ContactInContactList struct {
	AddedAt      int64 `json:"addedAt"`
	Vid          int   `json:"vid"`
	CanonicalVid int   `json:"canonical-vid"`
	MergedVids   []int `json:"merged-vids"`
	PortalId     int   `json:"portal-id"`
	IsContact    bool  `json:"is-contact"`
	Properties   struct {
		Firstname struct {
			Value string `json:"value"`
		} `json:"firstname,omitempty"`
		Lastmodifieddate struct {
			Value string `json:"value"`
		} `json:"lastmodifieddate"`
		Company struct {
			Value string `json:"value"`
		} `json:"company,omitempty"`
		Lastname struct {
			Value string `json:"value"`
		} `json:"lastname,omitempty"`
	} `json:"properties"`
	FormSubmissions  []interface{} `json:"form-submissions"`
	IdentityProfiles []struct {
		Vid                     int   `json:"vid"`
		SavedAtTimestamp        int64 `json:"saved-at-timestamp"`
		DeletedChangedTimestamp int   `json:"deleted-changed-timestamp"`
		Identities              []struct {
			Type      string `json:"type"`
			Value     string `json:"value"`
			Timestamp int64  `json:"timestamp"`
		} `json:"identities"`
	} `json:"identity-profiles"`
	MergeAudits []interface{} `json:"merge-audits"`
}

func (service *Service) GetContactsInContactList(config *GetContactsInContactListConfig) (*[]ContactInContactList, *errortools.Error) {
	values := url.Values{}
	endpoint := "lists"

	var contacts []ContactInContactList
	var offset int = 0

	for {
		if offset > 0 {
			values.Set("vidOffset", fmt.Sprintf("%v", offset))
		}

		res := GetContactsInContactListResponse{}

		requestConfig := go_http.RequestConfig{
			Method:        http.MethodGet,
			Url:           service.urlContacts(fmt.Sprintf("%s/%v/contacts/all?%s", endpoint, config.ListId, values.Encode())),
			ResponseModel: &res,
		}

		_, _, e := service.httpRequest(&requestConfig)
		if e != nil {
			return nil, e
		}

		contacts = append(contacts, res.Contacts...)

		if !res.HasMore {
			break
		}

		offset = res.VidOffset
	}
	return &contacts, nil
}
