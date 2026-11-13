package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/B-Shresth12/DocuMind/backend/internal/app"
	"github.com/B-Shresth12/DocuMind/backend/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("fatal: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	application, err := app.New(ctx, &cfg)
	if err != nil {
		return err
	}
	defer func() {
		if err := application.Close(); err != nil {
			log.Printf("database connection cleanup failed: %v", err)
		}
	}()

	return application.Run(ctx)
}
