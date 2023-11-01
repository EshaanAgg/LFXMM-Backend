package main

import (
	"eshaanagg/lfx/api"
	"eshaanagg/lfx/scrapers/project"
	projectid "eshaanagg/lfx/scrapers/projectId"
	"fmt"

	"github.com/joho/godotenv"
)

func main() {
	// Populate the environment variables
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Cannot load the environment variables")
		fmt.Println(err)
	}
	// Scrape()
	api.Start()
}

func Scrape() {
	// Generate all the project ids in a .csv file
	projectid.GenerateProjectIds()
	project.Merge()
	project.Parse()
}
