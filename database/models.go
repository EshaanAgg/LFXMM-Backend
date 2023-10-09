package database

type ParentOrg struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logoUrl"`
	Description    string  `json:"description"`
	Year    int  `json:"year"`
	Term    string  `json:"term"`
	Website        string  `json:"website"`

}
type Project struct {
	ID             string
	LFXProjectID   string
	Name           string
	Industry       []string
	Description    string
	Skills         []string
	Repository     string
	ProgramYear    int
	ProgramTerm    string
	AmountRaised   float64
	OrganizationID string
	Website        string
}
