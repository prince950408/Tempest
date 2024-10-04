package models

type Artist struct {
	Anv         string `json:"anv"`
	ID          int    `json:"id"`
	Join        string `json:"join"`
	Name        string `json:"name"`
	ResourceURL string `json:"resource_url"`
	Role        string `json:"role"`
	Tracks      string `json:"tracks"`
}

type Release struct {
	ID        int
	Title     string
	Artists   []Artist `json:"artists"`
	Genres    []string
	Styles    []string
	Year      int
	Status    string
	Thumb     string
	ArtistIDs []int
	GenreIDs  []int
	StyleIDs  []int
}

type PaginationInfo struct {
	PerPage int `json:"per_page"`
	Pages   int `json:"pages"`
	Page    int `json:"page"`
	Items   int `json:"items"`
	URLs    struct {
		Last string `json:"last"`
		Next string `json:"next"`
	} `json:"urls"`
}

type DiscogsResponse struct {
	Pagination PaginationInfo `json:"pagination"`
	Releases   []Release      `json:"releases"`
}

type QueryParams struct {
	ArtistID      int
	StyleID       int
	GenreID       int
	SortField     string
	SortDirection string
	Page          int
	Limit         int
	Offset        int
}

type SearchResponse struct {
	Releases       []Release      `json:"releases"`
	CountPerArtist map[string]int `json:"count_per_artist"`
	CountPerGenre  map[string]int `json:"count_per_genre"`
	CountPerStyle  map[string]int `json:"count_per_style"`
	Total          int            `json:"total"`
}

type Filter struct {
	Artists map[int]string `json:"artists"`
	Genres  map[int]string `json:"genres"`
	Styles  map[int]string `json:"styles"`
}
