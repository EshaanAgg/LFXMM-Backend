package handlers

import (
	"database/sql"
	"eshaanagg/lfx/database"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

/*
 * Get all the project (thumbnails) whose `name` or `description` matches the provided `filterText`
 */
func (client Client) GetProjectsByFilter(filterText string) ([]database.ProjectThumbnail, error) {
	rowsRs, err := client.Query(`
		SELECT id, name, description, programYear, programTerm 
		FROM projects 
		WHERE name LIKE '%$1%' OR description LIKE '%$1$'
		ORDER BY name;
		`, filterText)

	if err != nil {
		fmt.Println("[ERROR] GetProjectsByFilter query failed")
		fmt.Println(err)
		return nil, err
	}
	defer rowsRs.Close()

	projects, err := parseAsProjectThumbnailSlice(rowsRs)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

/*
 * Used to save a project to the database
 */
func (client Client) CreateProject(proj database.Project) *database.Project {
	insertStmt := `
        INSERT INTO projects 
        (lfxProjectId, name, industry, description, skills, programYear, programTerm,  website, repository, amountRaised, organizationId) 
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
        RETURNING id;
    `

	err := client.QueryRow(
		insertStmt,
		proj.LFXProjectID,
		proj.Name,
		pq.Array(proj.Industry),
		proj.Description,
		pq.Array(proj.Skills),
		proj.ProgramYear,
		proj.ProgramTerm,
		proj.Website,
		proj.Repository,
		proj.AmountRaised,
		proj.OrganizationID,
	).Scan(&proj.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add new project.")
		fmt.Println(err)
		return nil
	}

	return &proj
}

/*
 * Get all the project (thumbnails) for a particular organization
 */
func (client Client) GetProjectsByParentOrgID(id string) []database.ProjectThumbnail {
	queryStmt := `
        SELECT id, lfxProjectId, name, description, programYear, programTerm, skills
		FROM projects WHERE organizationId = $1
    `

	rowsRs, err := client.Query(queryStmt, id)
	if err != nil {
		fmt.Println("[ERROR] GetProjectsByParentOrgID query failed.")
		fmt.Println(err)
		return nil
	}

	projects, err := parseAsProjectThumbnailSlice(rowsRs)
	if err != nil {
		return nil
	}

	return projects
}

/*
 * Returns the complete project object from the database
 */
func (client Client) GetProjectById(projectID string) (*database.ProjectDetails, error) {
	queryStmt :=
		`
    	SELECT id, lfxProjectId, name, description, industry, website, amountRaised, skills, organizationId, repository
    	FROM projects
    	WHERE id = $1
    `

	rows, err := client.Query(queryStmt, projectID)
	if err != nil {
		fmt.Println("[ERROR] GetProjectByProjectId query failed.")
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	projects := []database.ProjectDetails{}

	for rows.Next() {
		project := database.ProjectDetails{}
		lfxId := ""

		err := rows.Scan(
			&project.ProjectID,
			&lfxId,
			&project.Name,
			&project.Description,
			&project.Industry,
			&project.Website,
			&project.AmountRaised,
			pq.Array(&project.Skills),
			&project.OrganizationID,
			&project.Repository,
		)
		project.LFXProjectUrl = "https://mentorship.lfx.linuxfoundation.org/project/" + lfxId

		if err != nil {
			fmt.Println("[ERROR] GetProjectByProjectId scan failed.")
			fmt.Println(err)
			return nil, err
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("[ERROR] GetProjectByProjectId rows error.")
		fmt.Println(err)
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("project not found with the provided id")
	}

	return &projects[0], nil
}

/*
 * Returns project (thumbnails) filtered by the parent organization ID (`id`) and the year they came in (`year`)
 */
func (client Client) GetProjectsByYear(id string, year int) []database.ProjectThumbnail {

	queryStmt :=
		`
		SELECT id, lfxProjectId, name, description, programYear, programTerm, skills
    	FROM projects
    	WHERE organizationId = $1 AND programYear = $2
		ORDER BY name
		`

	rowsRs, err := client.Query(queryStmt, id, year)
	if err != nil {
		fmt.Println("[ERROR] GetProjectsByYear query failed.")
		fmt.Println(err)
		return nil
	}

	projects, err := parseAsProjectThumbnailSlice(rowsRs)
	if err != nil {
		return nil
	}

	return projects
}

/*
 * Returns project (thumbnails) filtered by the parent organization ID (`id`)
 */
func (client Client) GetProjectsByOrganization(id string) []database.ProjectThumbnail {

	queryStmt :=
		`
		SELECT id, lfxProjectId, name, description, programYear, programTerm, skills
    	FROM projects
    	WHERE organizationId = $1
		ORDER BY name
		`

	rowsRs, err := client.Query(queryStmt, id)
	if err != nil {
		fmt.Println("[ERROR] GetProjectsByOrganization query failed.")
		fmt.Println(err)
		return nil
	}

	projects, err := parseAsProjectThumbnailSlice(rowsRs)
	if err != nil {
		return nil
	}

	return projects
}

/*
 * Returns the count of projects per year for a given organization.
 */

func (client Client) GetCountOfProjectsByParentOrgID(id string) []database.ProjectCountByYear {
	queryStmt :=
		`
        SELECT programYear, COUNT(*) as count
		FROM projects WHERE organizationId = $1
		GROUP BY programYear
		ORDER BY programYear;
        `
	// Create a placeholder object, type-defined in responses.go
	counts := make([]database.ProjectCountByYear, 0)

	// Query the database
	rowsRs, err := client.Query(queryStmt, id)
	if err != nil {
		fmt.Println("[ERROR] GetCountOfProjectsByParentOrgID query failed.")
		fmt.Println(err)
		return counts
	}

	// Loop through rows of result appending the counts object
	for rowsRs.Next() {
		count := database.ProjectCountByYear{}

		err := rowsRs.Scan(&count.ProgramYear, &count.Count)
		if err != nil {
			fmt.Println("[ERROR] Can't save to Count struct")
			return counts
		}
		counts = append(counts, count)
	}

	return counts
}

/*
 * Used to update the skills list for a project
 */
func (client Client) SetSkillsForProject(id string, skillsArg []string) error {
	params := make([]string, 0, len(skillsArg))
	for i := range skillsArg {
		params = append(params, fmt.Sprintf("$%v", i+1))
	}

	updateStmt := fmt.Sprintf(`
	UPDATE projects
	SET skills = ARRAY[%s]
	WHERE id = $%v
	RETURNING id;
	`,
		strings.Join(params, ", "),
		len(params)+1,
	)

	// Convert the slice from []string to []any explicitly
	skills := make([]any, 0)
	for _, skill := range skillsArg {
		skills = append(skills, skill)
	}
	skills = append(skills, id)

	_, err := client.Exec(
		updateStmt,
		skills...,
	)

	if err != nil {
		fmt.Printf("[ERROR] Can't update the skills for project %v.\n", id)
		fmt.Println(err)
		return nil
	}

	return nil
}

// Helper function to convert the resultset of a SELECT * query to a slice of ProjectThumbail struct array.
func parseAsProjectThumbnailSlice(rowsRs *sql.Rows) ([]database.ProjectThumbnail, error) {
	projects := make([]database.ProjectThumbnail, 0)

	for rowsRs.Next() {
		proj := database.ProjectThumbnail{}
		lfxId := ""

		err := rowsRs.Scan(
			&proj.ID,
			&lfxId,
			&proj.Name,
			&proj.Description,
			&proj.ProgramYear,
			&proj.ProgramTerm,
			pq.Array(&proj.Skills),
		)

		if err != nil {
			fmt.Println("[ERROR] Can't save to ProjectThumbnail struct")
			return nil, err
		}

		proj.ProjectURL = "https://mentorship.lfx.linuxfoundation.org/project/" + lfxId
		projects = append(projects, proj)
	}

	return projects, nil
}
