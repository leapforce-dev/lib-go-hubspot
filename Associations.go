package hubspot

type AssociationsSet struct {
	Results []Association `json:"results"`
}

type Association struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
