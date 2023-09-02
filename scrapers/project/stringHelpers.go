package project

import (
	"regexp"
	"strings"
)

func isSingleWord(input string) bool {
	// Define a regular expression pattern for a single word
	pattern := "^[A-Za-z]+$"
	regex, err := regexp.Compile(pattern)

	if err != nil {
		return false // Invalid pattern
	}

	return regex.MatchString(input)
}

func trim(input string) string {
	return strings.Trim(input, " \t\n\r")
}
