package hubspot

type WebhookPayload struct {
	AppId            int64  `json:"appId"`
	EventId          int64  `json:"eventId"`
	SubscriptionId   int64  `json:"subscriptionId"`
	PortalId         int64  `json:"portalId"`
	OccurredAt       int64  `json:"occurredAt"`
	SubscriptionType string `json:"subscriptionType"`
	AttemptNumber    int64  `json:"attemptNumber"`
	ObjectId         int64  `json:"objectId"`
	ChangeSource     string `json:"changeSource"`
	PropertyName     string `json:"propertyName"`
	PropertyValue    string `json:"propertyValue"`
	IsSensitive      bool   `json:"isSensitive"`
}
