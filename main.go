package main

import (
	projectid "eshaanagg/lfx/scrapers/projectId"

	"github.com/joho/godotenv"
)

func main() {
	// Populate the environment variables
	godotenv.Load(".env")
	projectid.GenerateProjectIds()
}
