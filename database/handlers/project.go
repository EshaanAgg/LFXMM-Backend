package handlers

import (
	"eshaanagg/lfx/database"
	"fmt"

	"github.com/lib/pq"
)

func (client Client) CreateProject(proj database.Project) *database.Project {
	insertStmt :=
		`
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
	queryStmt :=
		`
        SELECT id, lfxProjectId, name, description, programYear, programYear 
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

		err := rowsRs.Scan(&proj.ID, &lfxId, &proj.Name, &proj.Description, &proj.ProgramYear, &proj.ProgramYear)
		if err != nil {
			fmt.Println("[ERROR] Can't save to Project struct")
			return projects
		}
		proj.ProjectURL = "https://mentorship.lfx.linuxfoundation.org/project/" + lfxId
		projects = append(projects, proj)
	}

	return projects
}
