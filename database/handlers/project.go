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

	err := client.QueryRow(insertStmt, proj.LFXProjectID, proj.Name, pq.Array(proj.Industry), proj.Description, pq.Array(proj.Skills), proj.ProgramYear, proj.ProgramYear, proj.Website, proj.Repository, proj.AmountRaised, proj.OrganizationID).Scan(&proj.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add new project.")
		fmt.Println(err)
		return nil
	}

	return &proj
}
