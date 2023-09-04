package database

type ProjectThumbail struct {
	ID          string `json:"id"`
	ProjectURL  string `json:"projectUrl"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ProgramYear int    `json:"programYear"`
	ProgramTerm string `json:"programTerm"`
}
