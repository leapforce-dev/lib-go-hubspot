package hubspot

type AssociationsSet struct {
	Results []Association `json:"results"`
}

type Association struct {
	Id   string `json:"id"`
	Type string `json:"type"`
}
