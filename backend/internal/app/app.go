package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/B-Shresth12/DocuMind/backend/internal/config"
	"github.com/B-Shresth12/DocuMind/backend/internal/database"
	"github.com/B-Shresth12/DocuMind/backend/internal/router"
	"github.com/gin-gonic/gin"
)

const (
	startupTimeout  = 10 * time.Second
	shutdownTimeout = 10 * time.Second
)

type App struct {
	cfg    *config.Config
	sqlDB  *sql.DB
	server *http.Server
}

// New wires all dependencies. On error it cleans up anything it already opened.
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	startupCtx, cancel := context.WithTimeout(ctx, startupTimeout)
	defer cancel()

	db, err := database.NewPostgres(startupCtx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("database initialization: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get database connection pool: %w", err)
	}

	if err := database.Migrate(db); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("database migration: %w", err)
	}

	provider := &router.RouterProvider{DB: db, AppEnv: cfg.AppEnv}

	handler, err := provider.Register()
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("router initialization: %w", err)
	}

	return &App{
		cfg:   cfg,
		sqlDB: sqlDB,
		server: &http.Server{
			Addr:              ":" + cfg.AppPort,
			Handler:           handler,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}, nil
}

// Run blocks until the server fails or ctx is cancelled, then shuts down gracefully.
func (a *App) Run(ctx context.Context) error {
	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- a.server.ListenAndServe()
	}()

	log.Printf("DocuMind API starting on port %s", a.cfg.AppPort)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("http server: %w", err)
		}
		return nil

	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("graceful shutdown: %w", err)
		}
		return nil
	}
}

func (a *App) Close() error {
	return a.sqlDB.Close()
}
