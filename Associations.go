package hubspot

type Associations struct {
	Companies struct {
		Results []*Association `json:"results"`
	} `json:"companies"`
	Contacts struct {
		Results []*Association `json:"results"`
	} `json:"contacts"`
}

type Association struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
