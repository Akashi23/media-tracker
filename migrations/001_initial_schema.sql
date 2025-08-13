-- Initial schema for media tracker
-- Run this migration to set up the database

-- Create custom types
CREATE TYPE media_type AS ENUM ('video', 'book', 'anime', 'game', 'tv', 'movie');

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    name TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Media items table
CREATE TABLE media_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    type media_type NOT NULL,
    title TEXT NOT NULL,
    original_title TEXT,
    year INTEGER,
    cover_url TEXT,
    creators JSONB,
    genres TEXT[],
    duration INTEGER,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- User entries table
CREATE TABLE entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    media_id UUID REFERENCES media_items(id) ON DELETE CASCADE,
    status TEXT CHECK (status IN ('planned', 'in_progress', 'completed', 'on_hold', 'dropped')),
    rating NUMERIC(3,1) CHECK (rating BETWEEN 0 AND 10),
    review_md TEXT,
    progress JSONB,
    started_at DATE,
    finished_at DATE,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, media_id)
);

-- Collections table
CREATE TABLE collections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Collection entries junction table
CREATE TABLE collection_entries (
    collection_id UUID REFERENCES collections(id) ON DELETE CASCADE,
    entry_id UUID REFERENCES entries(id) ON DELETE CASCADE,
    position INTEGER,
    PRIMARY KEY (collection_id, entry_id)
);

-- Share tokens table
CREATE TABLE share_tokens (
    token TEXT PRIMARY KEY,
    kind TEXT CHECK (kind IN ('collection', 'profile', 'snapshot')),
    target_id UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    expires_at TIMESTAMPTZ
);

-- Create indexes for performance
CREATE INDEX idx_entries_user_status ON entries(user_id, status);
CREATE INDEX idx_media_items_type_year ON media_items(type, year);
CREATE INDEX idx_share_tokens_expires ON share_tokens(expires_at);
CREATE INDEX idx_entries_updated_at ON entries(updated_at DESC);
CREATE INDEX idx_media_items_title ON media_items USING gin(to_tsvector('english', title));

-- Add default user for testing
INSERT INTO users (id, email, name, created_at) VALUES 
    ('550e8400-e29b-41d4-a716-446655440000', 'admin@example.com', 'Admin User', NOW());

-- Add some sample data for testing
INSERT INTO media_items (id, type, title, year, genres) VALUES 
    (gen_random_uuid(), 'movie', 'Mononoke Hime', 1997, ARRAY['Animation', 'Adventure', 'Fantasy']),
    (gen_random_uuid(), 'book', 'The Lord of the Rings', 1954, ARRAY['Fantasy', 'Adventure']),
    (gen_random_uuid(), 'anime', 'Attack on Titan', 2013, ARRAY['Action', 'Drama', 'Horror']),
    (gen_random_uuid(), 'game', 'The Witcher 3: Wild Hunt', 2015, ARRAY['RPG', 'Action', 'Adventure']),
    (gen_random_uuid(), 'tv', 'Breaking Bad', 2008, ARRAY['Crime', 'Drama', 'Thriller']);

-- Add sample entries for the default user
INSERT INTO entries (id, user_id, media_id, status, rating, review_md, started_at, finished_at, updated_at) 
SELECT 
    gen_random_uuid(),
    '550e8400-e29b-41d4-a716-446655440000',
    mi.id,
    'completed',
    9.0,
    'Great movie! Highly recommended.',
    '2024-01-01',
    '2024-01-02',
    NOW()
FROM media_items mi 
WHERE mi.title = 'Mononoke Hime'
LIMIT 1;
