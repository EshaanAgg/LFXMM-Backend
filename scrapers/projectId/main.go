package projectid

import (
	"fmt"
)

func GetProjectIds() []string {
	limit := 50
	start := 0
	allIds := make([]string, 0)

	for {
		ids, err := makeRequest(start, limit)
		countIds := len(ids)

		if err != nil {
			fmt.Printf("[ERROR] The request with start = %d failed.\n", start)
			fmt.Println("Error: ", err)
		} else {
			fmt.Printf("[SUCCESS] ScrapeProjectID start = %d count = %d\n", start, countIds)
			allIds = append(allIds, ids...)
		}

		if countIds < limit {
			break
		}
		start += 25
	}

	return allIds
}
