package project

import (
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
