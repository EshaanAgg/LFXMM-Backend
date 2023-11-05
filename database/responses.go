package database

type ProjectThumbnail struct {
	ID          string   `json:"id"`
	ProjectURL  string   `json:"projectUrl"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ProgramYear int      `json:"programYear"`
	ProgramTerm string   `json:"programTerm"`
	Skills      []string `json:"skills"`
}

type ProjectDetails struct {
	ProjectID      string   `json:"projectId"`
	OrganizationID string   `json:"orgId"`
	LFXProjectUrl  string   `json:"lfxProjectUrl"`
	Name           string   `json:"name"`
	Industry       string   `json:"industry"`
	Description    string   `json:"description"`
	Skills         []string `json:"skills"`
	Repository     string   `json:"repoLink"`
	Website        string   `json:"websiteUrl"`
	LogoURL        string   `json:"logoUrl"`
	CreatedOn      string   `json:"createdOn"`
	AmountRaised   float64  `json:"amountRaised"`
}

type ProjectCountByYear struct {
	ProgramYear int `json:"programYear"`
	Count       int `json:"count"`
}

type ProjectDescription struct {
	ProjectID    string   `json:"projectId"`
	LFXId        string   `json:"lfid"`
	Status       string   `json:"status"`
	Industry     string   `json:"industry"`
	Description  string   `json:"description"`
	Repository   string   `json:"repolink"`
	Skills       []string `json:"skills"`
	AmountRaised float64  `json:"amountRaised"`
}
