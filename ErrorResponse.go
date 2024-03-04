package hubspot

// ErrorResponse stores general API error response
type ErrorResponse struct {
	Status        string            `json:"status"`
	Message       string            `json:"message"`
	CorrelationId string            `json:"correlationId"`
	Category      string            `json:"category"`
	Links         map[string]string `json:"links"`
}

type PropertyError struct {
	IsValid bool   `json:"isValid"`
	Message string `json:"message"`
	Error   string `json:"error"`
	Name    string `json:"name"`
}
