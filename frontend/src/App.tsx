import React, { useEffect, useState } from 'react';
import Filters from './Components/Filter';
import { fetchFilteredData, fetchFilters } from './api/api';
import { SearchResponse } from './interfaces/interfaces';

const App: React.FC = () => {
  const [filteredData, setFilteredData] = useState<SearchResponse>({
    releases: [],
    count_per_artist: {},
    count_per_genre: {},
    count_per_style: {},
    total: 0,
  });

  const [artists, setArtists] = useState<{ [key: number]: string }>({});
  const [genres, setGenres] = useState<{ [key: number]: string }>({});
  const [styles, setStyles] = useState<{ [key: number]: string }>({});
  const [page, setPage] = useState<number>(1);
  const [filter, setFilter] = useState<{ style_id: string; artist_id: string; genre_id: string }>({
    style_id: '',
    artist_id: '',
    genre_id: ''
  });
  const [limit, setLimit] = useState<number>(10);
  const [sortConfig, setSortConfig] = useState<{ key: string; direction: 'asc' | 'desc' }>({
    key: '',
    direction: 'asc'
  });

  const handleFilterChange = (filters: { style_id: string; artist_id: string; genre_id: string }) => {
    setFilter(filters);
    setPage(1);
    fetchFilteredData(filters, page, limit, sortConfig.key, sortConfig.direction)
      .then((response) => setFilteredData(response.data))
      .catch((error) => console.error('Error fetching filtered data:', error));
    
  };

  useEffect(() => {
    fetchFilters()
      .then((response) => {
        setArtists(response.data['artists']);
        setGenres(response.data['genres']);
        setStyles(response.data['styles']);
      })
      .catch((error) => {
        console.error('Error fetching filter data:', error);
      });
  }, []);

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
    fetchFilteredData(filter, newPage, limit, sortConfig.key, sortConfig.direction)
      .then((response) => setFilteredData(response.data))
      .catch((error) => console.error('Error fetching paginated data:', error));
  };

  const handleSort = (key: string) => {
    const direction = sortConfig.key === key && sortConfig.direction === 'asc' ? 'desc' : 'asc';
    setSortConfig({ key, direction });
    fetchFilteredData(filter, page, limit, key, direction)
      .then((response) => setFilteredData(response.data))
      .catch((error) => console.error('Error fetching sorted data:', error));
  };


  return (
    <div className="container mx-auto p-4">
      <h1 className="text-3xl font-bold text-center mb-8">Discog Data Filter</h1>
      <Filters 
        onFilterChange={handleFilterChange} 
        artists={artists} genres={genres} 
        styles={styles} 
        count_per_artist={filteredData.count_per_artist}
        count_per_genre={filteredData.count_per_genre}
        count_per_style={filteredData.count_per_style}
      />
      <div className="mt-8">
        <h2 className="text-2xl font-semibold mb-4">Filtered Data : {filteredData.total}</h2>
        {filteredData.releases && filteredData.releases.length > 0 ? (
          <div>
            <div className="flex justify-center gap-8 items-center mt-6">
              <button
                disabled={page === 1}
                onClick={() => handlePageChange(page - 1)}
                className={`px-4 py-2 bg-blue-500 text-white rounded ${page === 1 && 'opacity-50 cursor-not-allowed'}`}
              >
                Previous
              </button>
              <span className="text-gray-700">Page {page}</span>
              <button
                disabled={page * limit >= filteredData.total}
                onClick={() => handlePageChange(page + 1)}
                className={`px-4 py-2 bg-blue-500 text-white rounded ${
                  page * limit >= filteredData.total && 'opacity-50 cursor-not-allowed'
                }`}
              >
                Next
              </button>
            </div>
            <table className="min-w-full bg-white">
              <thead>
                <tr>
                  <th className='w-[10%]'>Image</th>
                  <th onClick={() => handleSort('title')} className="cursor-pointer w-[10%]">
                    Title {sortConfig.key === 'title' ? (sortConfig.direction === 'asc' ? '↑' : '↓') : '↑↓'}
                  </th>
                  <th onClick={() => handleSort('year')} className="cursor-pointer w-[10%]">
                    Year {sortConfig.key === 'year' ? (sortConfig.direction === 'asc' ? '↑' : '↓') : '↑↓'}
                  </th>
                  <th onClick={() => handleSort('status')} className="cursor-pointer w-[10%]">
                    Status {sortConfig.key === 'status' ? (sortConfig.direction === 'asc' ? '↑' : '↓') : '↑↓'}
                  </th>
                  <th className="cursor-pointer w-[20%]">
                    Artists
                  </th>
                  <th className="cursor-pointer w-[20%]">
                    Genres
                  </th>
                  <th className="cursor-pointe w-[20%]">
                    Styles
                  </th>
                </tr>
              </thead>
              <tbody>
              {filteredData.releases.map((item) => (
                <tr key={item.ID}>
                  <td>
                  <img
                    src={item.Thumb}
                    alt={item.Title}
                    className="w-20 h-20 rounded-lg mx-auto object-cover"
                  />
                  </td>
                  <td className='text-center'>{item.Title}</td>
                  <td className='text-center'>{item.Year}</td>
                  <td className='text-center'>{item.Status}</td>
                  <td className='text-center'>{item.ArtistIDs?.map((id: number) => artists[id]).join(', ')}</td>
                  <td className='text-center'>{item.GenreIDs?.map((id: number) => genres[id]).join(', ')}</td>
                  <td className='text-center'>{item.StyleIDs?.map((id: number) => styles[id]).join(', ')}</td>
                </tr>
              ))}
              </tbody>
            </table>
            <div className="flex justify-center gap-8 items-center mt-6">
              <button
                disabled={page === 1}
                onClick={() => handlePageChange(page - 1)}
                className={`px-4 py-2 bg-blue-500 text-white rounded ${page === 1 && 'opacity-50 cursor-not-allowed'}`}
              >
                Previous
              </button>
              <span className="text-gray-700">Page {page}</span>
              <button
                disabled={page * limit >= filteredData.total}
                onClick={() => handlePageChange(page + 1)}
                className={`px-4 py-2 bg-blue-500 text-white rounded ${
                  page * limit >= filteredData.total && 'opacity-50 cursor-not-allowed'
                }`}
              >
                Next
              </button>
            </div>
          </div>
        ) : (
          <p className="text-gray-600">No results found.</p>
        )}
      </div>
    </div>
  );
};

export default App;
