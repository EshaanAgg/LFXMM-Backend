package database

type ParentOrg struct {
	ID   string
	Name string
	Logo string
}

type Organization struct {
	ID          string
	Name        string
	Logo        string
	ParentOrgID string
}

type Project struct {
	ID             string
	LFXProjectID   string
	Name           string
	Industry       []string
	Description    string
	Skills         []string
	ProgramYear    int
	ProgramTerm    string
	OrganizationID string
}
