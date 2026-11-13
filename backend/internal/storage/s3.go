package storage

import "github.com/B-Shresth12/DocuMind/backend/internal/config"

func NewS3Storage(cfg config.Config) (*Storage, error) {
	return &Storage{}, nil
}
