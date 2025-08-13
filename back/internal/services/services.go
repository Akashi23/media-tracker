package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"media-tracker/internal/config"
	"media-tracker/internal/models"
	"media-tracker/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// AuthService
type AuthService struct {
	userRepo *repository.UserRepository
	redis    *redis.Client
	jwtConfig config.JWTConfig
}

func NewAuthService(userRepo *repository.UserRepository, redis *redis.Client, jwtConfig config.JWTConfig) *AuthService {
	return &AuthService{userRepo: userRepo, redis: redis, jwtConfig: jwtConfig}
}

func (s *AuthService) Login(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Create new user
			user = &models.User{
				ID:        uuid.New(),
				Email:     email,
				Name:      email, // Default to email
				CreatedAt: time.Now(),
			}
			if err := s.userRepo.Create(ctx, user); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(time.Duration(s.jwtConfig.Expiry) * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.jwtConfig.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) GetUser(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

// MediaService
type MediaService struct {
	mediaRepo *repository.MediaRepository
}

func NewMediaService(mediaRepo *repository.MediaRepository) *MediaService {
	return &MediaService{mediaRepo: mediaRepo}
}

func (s *MediaService) Create(ctx context.Context, req *models.CreateMediaRequest) (*models.MediaItem, error) {
	media := &models.MediaItem{
		ID:            uuid.New(),
		Type:          req.Type,
		Title:         req.Title,
		OriginalTitle: req.OriginalTitle,
		Year:          req.Year,
		CoverURL:      req.CoverURL,
		Creators:      req.Creators,
		Genres:        req.Genres,
		Duration:      req.Duration,
		Metadata:      req.Metadata,
		CreatedAt:     time.Now(),
	}

	if err := s.mediaRepo.Create(ctx, media); err != nil {
		return nil, err
	}

	return media, nil
}

func (s *MediaService) Search(ctx context.Context, query string, mediaType *models.MediaType) ([]*models.MediaItem, error) {
	return s.mediaRepo.Search(ctx, query, mediaType)
}

// EntryService
type EntryService struct {
	entryRepo *repository.EntryRepository
	mediaRepo *repository.MediaRepository
}

func NewEntryService(entryRepo *repository.EntryRepository, mediaRepo *repository.MediaRepository) *EntryService {
	return &EntryService{entryRepo: entryRepo, mediaRepo: mediaRepo}
}

func (s *EntryService) Create(ctx context.Context, userID uuid.UUID, req *models.CreateEntryRequest) (*models.Entry, error) {
	// Verify media exists
	media, err := s.mediaRepo.GetByID(ctx, req.MediaID)
	if err != nil {
		return nil, err
	}

	entry := &models.Entry{
		ID:         uuid.New(),
		UserID:     userID,
		MediaID:    req.MediaID,
		Status:     req.Status,
		Rating:     req.Rating,
		ReviewMD:   req.ReviewMD,
		Progress:   req.Progress,
		StartedAt:  req.StartedAt,
		FinishedAt: req.FinishedAt,
		UpdatedAt:  time.Now(),
		Media:      media,
	}

	if err := s.entryRepo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *EntryService) List(ctx context.Context, userID uuid.UUID, status *models.Status, mediaType *models.MediaType) ([]*models.Entry, error) {
	return s.entryRepo.ListByUser(ctx, userID, status, mediaType)
}

func (s *EntryService) Get(ctx context.Context, id uuid.UUID) (*models.Entry, error) {
	return s.entryRepo.GetByID(ctx, id)
}

func (s *EntryService) Update(ctx context.Context, id uuid.UUID, req *models.CreateEntryRequest) (*models.Entry, error) {
	entry, err := s.entryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	entry.Status = req.Status
	entry.Rating = req.Rating
	entry.ReviewMD = req.ReviewMD
	entry.Progress = req.Progress
	entry.StartedAt = req.StartedAt
	entry.FinishedAt = req.FinishedAt
	entry.UpdatedAt = time.Now()

	if err := s.entryRepo.Update(ctx, entry); err != nil {
		return nil, err
	}

	return entry, nil
}

func (s *EntryService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.entryRepo.Delete(ctx, id)
}

// CollectionService
type CollectionService struct {
	collectionRepo *repository.CollectionRepository
	entryRepo      *repository.EntryRepository
}

func NewCollectionService(collectionRepo *repository.CollectionRepository, entryRepo *repository.EntryRepository) *CollectionService {
	return &CollectionService{collectionRepo: collectionRepo, entryRepo: entryRepo}
}

func (s *CollectionService) Create(ctx context.Context, userID uuid.UUID, req *models.CreateCollectionRequest) (*models.Collection, error) {
	collection := &models.Collection{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     req.Title,
		IsPublic:  req.IsPublic,
		CreatedAt: time.Now(),
	}

	if err := s.collectionRepo.Create(ctx, collection); err != nil {
		return nil, err
	}

	return collection, nil
}

func (s *CollectionService) Get(ctx context.Context, id uuid.UUID) (*models.Collection, error) {
	return s.collectionRepo.GetByID(ctx, id)
}

// ShareService
type ShareService struct {
	shareRepo      *repository.ShareRepository
	collectionRepo *repository.CollectionRepository
	entryRepo      *repository.EntryRepository
}

func NewShareService(shareRepo *repository.ShareRepository, collectionRepo *repository.CollectionRepository, entryRepo *repository.EntryRepository) *ShareService {
	return &ShareService{shareRepo: shareRepo, collectionRepo: collectionRepo, entryRepo: entryRepo}
}

func (s *ShareService) CreateShareToken(ctx context.Context, kind string, targetID uuid.UUID) (*models.ShareToken, error) {
	token := generateToken()
	
	share := &models.ShareToken{
		Token:     token,
		Kind:      kind,
		TargetID:  targetID,
		CreatedAt: time.Now(),
		ExpiresAt: &[]time.Time{time.Now().AddDate(0, 1, 0)}[0], // 1 month
	}

	if err := s.shareRepo.Create(ctx, share); err != nil {
		return nil, err
	}

	return share, nil
}

func (s *ShareService) GetPublicShare(ctx context.Context, token string) (interface{}, error) {
	share, err := s.shareRepo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	switch share.Kind {
	case "collection":
		return s.collectionRepo.GetByID(ctx, share.TargetID)
	case "profile":
		// Return user's entries
		return s.entryRepo.ListByUser(ctx, share.TargetID, nil, nil)
	default:
		return nil, errors.New("unknown share kind")
	}
}

// GuestService
type GuestService struct {
	entryRepo *repository.EntryRepository
	mediaRepo *repository.MediaRepository
	shareRepo *repository.ShareRepository
}

func NewGuestService(entryRepo *repository.EntryRepository, mediaRepo *repository.MediaRepository, shareRepo *repository.ShareRepository) *GuestService {
	return &GuestService{entryRepo: entryRepo, mediaRepo: mediaRepo, shareRepo: shareRepo}
}

func (s *GuestService) CreateSnapshot(ctx context.Context, req *models.GuestSnapshotRequest) (*models.ShareToken, error) {
	// Create snapshot ID
	snapshotID := uuid.New()
	
	// Create share token for snapshot
	share := &models.ShareToken{
		Token:     generateToken(),
		Kind:      "snapshot",
		TargetID:  snapshotID,
		CreatedAt: time.Now(),
		ExpiresAt: &[]time.Time{time.Now().AddDate(0, 1, 0)}[0], // 1 month
	}

	if err := s.shareRepo.Create(ctx, share); err != nil {
		return nil, err
	}

	return share, nil
}

func (s *GuestService) MergeToAccount(ctx context.Context, userID uuid.UUID, guestEntries []models.Entry) error {
	// Simple merge: create new entries for user
	for _, guestEntry := range guestEntries {
		entry := &models.Entry{
			ID:         uuid.New(),
			UserID:     userID,
			MediaID:    guestEntry.MediaID,
			Status:     guestEntry.Status,
			Rating:     guestEntry.Rating,
			ReviewMD:   guestEntry.ReviewMD,
			Progress:   guestEntry.Progress,
			StartedAt:  guestEntry.StartedAt,
			FinishedAt: guestEntry.FinishedAt,
			UpdatedAt:  time.Now(),
		}

		if err := s.entryRepo.Create(ctx, entry); err != nil {
			return err
		}
	}

	return nil
}

// Utility functions
func generateToken() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
