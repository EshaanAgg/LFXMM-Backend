package database

// type ProjectThumbail struct {
// 	ID          string `json:"id"`
// 	ProjectURL  string `json:"projectUrl"`
// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	ProgramYear int    `json:"programYear"`
// 	ProgramTerm string `json:"programTerm"`
// }

type ProjectThumbail struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Industry    []string `json:"industry"`
	Description string   `json:"description"`
	Skills      []string `json:"skills"`
	Website     string   `json:"website"`
	Repository  string   `json:"repository"`
	ProgramYear int      `json:"programYear"`
	ProgramTerm string   `json:"programTerm"`
	ProjectURL  string   `json:"projectUrl"`
}
