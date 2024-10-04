import axios from 'axios';

const BASE_URL = process.env.REACT_APP_API_BASE_URL || 'http://localhost:8080'

export const fetchFilteredData = (
    filters: { 
        style_id: string,
        artist_id: string, 
        genre_id: string,
    },
    page: number,
    limit: number,
    sort: string,
    direction: string,
) => {
    return axios.get(`${BASE_URL}/search`, {
        params: { ...filters, page, limit, sort, direction },
    });
};

export const fetchFilters = () => {
    return axios.get(`${BASE_URL}/get-filter`);
};