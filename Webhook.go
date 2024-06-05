package hubspot

type WebhookPayload struct {
	AppId            int    `json:"appId"`
	EventId          int    `json:"eventId"`
	SubscriptionId   int    `json:"subscriptionId"`
	PortalId         int    `json:"portalId"`
	OccurredAt       int64  `json:"occurredAt"`
	SubscriptionType string `json:"subscriptionType"`
	AttemptNumber    int    `json:"attemptNumber"`
	ObjectId         int    `json:"objectId"`
	ChangeSource     string `json:"changeSource"`
	PropertyName     string `json:"propertyName"`
	PropertyValue    string `json:"propertyValue"`
	IsSensitive      bool   `json:"isSensitive"`
}
