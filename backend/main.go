package main

import (
	"log"
	"net/http"
	"os"

	"tempest_backend/api"
	"tempest_backend/database"

	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins, you can restrict this to specific domains
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

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

	fetchAtFirst := os.Getenv("FETCH_AT_FIRST")
	if fetchAtFirst == "true" {
		log.Println("Fetching data on app startup...")
		err = api.FetchAndStoreReleases()
		if err != nil {
			log.Fatalf("Failed to fetch releases on startup: %v", err)
		}
		log.Println("Initial data fetch completed.")
	}

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

	http.Handle("/search", corsMiddleware(http.HandlerFunc(api.SearchReleases)))
	http.Handle("/get-filter", corsMiddleware(http.HandlerFunc(api.FetchFilter)))

	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	select {}
}
