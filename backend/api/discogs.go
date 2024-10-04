package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"tempest_backend/models"
	"tempest_backend/repository"
)

func FetchReleasesPage(labelID int, page, perPage int) (models.DiscogsResponse, error) {
	var response models.DiscogsResponse

	url := fmt.Sprintf("https://api.discogs.com/labels/%s/releases?token=%s&page=%d&per_page=%d",
		os.Getenv("LABEL_ID"),
		os.Getenv("DISCOGS_TOKEN"),
		page,
		perPage)

	resp, err := http.Get(url)
	if err != nil {
		return response, fmt.Errorf("failed to fetch releases: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, fmt.Errorf("failed to decode response: %v", err)
	}

	return response, nil
}

func fetchReleaseDetails(releaseID int) (models.Release, error) {
	var response models.Release

	url := fmt.Sprintf("https://api.discogs.com/releases/%d?token=%s", releaseID, os.Getenv("DISCOGS_TOKEN"))

	log.Printf("Fetching: %s", url)
	resp, err := http.Get(url)
	if err != nil {
		return response, fmt.Errorf("failed to fetch releases details: %v", err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return response, fmt.Errorf("failed to decode response: %v", err)
	}

	return response, nil
}

func ReleasesGenerator(labelID, perPage int) chan models.Release {
	releaseChan := make(chan models.Release)

	go func() {
		defer close(releaseChan)

		page := 1
		for {
			response, err := FetchReleasesPage(labelID, page, perPage)
			if err != nil {
				log.Printf("Error fetching page %d: %v", page, err)
				return
			}

			for _, release := range response.Releases {
				releaseChan <- release
			}

			if page >= response.Pagination.Pages {
				break
			}

			page++
			time.Sleep(1 * time.Second) // To avoid hitting API rate limits
		}
	}()

	return releaseChan
}

func FetchAndStoreReleases() error {
	labelIDStr := os.Getenv("LABEL_ID")
	perPageStr := os.Getenv("RESULTS_PER_PAGE")

	labelID, err := strconv.Atoi(labelIDStr)
	if err != nil {
		log.Fatalf("Invalid LABEL_ID value: %v", err)
		return err
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		log.Fatalf("Invalid RESULTS_PER_PAGE value: %v", err)
		return err
	}

	log.Println("Fetching and storing releases...")

	releaseChan := ReleasesGenerator(labelID, perPage)

	for release := range releaseChan {
		releaseDetail, err := fetchReleaseDetails(release.ID)
		if err != nil {
			log.Printf("Error fetching details for release %d: %v", releaseDetail.ID, err)
			continue
		}

		err = repository.SaveRelease(releaseDetail)
		if err != nil {
			log.Printf("Error saving release: %v", err)
			continue
		}

		time.Sleep(1 * time.Second)
	}

	log.Println("All releases successfully fetched and stored.")
	return nil
}
