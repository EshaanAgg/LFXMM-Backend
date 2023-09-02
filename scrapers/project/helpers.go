package project

import (
	"fmt"
	"time"
)

const dateTimeFormatString = "2006-01-02 15:04:05 +0000"

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

func getProjectTerm(project projectResponse) (*int, *string) {
	timeString := project.CreatedOn

	createdDate, err := time.Parse(dateTimeFormatString, timeString)
	if err != nil {
		fmt.Printf("[ERROR] Can't parse date object from the provided string. Recieved: %s\n", timeString)
		fmt.Println(err)
		return nil, nil
	}

	year := createdDate.Year()
	month := createdDate.Month()
	term := getTermFromMonth(month)

	if month >= 9 {
		year++
	}

	return &year, &term
}
