package project

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"eshaanagg/lfx/database"
	"eshaanagg/lfx/database/handlers"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const apiURL = "https://api.mentorship.lfx.linuxfoundation.org/projects/%s"

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
	f, err := os.Open(pwd + "/scrapers/assets/projectIDs.csv")

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

func addToDatabase(projRes *projectResponse) {
	project := database.Project{}

	// One-to-one mapped fields
	project.LFXProjectID = projRes.ProjectID
	project.Name = trim(projRes.Name)
	project.Description = trim(projRes.Description)
	project.AmountRaised = projRes.AmountRaised
	project.Repository = projRes.Repository
	project.Website = projRes.Website

	// Fields that can be mapped via simple interpolation
	project.Industry = strings.Split(projRes.Description, "/")

	// Map skills that are single word
	project.Skills = make([]string, 0)
	for _, skill := range projRes.ApprenticeNeeds.Skills {
		if isSingleWord(trim(skill)) {
			project.Skills = append(project.Skills, trim(skill))
		}
	}

	// Map program term and year
	year, term := getProjectTerm(*projRes)
	if term != nil {
		project.ProgramYear = *year
		project.ProgramTerm = *term
	} else {
		fmt.Printf("[ERROR] Project Term and Year could not be parsed for project %s. It will be populated by standard zero values.\n", projRes.ProjectID)
		project.ProgramYear = 0
		project.ProgramTerm = "Uncategorized"
	}

	project.OrganizationID = getOrganizationID(projRes)

	client := handlers.New()
	proj := client.CreateProject(project)
	if proj == nil {
		fmt.Println("[ERROR] Can't save this project to database.")
	} else {
		fmt.Println("[SUCCESS] Project added.")
	}
}

func getOrganizationID(proj *projectResponse) string {
	name := proj.Name
	var orgName string

	if strings.Contains(name, ":") {
		orgName = trim(strings.Split(name, ":")[0])
	} else if strings.Contains(name, "-") {
		orgName = trim(strings.Split(name, "-")[0])
	} else {
		client := handlers.New()
		allOrgs := client.GetAllOrgNames()

		fmt.Printf("Project Name: %s. \n", name)
		fmt.Println("No organization can be found for this project. The most probable organizations for the same can be: ")

		probableOrgs := getMostProbableOrgs(allOrgs, name)
		for i, org := range probableOrgs {
			fmt.Printf("%d: %s\n", i+1, org)
		}

		// Read the user's input for the organization name.
		fmt.Println("Enter the name/index for the organization: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		sentence := trim(scanner.Text())

		index, err := strconv.Atoi(sentence)
		if err != nil || index >= len(probableOrgs) {
			orgName = sentence
		} else {
			orgName = sentence
		}
	}

	return getOrganization(orgName, proj)
}

func getOrganization(orgName string, proj *projectResponse) string {
	client := handlers.New()
	org := client.GetOrganizationByName(orgName)

	if org != nil {
		return org.ID
	}

	fmt.Printf("[INFO] Creating a new organization: %s \n", orgName)
	newOrg := client.CreateParentOrg(orgName, proj.LogoURL)

	if newOrg == nil {
		fmt.Println("[ERROR] New organization generation failed. Providing the organization ID as 0.")
		return "0"
	}

	return newOrg.ID
}
