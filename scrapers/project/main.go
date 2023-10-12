package project

import (
	"eshaanagg/lfx/database/handlers"
	"fmt"
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

func UpdateSkillsForOrgs() {
	client := handlers.New()
	defer client.Close()

	orgs := client.GetAllParentOrgs()
	for _, org := range orgs {
		projects := client.GetProjectsByOrganization(org.ID)
		frequencyMap := make(map[string]int)

		for _, project := range projects {
			for _, skill := range project.Skills {
				if _, exists := frequencyMap[skill]; exists {
					frequencyMap[skill]++
				} else {
					frequencyMap[skill] = 1
				}
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
