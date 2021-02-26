package hubspot

type ObjectType string

const (
	ObjectTypeCompanies           ObjectType = "companies"
	ObjectTypeContacts            ObjectType = "contacts"
	ObjectTypeDeals               ObjectType = "deals"
	ObjectTypeFeedbackSubmissions ObjectType = "feedback_submissions"
	ObjectTypeLineItems           ObjectType = "line_items"
	ObjectTypeProducts            ObjectType = "products"
	ObjectTypeTickets             ObjectType = "tickets"
	ObjectTypeQuotes              ObjectType = "quotes"
)
