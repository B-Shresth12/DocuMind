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

	"github.com/B-Shresth12/DocuMind/backend/internal/config"
	"github.com/B-Shresth12/DocuMind/backend/internal/database"
	"github.com/B-Shresth12/DocuMind/backend/internal/router"
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

	routerProvider := &router.RouterProvider{
		DB:     db,
		AppEnv: cfg.AppEnv,
	}

	httpRouter, err := routerProvider.Register()
	if err != nil {
		log.Fatalf("router initialization failedL: %v", err)
	}

	server := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           httpRouter,
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
