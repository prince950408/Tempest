export interface Release {
    ID: number;
    Title: string;
    Year: number;
    Status: string;
    Thumb: string;
    ArtistIDs: number[];
    GenreIDs: number[];
    StyleIDs: number[];
}

export interface SearchResponse {
    releases: Release[];
    count_per_artist: { [key: string]: number };
    count_per_genre: { [key: string]: number };
    count_per_style: { [key: string]: number };
    total: number;
}

export interface FiltersProps {
    onFilterChange: (filters: { style_id: string; artist_id: string; genre_id: string }) => void;
    artists: { [key: string]: string };
    genres: { [key: string]: string };
    styles: { [key: string]: string };
    count_per_artist: { [key: string]: number };
    count_per_genre: { [key: string]: number };
    count_per_style: { [key: string]: number };
  }
  