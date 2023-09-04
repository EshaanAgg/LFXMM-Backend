package main

import (
	"eshaanagg/lfx/api"

	"github.com/joho/godotenv"
)

func main() {
	// Populate the environment variables
	godotenv.Load(".env")
	// Scrape()

	api.Start()
}

func Scrape() {
	// Generate all the project ids in a .csv file
	// projectIdScrapper.GenerateProjectIds()
	// project.Parse()
	// project.Merge()
}
