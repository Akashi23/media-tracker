package repository

import (
	"context"
	"database/sql"
	"media-tracker/internal/models"
	"strings"

	"github.com/google/uuid"
)

// UserRepository
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name, user.CreatedAt)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, name, created_at FROM users WHERE email = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, name, created_at FROM users WHERE id = $1`
	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// MediaRepository
type MediaRepository struct {
	db *sql.DB
}

func NewMediaRepository(db *sql.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(ctx context.Context, media *models.MediaItem) error {
	query := `INSERT INTO media_items (id, type, title, original_title, year, cover_url, creators, genres, duration, metadata, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.ExecContext(ctx, query, media.ID, media.Type, media.Title, media.OriginalTitle, media.Year,
		media.CoverURL, media.Creators, media.Genres, media.Duration, media.Metadata, media.CreatedAt)
	return err
}

func (r *MediaRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.MediaItem, error) {
	query := `SELECT id, type, title, original_title, year, cover_url, creators, genres, duration, metadata, created_at 
			  FROM media_items WHERE id = $1`
	media := &models.MediaItem{}
	var genresBytes []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(&media.ID, &media.Type, &media.Title, &media.OriginalTitle,
		&media.Year, &media.CoverURL, &media.Creators, &genresBytes, &media.Duration, &media.Metadata, &media.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Convert PostgreSQL array bytes to string slice
	if genresBytes != nil {
		genresStr := string(genresBytes)
		// Remove curly braces and split by comma
		if len(genresStr) >= 2 && genresStr[0] == '{' && genresStr[len(genresStr)-1] == '}' {
			genresStr = genresStr[1 : len(genresStr)-1]
			if genresStr != "" {
				media.Genres = strings.Split(genresStr, ",")
				// Trim whitespace from each genre
				for i, genre := range media.Genres {
					media.Genres[i] = strings.TrimSpace(genre)
				}
			}
		}
	}

	return media, nil
}

func (r *MediaRepository) Search(ctx context.Context, query string, mediaType *models.MediaType) ([]*models.MediaItem, error) {
	baseQuery := `SELECT id, type, title, original_title, year, cover_url, creators, genres, duration, metadata, created_at 
				  FROM media_items WHERE title ILIKE $1`
	args := []interface{}{"%" + query + "%"}

	if mediaType != nil {
		baseQuery += " AND type = $2"
		args = append(args, *mediaType)
	}

	baseQuery += " ORDER BY title LIMIT 20"

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*models.MediaItem
	for rows.Next() {
		media := &models.MediaItem{}
		var genresBytes []byte
		err := rows.Scan(&media.ID, &media.Type, &media.Title, &media.OriginalTitle, &media.Year,
			&media.CoverURL, &media.Creators, &genresBytes, &media.Duration, &media.Metadata, &media.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Convert PostgreSQL array bytes to string slice
		if genresBytes != nil {
			genresStr := string(genresBytes)
			// Remove curly braces and split by comma
			if len(genresStr) >= 2 && genresStr[0] == '{' && genresStr[len(genresStr)-1] == '}' {
				genresStr = genresStr[1 : len(genresStr)-1]
				if genresStr != "" {
					media.Genres = strings.Split(genresStr, ",")
					// Trim whitespace from each genre
					for i, genre := range media.Genres {
						media.Genres[i] = strings.TrimSpace(genre)
					}
				}
			}
		}

		items = append(items, media)
	}
	return items, nil
}

// EntryRepository
type EntryRepository struct {
	db *sql.DB
}

func NewEntryRepository(db *sql.DB) *EntryRepository {
	return &EntryRepository{db: db}
}

func (r *EntryRepository) Create(ctx context.Context, entry *models.Entry) error {
	query := `INSERT INTO entries (id, user_id, media_id, status, rating, review_md, progress, started_at, finished_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err := r.db.ExecContext(ctx, query, entry.ID, entry.UserID, entry.MediaID, entry.Status, entry.Rating,
		entry.ReviewMD, entry.Progress, entry.StartedAt, entry.FinishedAt, entry.UpdatedAt)
	return err
}

func (r *EntryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Entry, error) {
	query := `SELECT e.id, e.user_id, e.media_id, e.status, e.rating, e.review_md, e.progress, e.started_at, e.finished_at, e.updated_at,
			  m.id, m.type, m.title, m.original_title, m.year, m.cover_url, m.creators, m.genres, m.duration, m.metadata, m.created_at
			  FROM entries e 
			  JOIN media_items m ON e.media_id = m.id 
			  WHERE e.id = $1`

	entry := &models.Entry{Media: &models.MediaItem{}}
	var genresBytes []byte
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&entry.ID, &entry.UserID, &entry.MediaID, &entry.Status, &entry.Rating, &entry.ReviewMD, &entry.Progress,
		&entry.StartedAt, &entry.FinishedAt, &entry.UpdatedAt,
		&entry.Media.ID, &entry.Media.Type, &entry.Media.Title, &entry.Media.OriginalTitle, &entry.Media.Year,
		&entry.Media.CoverURL, &entry.Media.Creators, &genresBytes, &entry.Media.Duration, &entry.Media.Metadata, &entry.Media.CreatedAt)
	if err != nil {
		return nil, err
	}

	// Convert PostgreSQL array bytes to string slice
	if genresBytes != nil {
		genresStr := string(genresBytes)
		// Remove curly braces and split by comma
		if len(genresStr) >= 2 && genresStr[0] == '{' && genresStr[len(genresStr)-1] == '}' {
			genresStr = genresStr[1 : len(genresStr)-1]
			if genresStr != "" {
				entry.Media.Genres = strings.Split(genresStr, ",")
				// Trim whitespace from each genre
				for i, genre := range entry.Media.Genres {
					entry.Media.Genres[i] = strings.TrimSpace(genre)
				}
			}
		}
	}

	return entry, nil
}

func (r *EntryRepository) ListByUser(ctx context.Context, userID uuid.UUID, status *models.Status, mediaType *models.MediaType) ([]*models.Entry, error) {
	query := `SELECT e.id, e.user_id, e.media_id, e.status, e.rating, e.review_md, e.progress, e.started_at, e.finished_at, e.updated_at,
			  m.id, m.type, m.title, m.original_title, m.year, m.cover_url, m.creators, m.genres, m.duration, m.metadata, m.created_at
			  FROM entries e 
			  JOIN media_items m ON e.media_id = m.id 
			  WHERE e.user_id = $1`
	args := []interface{}{userID}

	if status != nil {
		query += " AND e.status = $2"
		args = append(args, *status)
	}

	if mediaType != nil {
		query += " AND m.type = $3"
		args = append(args, *mediaType)
	}

	query += " ORDER BY e.updated_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*models.Entry
	for rows.Next() {
		entry := &models.Entry{Media: &models.MediaItem{}}
		var genresBytes []byte
		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.MediaID, &entry.Status, &entry.Rating, &entry.ReviewMD, &entry.Progress,
			&entry.StartedAt, &entry.FinishedAt, &entry.UpdatedAt,
			&entry.Media.ID, &entry.Media.Type, &entry.Media.Title, &entry.Media.OriginalTitle, &entry.Media.Year,
			&entry.Media.CoverURL, &entry.Media.Creators, &genresBytes, &entry.Media.Duration, &entry.Media.Metadata, &entry.Media.CreatedAt)
		if err != nil {
			return nil, err
		}

		// Convert PostgreSQL array bytes to string slice
		if genresBytes != nil {
			genresStr := string(genresBytes)
			// Remove curly braces and split by comma
			if len(genresStr) >= 2 && genresStr[0] == '{' && genresStr[len(genresStr)-1] == '}' {
				genresStr = genresStr[1 : len(genresStr)-1]
				if genresStr != "" {
					entry.Media.Genres = strings.Split(genresStr, ",")
					// Trim whitespace from each genre
					for i, genre := range entry.Media.Genres {
						entry.Media.Genres[i] = strings.TrimSpace(genre)
					}
				}
			}
		}

		entries = append(entries, entry)
	}
	return entries, nil
}

func (r *EntryRepository) Update(ctx context.Context, entry *models.Entry) error {
	query := `UPDATE entries SET status = $1, rating = $2, review_md = $3, progress = $4, started_at = $5, finished_at = $6, updated_at = $7 
			  WHERE id = $8`
	_, err := r.db.ExecContext(ctx, query, entry.Status, entry.Rating, entry.ReviewMD, entry.Progress,
		entry.StartedAt, entry.FinishedAt, entry.UpdatedAt, entry.ID)
	return err
}

func (r *EntryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM entries WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// CollectionRepository
type CollectionRepository struct {
	db *sql.DB
}

func NewCollectionRepository(db *sql.DB) *CollectionRepository {
	return &CollectionRepository{db: db}
}

func (r *CollectionRepository) Create(ctx context.Context, collection *models.Collection) error {
	query := `INSERT INTO collections (id, user_id, title, is_public, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, collection.ID, collection.UserID, collection.Title, collection.IsPublic, collection.CreatedAt)
	return err
}

func (r *CollectionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	query := `SELECT id, user_id, title, is_public, created_at FROM collections WHERE id = $1`
	collection := &models.Collection{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&collection.ID, &collection.UserID, &collection.Title, &collection.IsPublic, &collection.CreatedAt)
	if err != nil {
		return nil, err
	}
	return collection, nil
}

// ShareRepository
type ShareRepository struct {
	db *sql.DB
}

func NewShareRepository(db *sql.DB) *ShareRepository {
	return &ShareRepository{db: db}
}

func (r *ShareRepository) Create(ctx context.Context, share *models.ShareToken) error {
	query := `INSERT INTO share_tokens (token, kind, target_id, created_at, expires_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, share.Token, share.Kind, share.TargetID, share.CreatedAt, share.ExpiresAt)
	return err
}

func (r *ShareRepository) GetByToken(ctx context.Context, token string) (*models.ShareToken, error) {
	query := `SELECT token, kind, target_id, created_at, expires_at FROM share_tokens WHERE token = $1`
	share := &models.ShareToken{}
	err := r.db.QueryRowContext(ctx, query, token).Scan(&share.Token, &share.Kind, &share.TargetID, &share.CreatedAt, &share.ExpiresAt)
	if err != nil {
		return nil, err
	}
	return share, nil
}
