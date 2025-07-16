package hubspot

type ObjectType string

const (
	ObjectTypeCalls               ObjectType = "calls"
	ObjectTypeCompanies           ObjectType = "companies"
	ObjectTypeContacts            ObjectType = "contacts"
	ObjectTypeCourses             ObjectType = "0-410"
	ObjectTypeDeals               ObjectType = "deals"
	ObjectTypeEmails              ObjectType = "emails"
	ObjectTypeFeedbackSubmissions ObjectType = "feedback_submissions"
	ObjectTypeLineItems           ObjectType = "line_items"
	ObjectTypeMeetings            ObjectType = "meetings"
	ObjectTypeNotes               ObjectType = "notes"
	ObjectTypeProducts            ObjectType = "products"
	ObjectTypeQuotes              ObjectType = "quotes"
	ObjectTypeTasks               ObjectType = "tasks"
	ObjectTypeTickets             ObjectType = "tickets"
)

type CreateObjectConfig struct {
	ObjectType   string             `json:"-"`
	Properties   map[string]string  `json:"properties"`
	Associations *[]AssociationToV4 `json:"associations,omitempty"`
}

type BatchObjectInput struct {
	Id           *string            `json:"id,omitempty"`
	Properties   map[string]string  `json:"properties"`
	Associations *[]AssociationToV4 `json:"associations,omitempty"`
}

type BatchObjectsConfig struct {
	ObjectType string             `json:"-"`
	Inputs     []BatchObjectInput `json:"inputs"`
}

type UpdateObjectConfig struct {
	ObjectType string            `json:"-"`
	ObjectId   string            `json:"-"`
	Properties map[string]string `json:"properties"`
}
