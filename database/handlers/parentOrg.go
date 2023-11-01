package handlers

import (
	"database/sql"
	"eshaanagg/lfx/database"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

/*
 * Returns all organizations.
 */
func (client Client) GetAllParentOrgs() []database.ParentOrg {

	// Get all orgs as sql.Rows object.
	rowsRs, err := client.Query("SELECT * FROM parentOrgs ORDER BY name")

	if err != nil {
		fmt.Println("[ERROR] GetAllParentOrgs query failed")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}
	defer rowsRs.Close()

	orgs, err := parseAsParentOrgSlice(rowsRs)
	if err != nil {
		fmt.Println("[ERROR] Can't convert to result set in GetAllParentOrgs function.")
		fmt.Println(err)
		return make([]database.ParentOrg, 0)
	}

	return orgs
}

/*
 * Returns all the organization names.
 */
func (client Client) GetAllOrgNames() []string {
	orgs := client.GetAllParentOrgs()
	names := make([]string, 0)

	for _, org := range orgs {
		names = append(names, org.Name)
	}

	return names
}

/*
 * Inserts a new organization into the database
 * Skills are not populated manually. They are calculated from the projects that are conducted under an organization.
 */
func (client Client) CreateParentOrg(name string, logo string) *database.ParentOrg {
	// Create a placeholder object.
	org := database.ParentOrg{
		ID:          "0",
		Name:        name,
		Logo:        logo,
		Description: "",
	}

	updateStmt :=
		`
        INSERT INTO parentOrgs (name, logo) 
        VALUES($1, $2) 
        RETURNING id;
        `

	err := client.QueryRow(updateStmt, org.Name, org.Logo).Scan(&org.ID)

	if err != nil {
		fmt.Println("[ERROR] Can't add a new parent organization.")
		fmt.Println(err)
		return nil
	}

	return &org
}

/*
 * Function to search for organizations by name.
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
	orgs, err := parseAsParentOrgSlice(rowsRs) // This function is defined in this file.
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
 * Returns an organization with the provided ID
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

	orgs, err := parseAsParentOrgSlice(rowsRs) // This function is defined below.
	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(orgs) == 0 {
		return nil
	}

	return &orgs[0]
}

/*
 * Used to update the skills list for an organization
 */
func (client Client) SetSkillsForOrg(id string, skills []interface{}) error {
	params := make([]string, 0, len(skills))
	for i := range skills {
		params = append(params, fmt.Sprintf("$%v", i+1))
	}

	updateStmt := fmt.Sprintf(`
		UPDATE parentOrgs
		SET skills = ARRAY[%s]
		WHERE id = $%v
		RETURNING id;
		`,
		strings.Join(params, ", "),
		len(skills)+1,
	)
	// We append the id to the skills array as well so that the same can be destructed while executing the query
	skills = append(skills, id)

	_, err := client.Exec(
		updateStmt,
		skills...,
	)

	if err != nil {
		fmt.Printf("[ERROR] Can't update the skills for the organization %v.\n", id)
		fmt.Println(err)
		return nil
	}

	return nil
}

/*
 * Used to update the description an organization
 */
func (client Client) SetDescForOrg(id string, description string) error {
	updateStmt := `
		UPDATE parentOrgs
		SET description = $1
		WHERE id = $2
		RETURNING id;
		`

	_, err := client.Exec(
		updateStmt,
		description,
		id,
	)

	if err != nil {
		fmt.Printf("[ERROR] Can't update the description for the organization %v.\n", id)
		fmt.Println(err)
		return nil
	}

	return nil
}

func (client Client) SaveUniqueSkillsMaptoDb(skillsMap map[string]int) error {
	insertStmt := `
		INSERT INTO uniqueSkills (skill, frequency)
		VALUES ($1, $2)
		ON CONFLICT (skill) DO UPDATE SET frequency = EXCLUDED.frequency
	`

	stmt, err := client.Prepare(insertStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for skill, frequency := range skillsMap {
		_, err := stmt.Exec(skill, frequency)
		if err != nil {
			return err
		}
	}

	return nil
}

func (client Client) GetAllSkills() ([]database.UniqueSkills, error) {
	query := `
		SELECT skill, frequency 
		FROM uniqueSkills
		ORDER BY frequency DESC
	`

	rows, err := client.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allSkills []database.UniqueSkills

	for rows.Next() {
		var skill database.UniqueSkills
		err := rows.Scan(&skill.Skill, &skill.Frequency)
		if err != nil {
			return nil, err
		}
		allSkills = append(allSkills, skill)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return allSkills, nil
}

// Helper function to convert the resultset of a SELECT * query to a slice of ParentOrg struct.
func parseAsParentOrgSlice(rowsRs *sql.Rows) ([]database.ParentOrg, error) {
	// Create a placeholder.
	orgs := make([]database.ParentOrg, 0)

	// Loop through the values of rows.
	for rowsRs.Next() {
		org := database.ParentOrg{}
		err := rowsRs.Scan(
			&org.ID,
			&org.Name,
			&org.Logo,
			&org.Description,
			pq.Array(&org.Skills),
		)

		if err != nil {
			fmt.Println("[ERROR] Can't save to ParentOrg struct")
			return nil, err
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}
