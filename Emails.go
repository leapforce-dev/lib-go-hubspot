package hubspot

import (
	"encoding/json"
)

type EmailAssociations struct {
	Contacts AssociationsSet `json:"contacts"`
}

type EmailAssociation struct {
	To struct {
		Id string `json:"id"`
	} `json:"to"`
	Types []AssociationTypeV4 `json:"types"`
}

func NewEmailAssociation(toId string, category string, typeId int64) EmailAssociation {
	return EmailAssociation{
		To: struct {
			Id string `json:"id"`
		}{toId},
		Types: []AssociationTypeV4{
			{
				AssociationCategory: category,
				AssociationTypeId:   typeId,
			},
		},
	}
}

type EmailHeaderItem struct {
	Email     string  `json:"email"`
	FirstName *string `json:"firstName,omitempty"`
	LastName  *string `json:"lastName,omitempty"`
}

type EmailHeaders struct {
	From EmailHeaderItem    `json:"from"`
	To   *[]EmailHeaderItem `json:"to,omitempty"`
	Cc   *[]EmailHeaderItem `json:"cc,omitempty"`
	Bcc  *[]EmailHeaderItem `json:"bcc,omitempty"`
}

func SetEmailHeaders(properties map[string]string, headers *EmailHeaders) error {
	if headers == nil {
		properties["hs_email_headers"] = ""
		return nil
	}

	b, err := json.Marshal(headers)
	if err != nil {
		return err
	}

	properties["hs_email_headers"] = string(b)

	return nil
}
