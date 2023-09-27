package handlers

import (
	"database/sql"
	"eshaanagg/lfx/database"
	"fmt"
)

/*
 * The following function makes a database query to get information of all
 * organizations and returns a slice containing the same.
 *
 * Method for:  Client (object - database instance)
 * Returns:     orgs (slice of ParentOrg objects, ParentOrg defined at ../models.go)
 */
func (client Client) GetAllParentOrgs() []database.ParentOrg {

	// Get all orgs as sql.Rows object.
	rowsRs, err := client.Query("SELECT * FROM parentOrgs;")

	if err != nil {
		fmt.Println("[ERROR] GetAllParentOrgs query failed")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}
	defer rowsRs.Close()

	/* Create slice from sql.Rows.
	 * This function is defined in this file.
	 */
	 orgs, err := parseResultSetToSlice(rowsRs)
	if err != nil {
		fmt.Println("[ERROR] Can't convert to result set in GetAllParentOrgs function.")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}

	return orgs
}

/*
 * This function gets all organization names.
 *
 * Method for: Client
 * Returns:    names (slice of organization names)
 */
func (client Client) GetAllOrgNames() []string {
	orgs := client.GetAllParentOrgs()
	names := make([]string, 0)

	// Loop over orgs slice to create names slice.
	for _, org := range orgs {
		names = append(names, org.Name)
	}

	return names
}

/*
 * The following function inserts new data into the database.
 *
 * Method for:  Client (object - database instance)
 * Args:        name, logo
 * Returns:     org (ParentOrg object with the data inserted)
 */
func (client Client) CreateParentOrg(name string, logo string) *database.ParentOrg {
	// Create a placeholder object.
	org := database.ParentOrg{ID: "0", Name: name, Logo: logo}

	insertStmt :=
		`
        INSERT INTO parentOrgs (name, logo) 
        VALUES($1, $2) 
        RETURNING id;
        `

	/* Insert the object/data into the database.
	 * Note: Scan is used only for the purposes of getting errors
	 *       as QueryRow doesn't return an error.
	 */
	 err := client.QueryRow(insertStmt, org.Name, org.Logo).Scan(&org.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add a new parent organization.")
		fmt.Println(err)
		return nil
	}

	return &org
}

/*
 * Function to search for organizations by name.
 *
 * Method for:  Client
 * Args:        name (search argument)
 * Returns:     orgs[0] (first ParentOrg object(row) in the database with the given name)
 */
func (client Client) GetOrganizationByName(name string) *database.ParentOrg {
	// Query into the database.
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

	// Create a ParentOrg slice.
	orgs, err := parseResultSetToSlice(rowsRs) // This function is defined in this file.
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// If no organization exists for the given name, return nil.
	if len(orgs) == 0 {
		return nil
	}

	// Return the first corresponding org.
	return &orgs[0]
}

/*
 * This function is mostly the same as the function written above (GetOrganizationByName)
 * with one difference, the arguments.
 *
 * Args: id
 */
func (client Client) GetOrganizationByID(id string) *database.ParentOrg {
	queryStmt :=
		`
        SELECT * FROM parentOrgs 
        WHERE id = $1;
        `

	rowsRs, err := client.Query(queryStmt, id)
	if err != nil {
		fmt.Println("[ERROR] GetOrganizationByID query failed")
		fmt.Println(err)
		return nil
	}
	defer rowsRs.Close()

	orgs, err := parseResultSetToSlice(rowsRs) // This function is defined below.
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(orgs) == 0 {
		return nil
	}

	return &orgs[0]
}


 // Helper function to convert the resultset of a SELECT * query to a slice of ParentOrg struct. 
func parseResultSetToSlice(rowsRs *sql.Rows) ([]database.ParentOrg, error) {
	// Create a placeholder.
	orgs := make([]database.ParentOrg, 0)

	// Loop through the values of rows.
	for rowsRs.Next() {
		org := database.ParentOrg{}
		err := rowsRs.Scan(&org.ID, &org.Name, &org.Logo)
		if err != nil {
			fmt.Println("[ERROR] Can't save to ParentOrg struct")
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}
