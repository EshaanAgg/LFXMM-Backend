package projectid

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiURL = "https://api.mentorship.lfx.linuxfoundation.org/projects/cache/paginate?from=%d&size=%d&sortby=updatedStamp"

func makeRequest(start int, limit int) ([]string, error) {

	// Make the API request
	res, err := http.Get(fmt.Sprintf(apiURL, start, limit))

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
	var response apiResponse
	err = json.Unmarshal(bodyBtyes, &response)
	if err != nil {
		return nil, err
	}

	projectIds := make([]string, 0)
	for _, project := range response.Hits.Hits {
		projectIds = append(projectIds, project.ID)
	}

	return projectIds, nil
}
