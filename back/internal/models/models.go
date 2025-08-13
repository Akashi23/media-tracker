package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeVideo MediaType = "video"
	MediaTypeBook  MediaType = "book"
	MediaTypeAnime MediaType = "anime"
	MediaTypeGame  MediaType = "game"
	MediaTypeTV    MediaType = "tv"
	MediaTypeMovie MediaType = "movie"
)

type Status string

const (
	StatusPlanned     Status = "planned"
	StatusInProgress  Status = "in_progress"
	StatusCompleted   Status = "completed"
	StatusOnHold      Status = "on_hold"
	StatusDropped     Status = "dropped"
)

type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type MediaItem struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	Type          MediaType  `json:"type" db:"type"`
	Title         string     `json:"title" db:"title"`
	OriginalTitle *string    `json:"original_title,omitempty" db:"original_title"`
	Year          *int       `json:"year,omitempty" db:"year"`
	CoverURL      *string    `json:"cover_url,omitempty" db:"cover_url"`
	Creators      JSONB      `json:"creators,omitempty" db:"creators"`
	Genres        []string   `json:"genres,omitempty" db:"genres"`
	Duration      *int       `json:"duration,omitempty" db:"duration"`
	Metadata      JSONB      `json:"metadata,omitempty" db:"metadata"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type Entry struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	MediaID   uuid.UUID  `json:"media_id" db:"media_id"`
	Status    Status     `json:"status" db:"status"`
	Rating    *float64   `json:"rating,omitempty" db:"rating"`
	ReviewMD  *string    `json:"review_md,omitempty" db:"review_md"`
	Progress  JSONB      `json:"progress,omitempty" db:"progress"`
	StartedAt *time.Time `json:"started_at,omitempty" db:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty" db:"finished_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
	Media     *MediaItem `json:"media,omitempty"`
}

type Collection struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	IsPublic  bool      `json:"is_public" db:"is_public"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Entries   []Entry   `json:"entries,omitempty"`
}

type ShareToken struct {
	Token     string     `json:"token" db:"token"`
	Kind      string     `json:"kind" db:"kind"`
	TargetID  uuid.UUID  `json:"target_id" db:"target_id"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" db:"expires_at"`
}

// JSONB type for PostgreSQL JSONB fields
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New("cannot scan non-string value into JSONB")
	}
	
	return json.Unmarshal(bytes, j)
}

// Request/Response DTOs
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type CreateEntryRequest struct {
	MediaID   uuid.UUID  `json:"media_id" binding:"required"`
	Status    Status     `json:"status" binding:"required"`
	Rating    *float64   `json:"rating,omitempty"`
	ReviewMD  *string    `json:"review_md,omitempty"`
	Progress  JSONB      `json:"progress,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}

type CreateMediaRequest struct {
	Type          MediaType `json:"type" binding:"required"`
	Title         string    `json:"title" binding:"required"`
	OriginalTitle *string   `json:"original_title,omitempty"`
	Year          *int      `json:"year,omitempty"`
	CoverURL      *string   `json:"cover_url,omitempty"`
	Creators      JSONB     `json:"creators,omitempty"`
	Genres        []string  `json:"genres,omitempty"`
	Duration      *int      `json:"duration,omitempty"`
	Metadata      JSONB     `json:"metadata,omitempty"`
}

type CreateCollectionRequest struct {
	Title    string     `json:"title" binding:"required"`
	IsPublic bool       `json:"is_public"`
	EntryIDs []uuid.UUID `json:"entry_ids,omitempty"`
}

type GuestSnapshotRequest struct {
	Entries []Entry `json:"entries" binding:"required"`
	Media   []MediaItem `json:"media" binding:"required"`
}

type MergeRequest struct {
	GuestEntries []Entry `json:"guest_entries" binding:"required"`
}
