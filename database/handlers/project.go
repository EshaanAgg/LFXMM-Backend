package handlers

import (
	"eshaanagg/lfx/database"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

func (client Client) CreateProject(proj database.Project) *database.Project {
	insertStmt := `
        INSERT INTO projects 
        (lfxProjectId, name, industry, description, skills, programYear, programTerm,  website, repository, amountRaised, organizationId) 
        VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
        RETURNING id;
    `

	err := client.QueryRow(insertStmt, proj.LFXProjectID, proj.Name, pq.Array(proj.Industry), proj.Description, pq.Array(proj.Skills), proj.ProgramYear, proj.ProgramTerm, proj.Website, proj.Repository, proj.AmountRaised, proj.OrganizationID).Scan(&proj.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add new project.")
		fmt.Println(err)
		return nil
	}

	return &proj
}

func (client Client) GetProjectsByParentOrgID(id string) []database.ProjectThumbail {
	queryStmt := `
        SELECT id, lfxProjectId, name, description, programYear, programTerm 
		FROM projects WHERE organizationId = $1
    `

	projects := make([]database.ProjectThumbail, 0)

	rowsRs, err := client.Query(queryStmt, id)
	if err != nil {
		fmt.Println("[ERROR] GetProjectsByParentOrgID query failed.")
		fmt.Println(err)
		return projects
	}

	for rowsRs.Next() {
		proj := database.ProjectThumbail{}
		lfxId := ""

		err := rowsRs.Scan(&proj.ID, &lfxId, &proj.Name, &proj.Description, &proj.ProgramYear, &proj.ProgramTerm)
		if err != nil {
			fmt.Println("[ERROR] Can't save to Project struct")
			return projects
		}
		proj.ProjectURL = "https://mentorship.lfx.linuxfoundation.org/project/" + lfxId
		projects = append(projects, proj)
	}

	return projects
}

func (client Client) GetProjectById(projectID string) ([]database.ProjectDetails, error) {
	queryStmt := `
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
		project := database.ProjectDetails{} // Create a single ProjectDetails struct to hold each row
		lfxId := ""
		var skillsStr string

		err := rows.Scan(
			&project.ProjectID,
			&lfxId,
			&project.Name,
			&project.Description,
			&project.Industry,
			&project.Website,
			&project.AmountRaised,
			&skillsStr,
			&project.OrganizationID,
			&project.Repository,
		)
		project.LFXProjectUrl = "https://mentorship.lfx.linuxfoundation.org/project/" + lfxId
		if err != nil {
			fmt.Println("[ERROR] GetProjectByProjectId scan failed.")
			fmt.Println(err)
			return nil, err
		}

		// Clean up and split the skills data
		cleanedSkills := strings.Split(strings.Trim(skillsStr, "{} "), ",")

		// Set the cleaned skills data in the project struct
		project.Skills = cleanedSkills

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("[ERROR] GetProjectByProjectId rows error.")
		fmt.Println(err)
		return nil, err
	}

	if len(projects) == 0 {
		return nil, fmt.Errorf("Project not found")
	}

	return projects, nil
}

func (client Client) GetProjectsByYear(id string, year int) []database.ProjectThumbail {

	queryStmt := `
		SELECT id, lfxProjectId, name, description, programYear, programTerm 
    	FROM projects 
    	WHERE organizationId = $1 AND programYear = $2
	`

	projects := make([]database.ProjectThumbail, 0)

	rowsRs, err := client.Query(queryStmt, id, year)
	if err != nil {
		fmt.Println("[ERROR] GetProjectsByYear query failed.")
		fmt.Println(err)
		return projects
	}

	for rowsRs.Next() {
		proj := database.ProjectThumbail{}
		lfxId := ""

		err := rowsRs.Scan(&proj.ID, &lfxId, &proj.Name, &proj.Description, &proj.ProgramYear, &proj.ProgramTerm)
		if err != nil {
			fmt.Println("[ERROR] Can't save to Project struct")
			return projects
		}
		proj.ProjectURL = "https://mentorship.lfx.linuxfoundation.org/project/" + lfxId
		projects = append(projects, proj)
	}
	return projects
}
