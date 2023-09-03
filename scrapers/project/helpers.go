package project

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"
)

const dateTimeFormatString = "2006-01-02 15:04:05 +0000"
const maxOptions = float64(5)

func getTermFromMonth(month time.Month) string {
	term := "Term 1"
	if month >= 3 {
		term = "Term 2"
	}
	if month >= 6 {
		term = "Term 3"
	}
	if month >= 9 {
		term = "Term 1"
	}

	return term
}

func getProjectTerm(project projectResponse) (int, string) {
	limit := math.Min(float64(len(dateTimeFormatString)), float64(len(project.CreatedOn)))
	timeString := project.CreatedOn[0:int(limit)]

	createdDate, err := time.Parse(dateTimeFormatString, timeString)
	if err != nil {
		fmt.Printf("[ERROR] Can't parse date object from the provided string. Recieved: %s\n", timeString)
		fmt.Println(err)

		var year int
		fmt.Print("Please enter the year: ")
		fmt.Scan(&year)

		fmt.Print("Enter the term: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		term := trim(scanner.Text())

		return year, term
	}

	year := createdDate.Year()
	month := createdDate.Month()
	term := getTermFromMonth(month)

	if month >= 9 {
		year++
	}

	return year, term
}

func getMostProbableOrgs(allOrgs []string, project string) []string {

	data := map[string]int{}
	for _, org := range allOrgs {
		data[org] = 0
	}

	// Compare each organization anf project name word by word, and assign a score
	for _, projWord := range strings.Split(project, " ") {
		for _, org := range allOrgs {
			for _, orgWord := range strings.Split(org, " ") {
				if strings.EqualFold(orgWord, projWord) {
					data[org]++
				}
			}
		}
	}

	// Create a struct to sort by values
	var keyValuePairs []struct {
		Key   string
		Value int
	}
	for key, value := range data {
		keyValuePairs = append(keyValuePairs, struct {
			Key   string
			Value int
		}{key, value})
	}

	// Define a custom sorting function to sort by values
	sort.Slice(keyValuePairs, func(i, j int) bool {
		return keyValuePairs[i].Value < keyValuePairs[j].Value
	})

	orgs := make([]string, 0)
	for i := 0; float64(i) < math.Min(maxOptions, float64(len(keyValuePairs))); i++ {
		orgs = append(orgs, keyValuePairs[i].Key)
	}
	return orgs
}

func findInExistingOrgs(allOrgs []string, name string) string {
	for _, org := range allOrgs {
		if strings.Contains(strings.ToLower(name), strings.ToLower(org)) {
			return org
		}
	}

	return ""
}
