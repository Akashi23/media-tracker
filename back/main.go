package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"media-tracker/internal/config"
	"media-tracker/internal/database"
	"media-tracker/internal/handlers"
	"media-tracker/internal/middleware"
	"media-tracker/internal/repository"
	"media-tracker/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize logger
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if gin.Mode() == gin.DebugMode {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Initialize database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := database.NewRedisConnection(cfg.Redis)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	mediaRepo := repository.NewMediaRepository(db)
	entryRepo := repository.NewEntryRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	shareRepo := repository.NewShareRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, redisClient, cfg.JWT)
	mediaService := services.NewMediaService(mediaRepo)
	entryService := services.NewEntryService(entryRepo, mediaRepo)
	collectionService := services.NewCollectionService(collectionRepo, entryRepo)
	shareService := services.NewShareService(shareRepo, collectionRepo, entryRepo)
	guestService := services.NewGuestService(entryRepo, mediaRepo, shareRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	mediaHandler := handlers.NewMediaHandler(mediaService)
	entryHandler := handlers.NewEntryHandler(entryService, mediaService)
	collectionHandler := handlers.NewCollectionHandler(collectionService, shareService)
	shareHandler := handlers.NewShareHandler(shareService)
	guestHandler := handlers.NewGuestHandler(guestService)

	// Setup router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(&logger))
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "timestamp": time.Now()})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.GET("/me", middleware.Auth(cfg.JWT), authHandler.GetProfile)
		}

		// Media routes
		media := api.Group("/media")
		{
			media.POST("", middleware.Auth(cfg.JWT), mediaHandler.Create)
			media.PUT("/:id", middleware.Auth(cfg.JWT), mediaHandler.Update)
			media.GET("/search", mediaHandler.Search)
		}

		// Entry routes
		entries := api.Group("/entries")
		{
			entries.GET("", middleware.Auth(cfg.JWT), entryHandler.List)
			entries.POST("", middleware.Auth(cfg.JWT), entryHandler.Create)
			entries.GET("/:id", middleware.Auth(cfg.JWT), entryHandler.Get)
			entries.PATCH("/:id", middleware.Auth(cfg.JWT), entryHandler.Update)
			entries.DELETE("/:id", middleware.Auth(cfg.JWT), entryHandler.Delete)
			entries.POST("/sync", middleware.Auth(cfg.JWT), entryHandler.Sync)
		}

		// Collection routes
		collections := api.Group("/collections")
		{
			collections.GET("", middleware.Auth(cfg.JWT), collectionHandler.List)
			collections.POST("", middleware.Auth(cfg.JWT), collectionHandler.Create)
			collections.GET("/:id", middleware.Auth(cfg.JWT), collectionHandler.Get)
			collections.PATCH("/:id", middleware.Auth(cfg.JWT), collectionHandler.Update)
			collections.DELETE("/:id", middleware.Auth(cfg.JWT), collectionHandler.Delete)
			collections.POST("/:id/share", middleware.Auth(cfg.JWT), collectionHandler.CreateShare)
		}

		// Guest routes
		guest := api.Group("/guest")
		{
			guest.POST("/snapshot", guestHandler.CreateSnapshot)
			guest.POST("/merge", middleware.Auth(cfg.JWT), guestHandler.MergeToAccount)
		}

		// Public share routes
		api.GET("/s/:token", shareHandler.GetPublicShare)
	}

	// Public share routes (also available without /api prefix for direct access)
	router.GET("/s/:token", shareHandler.GetPublicShare)

	// Setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info().Str("port", cfg.Server.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info().Msg("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited")
}
