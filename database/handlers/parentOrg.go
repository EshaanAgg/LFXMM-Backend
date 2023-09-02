package handlers

import (
	"database/sql"
	"eshaanagg/lfx/database"
	"fmt"
)

func (client Client) GetAllParentOrgs() []database.ParentOrg {
	rowsRs, err := client.Query("SELECT * FROM parentOrgs;")

	if err != nil {
		fmt.Println("[ERROR] GetAllParentOrgs query failed")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}
	defer rowsRs.Close()

	orgs, err := parseResultSetToSlice(rowsRs)
	if err != nil {
		fmt.Println("[ERROR] Can't convert to result set in GetAllParentOrgs function.")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}

	return orgs
}

func (client Client) GetAllOrgNames() []string {
	orgs := client.GetAllParentOrgs()
	names := make([]string, 0)

	for _, org := range orgs {
		names = append(names, org.Name)
	}

	return names
}

func (client Client) CreateParentOrg(name string, logo string) *database.ParentOrg {
	org := database.ParentOrg{ID: "0", Name: name, Logo: logo}

	insertStmt :=
		`
        INSERT INTO parentOrgs (name, logo) 
        VALUES($1, $2) 
        RETURNING id;
        `

	err := client.QueryRow(insertStmt, org.Name, org.Logo).Scan(&org.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add new parent organization")
		fmt.Println(err)
		return nil
	}

	return &org
}

func (client Client) GetOrganizationByName(name string) *database.ParentOrg {
	queryStmt :=
		`
        SELECT * FROM parentOrgs 
        WHERE name = $1;
        `
	rowsRs, err := client.Query(queryStmt, name)

	if err != nil {
		fmt.Println("[ERROR] GetOrganizationByName query failed")
		fmt.Println(err)
		return nil
	}
	defer rowsRs.Close()

	orgs, err := parseResultSetToSlice(rowsRs)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(orgs) == 0 {
		return nil
	}

	return &orgs[0]
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
