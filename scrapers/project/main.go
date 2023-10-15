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
			UpdateUniqueSkillsMap(project)
		}
		fmt.Println()
	}

	UpdateSkillsForOrgs()

	//Save unique skills to the database
	client := handlers.New()
	defer client.Close()

	addedSkill := client.SaveUniqueSkillsMaptoDb(uniqueSkillsMap)
	if addedSkill != nil {
		fmt.Println("[ERROR] Can't save this skills map to database.")
	} else {
		fmt.Println("[SUCCESS] Skills Map added.")
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
				frequencyMap[skill]++
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

var uniqueSkillsMap = make(map[string]int)

// UpdateUniqueSkillsMap updates the unique skills map.
func UpdateUniqueSkillsMap(projRes *projectResponse) {
	for _, skill := range projRes.ApprenticeNeeds.Skills {
		skill = trim(skill)
		if isSingleWord(skill) {
			if freq, exists := uniqueSkillsMap[skill]; exists {
				// Increment the frequency
				freq++
				uniqueSkillsMap[skill] = freq
			} else {
				// Initialize the frequency to 1
				uniqueSkillsMap[skill] = 1
			}
		}
	}
}
