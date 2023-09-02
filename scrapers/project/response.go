package project

type projectResponse struct {
	ProjectID       string `json:"projectId"`
	Name            string `json:"name"`
	Industry        string `json:"industry"`
	Description     string `json:"description"`
	ApprenticeNeeds struct {
		Skills []string `json:"skills"`
	} `json:"apprenticeNeeds"`
	Repository   string  `json:"repoLink"`
	Website      string  `json:"websiteUrl"`
	LogoURL      string  `json:"logoUrl"`
	CreatedOn    string  `json:"createdOn"`
	AmountRaised float64 `json:"amountRaised"`
}
