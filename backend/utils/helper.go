package utils

import (
	"net/http"
	"strconv"
	"tempest_backend/models"
)

func ParseQueryParams(r *http.Request) models.QueryParams {

	artistID, err := strconv.Atoi(r.URL.Query().Get("artist_id"))
	if err != nil || artistID < 1 {
		artistID = 0
	}

	styleID, err := strconv.Atoi(r.URL.Query().Get("style_id"))
	if err != nil || styleID < 1 {
		styleID = 0
	}

	genreID, err := strconv.Atoi(r.URL.Query().Get("genre_id"))
	if err != nil || genreID < 1 {
		genreID = 0
	}

	sortField := r.URL.Query().Get("sort")
	if sortField == "" {
		sortField = "id"
	}

	sortDirection := r.URL.Query().Get("direction")
	if sortDirection == "" {
		sortDirection = "asc"
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	return models.QueryParams{
		ArtistID:      artistID,
		StyleID:       styleID,
		GenreID:       genreID,
		SortField:     sortField,
		SortDirection: sortDirection,
		Page:          page,
		Limit:         limit,
		Offset:        offset,
	}
}
