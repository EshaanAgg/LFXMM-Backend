package main

import (
	"eshaanagg/lfx/api"
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

	api.Start()
}

func Scrape() {
	// Generate all the project ids in a .csv file
	// projectIdScrapper.GenerateProjectIds()
	// project.Parse()
	// project.Merge()
}
