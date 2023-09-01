package handlers

import (
	"database/sql"
	"eshaanagg/lfx/database"
	"fmt"
)

func (client Client) GetAllParentOrgs() {
	rowsRs, err := client.Query("SELECT * FROM parentOrgs;")

	if err != nil {
		fmt.Println("[ERROR] GetAllParentOrgs query failed")
		fmt.Println(err)
		return
	}
	defer rowsRs.Close()

	orgs, err := parseResultSetToSlice(rowsRs)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(orgs)
}

func (client Client) CreateParentOrg(org *database.ParentOrg) {
	insertStmt :=
		`
        INSERT INTO parentOrgs (name, avatar) 
        VALUES($1, $2) 
        RETURNING id;
        `

	err := client.QueryRow(insertStmt, org.Name, org.Logo).Scan(&org.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add new parent organization")
		fmt.Println(err)
	}
}

func (client Client) ParentOrgExists(name string) bool {
	queryStmt :=
		`
        SELECT * FROM parentOrgs 
        WHERE name = $1;
        `
	rowsRs, err := client.Query(queryStmt, name)

	if err != nil {
		fmt.Println("[ERROR] ParentOrgExists query failed")
		fmt.Println(err)
		return true
	}
	defer rowsRs.Close()

	orgs, err := parseResultSetToSlice(rowsRs)
	if err != nil {
		fmt.Println(err)
		return true
	}

	return len(orgs) == 1
}

// Helper function to convert the resultset of a SELECT * query to slice of ParentOrg struct
func parseResultSetToSlice(rowsRs *sql.Rows) ([]database.ParentOrg, error) {
	// Creates placeholder
	orgs := make([]database.ParentOrg, 0)

	// we loop through the values of rows
	for rowsRs.Next() {
		org := database.ParentOrg{}
		err := rowsRs.Scan(&org.ID, &org.Name, &org.Name)
		if err != nil {
			fmt.Println("[ERROR] Can't save to ParentOrg struct")
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}
