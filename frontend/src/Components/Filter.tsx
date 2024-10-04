import React, { useState, useEffect } from 'react';
import { FiltersProps } from './../interfaces/interfaces'

const Filters: React.FC<FiltersProps> = ({ 
  onFilterChange, 
  artists, 
  genres, 
  styles, 
  count_per_artist, 
  count_per_genre, 
  count_per_style 
}) => {
  const [selectedStyle, setSelectedStyle] = useState<string>('');
  const [selectedArtist, setSelectedArtist] = useState<string>('');
  const [selectedGenre, setSelectedGenre] = useState<string>('');

  useEffect(() => {
    onFilterChange({
      style_id: selectedStyle,
      artist_id: selectedArtist,
      genre_id: selectedGenre,
    });
  }, [selectedStyle, selectedArtist, selectedGenre]);

  return (
    <div className="flex items-center gap-8">
      {/* Artist Filter */}
      <div className="flex items-center space-x-2">
        <label className="font-medium text-gray-700">Artist:</label>
        <select
          className="border border-gray-300 rounded px-3 py-2 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={selectedArtist}
          onChange={(e) => setSelectedArtist(e.target.value)}
        >
          <option value=""> </option>
          {Object.entries(artists).map(([id, name]) => (
            <option key={id + name} value={id}>
              {name} : {count_per_artist[id] ?? 0}
            </option>
          ))}
        </select>
      </div>

      {/* Genre Filter */}
      <div className="flex items-center space-x-2">
        <label className="font-medium text-gray-700">Genre:</label>
        <select
          className="border border-gray-300 rounded px-3 py-2 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={selectedGenre}
          onChange={(e) => setSelectedGenre(e.target.value)}
        >
          <option value=""> </option>
          {Object.entries(genres).map(([id, name]) => (
            <option key={id + name} value={id}>
              {name} : {count_per_genre[id] ?? 0}
            </option>
          ))}
        </select>
      </div>

      {/* Style Filter */}
      <div className="flex items-center space-x-2">
        <label className="font-medium text-gray-700">Style:</label>
        <select
          className="border border-gray-300 rounded px-3 py-2 text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          value={selectedStyle}
          onChange={(e) => setSelectedStyle(e.target.value)}
        >
          <option value=""></option>
          {Object.entries(styles).map(([id, name]) => (
            <option key={id + name} value={id}>
              {name} : {count_per_style[id] ?? 0}
            </option>
          ))}
        </select>
      </div>
    </div>
  );
};

export default Filters;
