package handlers

import (
	"database/sql"
	"eshaanagg/lfx/database"
	"fmt"
)

// The following function makes a database query to get info of all orgs
// and returns a slice containing the same
//
// method for: 	client:  object(databse instance)
// returns:    	orgs:    slice of ParentOrg objects, ParentOrg defined at ../models.go
func (client Client) GetAllParentOrgs() []database.ParentOrg {

	//get all orgs as sql.Rows object
	rowsRs, err := client.Query("SELECT * FROM parentOrgs;")

	if err != nil {
		fmt.Println("[ERROR] GetAllParentOrgs query failed")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}
	defer rowsRs.Close()

	// create slice from sql.Rows
	// this function is defined at line:{154} of this file
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

	//loop over orgs slice to create names slice
	for _, org := range orgs {
		names = append(names, org.Name)
	}

	return names
}


// The following function creates inserts new data into the database
//
// method for:  client object(databse instance)
// args:    	name, logo
// returns: 	org: ParentOrg object with the data inserted
func (client Client) CreateParentOrg(name string, logo string) *database.ParentOrg {
	// create placeholder object
	org := database.ParentOrg{ID: "0", Name: name, Logo: logo}

	insertStmt :=
		`
        INSERT INTO parentOrgs (name, logo) 
        VALUES($1, $2) 
        RETURNING id;
        `

	// insert the object/data into the database
	// Note: Scan is used only for the purposes of getting errors
	//       as QueryRow doesn't return error
	err := client.QueryRow(insertStmt, org.Name, org.Logo).Scan(&org.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add new parent organization.")
		fmt.Println(err)
		return nil
	}

	return &org
}


// Function to search for organisations by name
//
// method for:     client
// args:           name (search argument)
// returns:        orgs[0], first parentOrg object(row) in database with the
//                 given name
func (client Client) GetOrganizationByName(name string) *database.ParentOrg {
	// query into database
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

	//create parentOrg slice
	orgs, err := parseResultSetToSlice(rowsRs) // this function is defined at
											   // line: {154} of this file
	if err != nil {
		fmt.Println(err)
		return nil
	}

	// if no org exists for given name return nil
	if len(orgs) == 0 {
		return nil
	}

	//return first correspoing org
	return &orgs[0]
}

// The following function is mostly same as the function written above (GetOrganizationByName)
// with 1 difference, the arguments
//
// args: id
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

	orgs, err := parseResultSetToSlice(rowsRs) // this fn is defined below
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
		err := rowsRs.Scan(&org.ID, &org.Name, &org.Logo)
		if err != nil {
			fmt.Println("[ERROR] Can't save to ParentOrg struct")
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}
