package hubspot

type Paging struct {
	Next struct {
		After string `json:"after"`
		Link  string `json:"link"`
	} `json:"next"`
}
