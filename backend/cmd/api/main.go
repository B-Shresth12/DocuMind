package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/B-Shresth12/DocuMind/backend/internal/auth"
	"github.com/B-Shresth12/DocuMind/backend/internal/config"
	"github.com/B-Shresth12/DocuMind/backend/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	startupCtx, cancelStartup := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)

	db, err := database.NewPostgres(startupCtx, cfg.DatabaseURL)
	cancelStartup()

	if err != nil {
		log.Fatalf("database initialization failed: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get underlying database connection pool: %v", err)
	}

	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("database connection cleanup failed: %v", err)
		}
	}()

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "alive",
		})
	})

	router.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "not ready",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "ready",
		})
	})

	registrationService := &auth.RegistrationService{
		DB: db,
	}

	registrationHandler := &auth.RegistrationHandler{
		Service: registrationService,
	}

	sessionService := &auth.SessionService{
		DB: db,
	}

	sessionHandler := &auth.SessionHandler{
		Service: sessionService,
		Env:     cfg.AppEnv,
		DB:      db,
	}

	api := router.Group("/api/v1")

	api.POST("/auth/register", registrationHandler.Register)
	api.POST("/auth/login", sessionHandler.Login)

	protected := api.Group("")
	protected.Use(sessionHandler.RequireAuth())
	protected.GET("/me", sessionHandler.Me)
	protected.POST("/auth/logout", sessionHandler.Logout)

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	log.Printf("DocuMind API starting on port %s", cfg.AppPort)

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server failed: %v", err)
		}

	case <-signalCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
		}
	}
}
