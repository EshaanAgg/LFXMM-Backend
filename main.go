package main

import (
	"eshaanagg/lfx/database/handlers"
	projectIdScrapper "eshaanagg/lfx/scrapers/projectId"

	"github.com/joho/godotenv"
)

func main() {
	// Populate the environment variables
	godotenv.Load(".env")

	client := handlers.New()
	client.GetAllParentOrgs()
}

func Scrape() {
	// Generate all the project ids in a .csv file
	projectIdScrapper.GenerateProjectIds()
}
