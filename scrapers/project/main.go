package project

import (
	"eshaanagg/lfx/database/handlers"
	"fmt"
	"strings"
	"sync"
)

func Parse() {
	projectIds, err := parseCSV()

	if err != nil {
		fmt.Println("[ERROR] Cancelling the parsing.")
		return
	}

	for ind, id := range projectIds {
		fmt.Printf("[INFO] Processing Mentorship Project %d\n", ind+1)
		project, err := makeRequest(id)
		if err != nil {
			fmt.Printf("[ERROR] Request failed for project %s. Try the same again later.\n", id)
			fmt.Println(err)
		} else {
			addToDatabase(project)
		}
		fmt.Println()
	}
}

// Sets the skill attribute for all the oraganizations
// The skills are converted to `lowercase` explicitly to as to ensure uniformity
func UpdateSkillsForOrgs() {
	client := handlers.New()
	defer client.Close()

	orgs := client.GetAllParentOrgs()
	for _, org := range orgs {
		projects := client.GetProjectsByOrganization(org.ID)
		frequencyMap := make(map[string]int)

		for _, project := range projects {
			for _, skill := range project.Skills {
				frequencyMap[strings.ToLower(skill)]++
			}
		}

		skills := getKeysSortedByFrequency(frequencyMap)

		if len(skills) == 0 {
			continue
		}

		skillInterface := make([]interface{}, 0)
		for _, skill := range skills {
			skillInterface = append(skillInterface, skill)
		}

		err := client.SetSkillsForOrg(org.ID, skillInterface)
		if err != nil {
			fmt.Println("[ERROR] There was an error in updating the skills for the parent organization.")
			fmt.Println(err)
		}
	}
}

// Converts the skills of all the projects to `lowercase` explicitly to as to ensure uniformity
func LowercaseSkillsForProjects() {
	client := handlers.New()
	defer client.Close()

	orgs := client.GetAllParentOrgs()
	for _, org := range orgs {
		projects := client.GetProjectsByOrganization(org.ID)

		for _, project := range projects {

			skills := make([]string, 0)
			for _, skill := range project.Skills {
				skills = append(skills, strings.ToLower(skill))
			}

			if len(skills) == 0 {
				continue
			}

			err := client.SetSkillsForProject(project.ID, skills)
			if err != nil {
				fmt.Println("[ERROR] There was an error in updating the skills for the parent organization.")
				fmt.Println(err)
			}
		}
	}
}

// Adds the description to organizations with the help of OpenAI
func UpdateOrganizationDescriptions() {
	client := handlers.New()
	defer client.Close()

	var wg sync.WaitGroup

	orgs := client.GetAllParentOrgs()

	for _, org := range orgs {
		wg.Add(1)

		// Create a IIFC to populate the data for the organization
		go func(orgID string, orgName string) {
			defer wg.Done()

			desc, err := getAIDescription(orgName)
			if err != nil {
				fmt.Println("[ERROR] There was an error in getting the description.")
				fmt.Println(err)
			}

			client := handlers.New()
			defer client.Close()

			err = client.SetDescForOrg(orgID, desc)
			if err != nil {
				fmt.Println("[ERROR] There was an error in saving the description to database.")
				fmt.Println(err)
			}
		}(org.ID, org.Name)
	}

	wg.Wait()
}
