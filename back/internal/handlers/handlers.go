package handlers

import (
	"fmt"
	"net/http"

	"media-tracker/internal/models"
	"media-tracker/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AuthHandler
type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.Login(c.Request.Context(), req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Simple logout - client should remove token
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.authService.GetUser(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// MediaHandler
type MediaHandler struct {
	mediaService *services.MediaService
}

func NewMediaHandler(mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService}
}

func (h *MediaHandler) Create(c *gin.Context) {
	var req models.CreateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media, err := h.mediaService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, media)
}

func (h *MediaHandler) Update(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter required"})
		return
	}

	var req models.UpdateMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	media, err := h.mediaService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, media)
}

func (h *MediaHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter required"})
		return
	}

	var mediaType *models.MediaType
	if typeStr := c.Query("type"); typeStr != "" {
		mt := models.MediaType(typeStr)
		mediaType = &mt
	}

	results, err := h.mediaService.Search(c.Request.Context(), query, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}

// EntryHandler
type EntryHandler struct {
	entryService *services.EntryService
	mediaService *services.MediaService
}

func NewEntryHandler(entryService *services.EntryService, mediaService *services.MediaService) *EntryHandler {
	return &EntryHandler{entryService: entryService, mediaService: mediaService}
}

func (h *EntryHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var status *models.Status
	if statusStr := c.Query("status"); statusStr != "" {
		s := models.Status(statusStr)
		status = &s
	}

	var mediaType *models.MediaType
	if typeStr := c.Query("type"); typeStr != "" {
		mt := models.MediaType(typeStr)
		mediaType = &mt
	}

	entries, err := h.entryService.List(c.Request.Context(), userID.(uuid.UUID), status, mediaType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entries)
}

func (h *EntryHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.entryService.Create(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, entry)
}

func (h *EntryHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	entry, err := h.entryService.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Entry not found"})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var req models.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	entry, err := h.entryService.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, entry)
}

func (h *EntryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.entryService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Entry deleted successfully"})
}

func (h *EntryHandler) Sync(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Entries []models.SyncEntryRequest `json:"entries" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Process each entry for sync
	var syncedEntries []models.Entry
	var errors []string

	for _, syncEntry := range req.Entries {
		// First, ensure media exists
		var mediaID uuid.UUID

		// Check if media already exists by title and type
		existingMedia, err := h.mediaService.Search(c.Request.Context(), syncEntry.Media.Title, &syncEntry.Media.Type)
		if err != nil || len(existingMedia) == 0 {
			// Create new media
			media, err := h.mediaService.Create(c.Request.Context(), &syncEntry.Media)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Error creating media %s: %v", syncEntry.Media.Title, err))
				continue
			}
			mediaID = media.ID
		} else {
			// Use existing media
			mediaID = existingMedia[0].ID
		}

		// Now create the entry
		entryReq := models.CreateEntryRequest{
			MediaID:    mediaID,
			Status:     syncEntry.Status,
			Rating:     syncEntry.Rating,
			ReviewMD:   syncEntry.ReviewMD,
			Progress:   syncEntry.Progress,
			StartedAt:  syncEntry.StartedAt,
			FinishedAt: syncEntry.FinishedAt,
		}

		// Check if entry already exists for this user and media
		existingEntries, err := h.entryService.ListByUserAndMedia(c.Request.Context(), userID.(uuid.UUID), mediaID)
		if err != nil {
			// Log the error but continue processing
			errors = append(errors, fmt.Sprintf("Error checking existing entries for media %s: %v", syncEntry.Media.Title, err))

			// Try to create new entry
			entry, createErr := h.entryService.Create(c.Request.Context(), userID.(uuid.UUID), &entryReq)
			if createErr != nil {
				errors = append(errors, fmt.Sprintf("Error creating entry for media %s: %v", syncEntry.Media.Title, createErr))
				continue
			}
			syncedEntries = append(syncedEntries, *entry)
			continue
		}

		if len(existingEntries) > 0 {
			// Update existing entry
			existingEntry := existingEntries[0]
			entry, err := h.entryService.Update(c.Request.Context(), existingEntry.ID, &entryReq)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Error updating entry %s: %v", existingEntry.ID, err))
				continue
			}
			syncedEntries = append(syncedEntries, *entry)
		} else {
			// Create new entry
			entry, err := h.entryService.Create(c.Request.Context(), userID.(uuid.UUID), &entryReq)
			if err != nil {
				errors = append(errors, fmt.Sprintf("Error creating entry for media %s: %v", syncEntry.Media.Title, err))
				continue
			}
			syncedEntries = append(syncedEntries, *entry)
		}
	}

	response := gin.H{
		"message": "Entries synced successfully",
		"entries": syncedEntries,
		"count":   len(syncedEntries),
	}

	if len(errors) > 0 {
		response["errors"] = errors
		response["message"] = fmt.Sprintf("Synced %d entries with %d errors", len(syncedEntries), len(errors))
	}

	c.JSON(http.StatusOK, response)
}

// CollectionHandler
type CollectionHandler struct {
	collectionService *services.CollectionService
	shareService      *services.ShareService
}

func NewCollectionHandler(collectionService *services.CollectionService, shareService *services.ShareService) *CollectionHandler {
	return &CollectionHandler{collectionService: collectionService, shareService: shareService}
}

func (h *CollectionHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	collections, err := h.collectionService.List(c.Request.Context(), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, collections)
}

func (h *CollectionHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection, err := h.collectionService.Create(c.Request.Context(), userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, collection)
}

func (h *CollectionHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	collection, err := h.collectionService.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Collection not found"})
		return
	}

	c.JSON(http.StatusOK, collection)
}

func (h *CollectionHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	var req models.CreateCollectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	collection, err := h.collectionService.Update(c.Request.Context(), id, userID.(uuid.UUID), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, collection)
}

func (h *CollectionHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid collection ID"})
		return
	}

	err = h.collectionService.Delete(c.Request.Context(), id, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted"})
}

func (h *CollectionHandler) CreateShare(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	share, err := h.shareService.CreateShareToken(c.Request.Context(), "collection", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"share_url": "/s/" + share.Token})
}

// ShareHandler
type ShareHandler struct {
	shareService *services.ShareService
}

func NewShareHandler(shareService *services.ShareService) *ShareHandler {
	return &ShareHandler{shareService: shareService}
}

func (h *ShareHandler) GetPublicShare(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token required"})
		return
	}

	data, err := h.shareService.GetPublicShare(c.Request.Context(), token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// GuestHandler
type GuestHandler struct {
	guestService *services.GuestService
}

func NewGuestHandler(guestService *services.GuestService) *GuestHandler {
	return &GuestHandler{guestService: guestService}
}

func (h *GuestHandler) CreateSnapshot(c *gin.Context) {
	var req models.GuestSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	share, err := h.guestService.CreateSnapshot(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"share_url": "/s/" + share.Token})
}

func (h *GuestHandler) MergeToAccount(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req models.MergeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.guestService.MergeToAccount(c.Request.Context(), userID.(uuid.UUID), req.GuestEntries); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Guest data merged successfully"})
}
