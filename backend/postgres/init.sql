-- Artists
CREATE TABLE IF NOT EXISTS artists (
    id SERIAL PRIMARY KEY,
    external_id INT UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL
);

-- Genres
CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Styles
CREATE TABLE IF NOT EXISTS styles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL
);

-- Releases
CREATE TABLE IF NOT EXISTS releases (
    id SERIAL PRIMARY KEY,
    external_id INT UNIQUE, 
    title VARCHAR(255),
    year INT,
    status VARCHAR(50),
    thumb VARCHAR(255),
    artist_ids JSONB DEFAULT '[]'::jsonb,
    genre_ids JSONB DEFAULT '[]'::jsonb,
    style_ids JSONB DEFAULT '[]'::jsonb
);
