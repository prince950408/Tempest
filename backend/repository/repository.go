package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"tempest_backend/database"
	"tempest_backend/models"
)

func InsertArtist(artists []models.Artist) ([]int, error) {
	db := database.GetDB()
	artistIDs := []int{}

	for _, artist := range artists {
		var artist_id int
		err := db.QueryRow(`SELECT id FROM artists WHERE external_id = $1`, artist.ID).Scan(&artist_id)
		if err == sql.ErrNoRows {
			err = db.QueryRow(`INSERT INTO artists (external_id, name) VALUES ($1, $2) RETURNING id`, artist.ID, artist.Name).Scan(&artist_id)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		artistIDs = append(artistIDs, artist_id)
	}

	return artistIDs, nil
}

func InsertGenres(genres []string) ([]int, error) {
	db := database.GetDB()
	genreIDs := []int{}

	for _, genre := range genres {
		var genreID int
		err := db.QueryRow(`SELECT id FROM genres WHERE name = $1`, genre).Scan(&genreID)
		if err == sql.ErrNoRows {
			err = db.QueryRow(`INSERT INTO genres (name) VALUES ($1) RETURNING id`, genre).Scan(&genreID)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		genreIDs = append(genreIDs, genreID)
	}

	return genreIDs, nil
}

func InsertStyles(styles []string) ([]int, error) {
	db := database.GetDB()
	styleIDs := []int{}

	for _, style := range styles {
		var styleID int
		err := db.QueryRow(`SELECT id FROM styles WHERE name = $1`, style).Scan(&styleID)
		if err == sql.ErrNoRows {
			err = db.QueryRow(`INSERT INTO styles (name) VALUES ($1) RETURNING id`, style).Scan(&styleID)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}
		styleIDs = append(styleIDs, styleID)
	}

	return styleIDs, nil
}

func SaveRelease(release models.Release) error {
	db := database.GetDB()

	var releaseID int

	artistIDs, err := InsertArtist(release.Artists)
	if err != nil {
		return err
	}

	genreIDs, err := InsertGenres(release.Genres)
	if err != nil {
		return err
	}

	styleIDs, err := InsertStyles(release.Styles)
	if err != nil {
		return err
	}

	artistIDsJSON, err := json.Marshal(artistIDs)
	if err != nil {
		return err
	}

	genreIDsJSON, err := json.Marshal(genreIDs)
	if err != nil {
		return err
	}

	styleIDsJSON, err := json.Marshal(styleIDs)
	if err != nil {
		return err
	}

	err = db.QueryRow(`SELECT id FROM releases WHERE external_id = $1`, release.ID).Scan(&releaseID)
	if err == sql.ErrNoRows {
		err = db.QueryRow(`
			INSERT INTO releases (external_id, title, year, status, thumb, artist_ids, genre_ids, style_ids) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) 
			ON CONFLICT (external_id) DO NOTHING RETURNING id`,
			release.ID, release.Title, release.Year, release.Status, release.Thumb, artistIDsJSON, genreIDsJSON, styleIDsJSON).
			Scan(&releaseID)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func FetchReleases(params models.QueryParams) ([]models.Release, int, error) {
	db := database.GetDB()

	countQuery := `
	SELECT COUNT(DISTINCT r.id)
	FROM
		releases r
	WHERE 1=1
	`

	query := `
	SELECT
		r.id, r.title, r.year, r.status, r.thumb, r.artist_ids, r.genre_ids, r.style_ids
	FROM
		releases r
	WHERE 1=1
	`

	var conditions []string
	if params.ArtistID != 0 {
		conditions = append(conditions, fmt.Sprintf("r.artist_ids @> to_jsonb(ARRAY[%d])", params.ArtistID))
	}
	if params.GenreID != 0 {
		conditions = append(conditions, fmt.Sprintf("r.genre_ids @> to_jsonb(ARRAY[%d])", params.GenreID))
	}
	if params.StyleID != 0 {
		conditions = append(conditions, fmt.Sprintf("r.style_ids @> to_jsonb(ARRAY[%d])", params.StyleID))
	}

	if len(conditions) > 0 {
		conditionStr := " AND " + strings.Join(conditions, " AND ")
		countQuery += conditionStr
		query += conditionStr
	}

	var total int
	err := db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute count query: %v", err)
	}

	query += fmt.Sprintf(" ORDER BY r.%s %s", params.SortField, params.SortDirection)
	query += fmt.Sprintf(" LIMIT %d", params.Limit)
	query += fmt.Sprintf(" OFFSET %d", params.Offset)

	rows, err := db.Query(query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	var releases []models.Release

	for rows.Next() {
		var release models.Release
		var artistIDsRaw, genreIDsRaw, styleIDsRaw []byte

		err := rows.Scan(
			&release.ID,
			&release.Title,
			&release.Year,
			&release.Status,
			&release.Thumb,
			&artistIDsRaw,
			&genreIDsRaw,
			&styleIDsRaw,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan row: %v", err)
		}
		err = json.Unmarshal(artistIDsRaw, &release.ArtistIDs)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal artist_ids: %v", err)
		}

		err = json.Unmarshal(genreIDsRaw, &release.GenreIDs)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal genre_ids: %v", err)
		}

		err = json.Unmarshal(styleIDsRaw, &release.StyleIDs)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal style_ids: %v", err)
		}

		releases = append(releases, release)
	}

	if rows.Err() != nil {
		return nil, 0, fmt.Errorf("error during row iteration: %v", rows.Err())
	}
	return releases, total, nil
}

func GetReleaseCount(artistID int, genreID int, styleID int, field string) (map[string]int, error) {
	db := database.GetDB()

	query := `
	WITH filtered_releases AS (
		SELECT
			r.*,
			jsonb_array_elements_text(r.artist_ids) AS artist_id,
			jsonb_array_elements_text(r.genre_ids) AS genre_id,
			jsonb_array_elements_text(r.style_ids) AS style_id
		FROM
			releases r
		WHERE 1 = 1
	`

	var conditions []string
	if artistID != 0 {
		conditions = append(conditions, fmt.Sprintf("artist_ids @> to_jsonb(ARRAY[%d])", artistID))
	}
	if genreID != 0 {
		conditions = append(conditions, fmt.Sprintf("genre_ids @> to_jsonb(ARRAY[%d])", genreID))
	}
	if styleID != 0 {
		conditions = append(conditions, fmt.Sprintf("style_ids @> to_jsonb(ARRAY[%d])", styleID))
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += `
	)
	`

	var table, joinColumn string

	switch field {
	case "artist":
		table = "artists"
		joinColumn = "artist_id"
	case "genre":
		table = "genres"
		joinColumn = "genre_id"
	case "style":
		table = "styles"
		joinColumn = "style_id"
	default:
		return nil, fmt.Errorf("unsupported field type: %s", field)
	}

	query += fmt.Sprintf(`
		SELECT
			COALESCE(t.id::TEXT, '') AS entity_id,
			COUNT(DISTINCT r.id) AS release_count
		FROM
			filtered_releases r
		LEFT JOIN %s t ON t.id::TEXT = r.%s
		GROUP BY t.id
	`, table, joinColumn)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	result := make(map[string]int)

	for rows.Next() {
		var entityId string
		var releaseCount int

		err := rows.Scan(&entityId, &releaseCount)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		result[entityId] = releaseCount
	}

	return result, nil
}

func GetFilter(table string) (map[int]string, error) {
	db := database.GetDB()
	query := fmt.Sprintf("SELECT id, name FROM %s", table)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	result := make(map[int]string)
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}
		result[id] = name
	}

	return result, nil
}
