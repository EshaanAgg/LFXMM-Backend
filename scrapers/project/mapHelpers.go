package project

import "sort"

func getKeysSortedByFrequency(frequencyMap map[string]int) []string {
	type frequencyItem struct {
		Key   string
		Value int
	}

	frequencyList := make([]frequencyItem, 0)

	for key, value := range frequencyMap {
		frequencyList = append(frequencyList, frequencyItem{key, value})
	}

	// Sort the slice by values in descending order
	sort.Slice(frequencyList, func(i, j int) bool {
		return frequencyList[i].Value > frequencyList[j].Value
	})

	// Get the most frequent values
	mostFrequentValues := make([]string, 0)
	for _, entry := range frequencyList {
		mostFrequentValues = append(mostFrequentValues, entry.Key)
	}

	return mostFrequentValues
}
