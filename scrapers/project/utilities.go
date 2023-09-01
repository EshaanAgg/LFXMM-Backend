package project

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const apiURL = "https://api.mentorship.lfx.linuxfoundation.org/project/%s"

func makeRequest(projectID string) (*projectResponse, error) {
	// Make the API request
	res, err := http.Get(fmt.Sprintf(apiURL, projectID))

	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code recieved: %v", res.StatusCode)
	}

	// Convert the HTTP response to bytes
	bodyBtyes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON and store the parsed data
	var response projectResponse
	err = json.Unmarshal(bodyBtyes, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func parseCSV() ([]string, error) {
	pwd, _ := os.Getwd()
	f, err := os.Create(pwd + "/scrapers/assets/projectIDs.csv")

	if err != nil {
		fmt.Printf("[ERROR] Can't open projectIDs.csv to read data.\n")
		fmt.Println("Error: ", err)
		return nil, err
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()

	if err != nil {
		fmt.Printf("[ERROR] Can't parse data to projectIDs.csv. Please ensure the file is a well formatted CSV.\n")
		fmt.Println("Error: ", err)
		return nil, err
	}

	if len(data) != 1 {
		fmt.Printf("[ERROR] The CSV has multiple rows. It should only have one single row of comma seperated project ids. Please ensure the same. Exiting.\n")
		return nil, errors.New("the file should have only one record")
	}
	ids := data[0]

	return ids, nil
}
