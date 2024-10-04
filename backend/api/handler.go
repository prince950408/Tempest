package api

import (
	"encoding/json"
	"log"
	"net/http"
	"tempest_backend/models"
	"tempest_backend/repository"
	"tempest_backend/utils"
)

func SearchReleases(w http.ResponseWriter, r *http.Request) {
	params := utils.ParseQueryParams(r)

	releases, total, err := repository.FetchReleases(params)
	if err != nil {
		log.Printf("Error fetching release: %v", err)
		http.Error(w, "Failed to fetch releases", http.StatusInternalServerError)
		return
	}

	countPerArtist, err := repository.GetReleaseCount(0, params.GenreID, params.StyleID, "artist")
	if err != nil {
		log.Printf("Error fetching count per artist: %v", err)
		http.Error(w, "Failed to fetch counts per artist", http.StatusInternalServerError)
		return
	}

	countPerGenre, err := repository.GetReleaseCount(params.ArtistID, 0, params.StyleID, "genre")
	if err != nil {
		log.Printf("Error fetching count per genre: %v", err)
		http.Error(w, "Failed to fetch counts per genre", http.StatusInternalServerError)
		return
	}

	countPerStyle, err := repository.GetReleaseCount(params.ArtistID, params.GenreID, 0, "style")
	if err != nil {
		log.Printf("Error fetching count per style: %v", err)
		http.Error(w, "Failed to fetch counts per style", http.StatusInternalServerError)
		return
	}

	response := models.SearchResponse{
		Releases:       releases,
		CountPerArtist: countPerArtist,
		CountPerGenre:  countPerGenre,
		CountPerStyle:  countPerStyle,
		Total:          total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func FetchFilter(w http.ResponseWriter, r *http.Request) {
	artists, err := repository.GetFilter("artists")
	if err != nil {
		log.Printf("Error fetching filter: %v", err)
		http.Error(w, "Failed to fetch filter", http.StatusInternalServerError)
		return
	}

	genres, err := repository.GetFilter("genres")
	if err != nil {
		log.Printf("Error fetching filter: %v", err)
		http.Error(w, "Failed to fetch filter", http.StatusInternalServerError)
		return
	}

	styles, err := repository.GetFilter("styles")
	if err != nil {
		log.Printf("Error fetching filter: %v", err)
		http.Error(w, "Failed to fetch filter", http.StatusInternalServerError)
		return
	}

	response := models.Filter{
		Artists: artists,
		Genres:  genres,
		Styles:  styles,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
