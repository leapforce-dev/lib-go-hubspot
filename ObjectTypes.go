package hubspot

type ObjectType string

const (
	ObjectTypeCalls               ObjectType = "calls"
	ObjectTypeCompanies           ObjectType = "companies"
	ObjectTypeContacts            ObjectType = "contacts"
	ObjectTypeDeals               ObjectType = "deals"
	ObjectTypeEmails              ObjectType = "emails"
	ObjectTypeFeedbackSubmissions ObjectType = "feedback_submissions"
	ObjectTypeLineItems           ObjectType = "line_items"
	ObjectTypeMeetings            ObjectType = "meetings"
	ObjectTypeNotes               ObjectType = "notes"
	ObjectTypeProducts            ObjectType = "products"
	ObjectTypeQuotes              ObjectType = "quotes"
	ObjectTypeTickets             ObjectType = "tickets"
)

type CreateObjectConfig struct {
	ObjectType   string             `json:"-"`
	Properties   map[string]string  `json:"properties"`
	Associations *[]AssociationToV4 `json:"associations,omitempty"`
}

type BatchCreateObjectInput struct {
	Properties   map[string]string  `json:"properties"`
	Associations *[]AssociationToV4 `json:"associations,omitempty"`
}

type BatchCreateObjectsConfig struct {
	ObjectType string                   `json:"-"`
	Inputs     []BatchCreateObjectInput `json:"inputs"`
}

type BatchUpdateObjectInput struct {
	Id         string            `json:"id"`
	Properties map[string]string `json:"properties"`
}

type BatchUpdateObjectsConfig struct {
	ObjectType string                   `json:"-"`
	Inputs     []BatchUpdateObjectInput `json:"inputs"`
}

type UpdateObjectConfig struct {
	ObjectType string            `json:"-"`
	ObjectId   string            `json:"-"`
	Properties map[string]string `json:"properties"`
}
