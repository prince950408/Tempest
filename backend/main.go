package main

import (
	"log"
	"net/http"

	"tempest_backend/api"
	"tempest_backend/database"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err = database.InitDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	defer database.CloseDB()

	log.Println("Fetching data on app startup...")
	err = api.FetchAndStoreReleases()
	if err != nil {
		log.Fatalf("Failed to fetch releases on startup: %v", err)
	}
	log.Println("Initial data fetch completed.")

	c := cron.New()

	c.AddFunc("@daily", func() {
		log.Println("Scheduled daily fetch triggered.")
		err := api.FetchAndStoreReleases()
		if err != nil {
			log.Printf("Error fetching releases: %v", err)
		} else {
			log.Println("Daily data fetch completed.")
		}
	})

	c.Start()

	// Set up HTTP routes for querying data
	http.HandleFunc("/search", api.SearchReleases)
	http.HandleFunc("/get-filter", api.FetchFilter)

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	select {}
}
