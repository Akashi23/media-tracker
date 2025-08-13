-- Seed data script for media tracker
-- Run this after the initial migration to add more test data

-- Add more test users
INSERT INTO users (id, email, name, created_at) VALUES 
    ('550e8400-e29b-41d4-a716-446655440001', 'john@example.com', 'John Doe', NOW()),
    ('550e8400-e29b-41d4-a716-446655440002', 'jane@example.com', 'Jane Smith', NOW()),
    ('550e8400-e29b-41d4-a716-446655440003', 'bob@example.com', 'Bob Johnson', NOW())
ON CONFLICT (email) DO NOTHING;

-- Add more media items
INSERT INTO media_items (id, type, title, year, genres, creators, metadata) VALUES 
    (gen_random_uuid(), 'movie', 'Spirited Away', 2001, ARRAY['Animation', 'Adventure', 'Fantasy'], 
     '{"director": "Hayao Miyazaki", "studio": "Studio Ghibli"}', 
     '{"runtime": 125, "language": "Japanese"}'),
    
    (gen_random_uuid(), 'book', '1984', 1949, ARRAY['Dystopian', 'Political Fiction'], 
     '{"author": "George Orwell"}', 
     '{"pages": 328, "isbn": "978-0451524935"}'),
    
    (gen_random_uuid(), 'anime', 'Death Note', 2006, ARRAY['Thriller', 'Psychological', 'Supernatural'], 
     '{"director": "Tetsurō Araki", "studio": "Madhouse"}', 
     '{"episodes": 37, "status": "completed"}'),
    
    (gen_random_uuid(), 'game', 'Red Dead Redemption 2', 2018, ARRAY['Action', 'Adventure', 'Western'], 
     '{"developer": "Rockstar Games", "publisher": "Rockstar Games"}', 
     '{"platforms": ["PS4", "Xbox One", "PC"], "genre": "Open World"}'),
    
    (gen_random_uuid(), 'tv', 'Game of Thrones', 2011, ARRAY['Fantasy', 'Drama', 'Adventure'], 
     '{"creators": ["David Benioff", "D.B. Weiss"], "network": "HBO"}', 
     '{"seasons": 8, "episodes": 73}')
ON CONFLICT DO NOTHING;

-- Add sample entries for different users
INSERT INTO entries (id, user_id, media_id, status, rating, review_md, started_at, finished_at, updated_at) 
SELECT 
    gen_random_uuid(),
    '550e8400-e29b-41d4-a716-446655440001',
    mi.id,
    'completed',
    8.5,
    'Amazing story and characters. The world-building is incredible.',
    '2024-02-01',
    '2024-02-15',
    NOW()
FROM media_items mi 
WHERE mi.title = 'Game of Thrones'
LIMIT 1;

INSERT INTO entries (id, user_id, media_id, status, rating, review_md, started_at, finished_at, updated_at) 
SELECT 
    gen_random_uuid(),
    '550e8400-e29b-41d4-a716-446655440002',
    mi.id,
    'in_progress',
    7.0,
    'Interesting concept, still reading...',
    '2024-03-01',
    NULL,
    NOW()
FROM media_items mi 
WHERE mi.title = '1984'
LIMIT 1;

INSERT INTO entries (id, user_id, media_id, status, rating, review_md, started_at, finished_at, updated_at) 
SELECT 
    gen_random_uuid(),
    '550e8400-e29b-41d4-a716-446655440003',
    mi.id,
    'planned',
    NULL,
    NULL,
    NULL,
    NULL,
    NOW()
FROM media_items mi 
WHERE mi.title = 'Death Note'
LIMIT 1;

-- Create a sample collection
INSERT INTO collections (id, user_id, title, is_public, created_at) VALUES 
    (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440000', 'My Favorite Movies', true, NOW());

-- Add entries to the collection
INSERT INTO collection_entries (collection_id, entry_id, position)
SELECT 
    c.id,
    e.id,
    1
FROM collections c, entries e
WHERE c.title = 'My Favorite Movies' 
AND e.user_id = '550e8400-e29b-41d4-a716-446655440000'
AND e.media_id IN (SELECT id FROM media_items WHERE title = 'Mononoke Hime')
LIMIT 1;
