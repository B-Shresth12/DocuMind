package storage

import "github.com/B-Shresth12/DocuMind/backend/internal/config"

func NewMinioStorage(cfg config.Config) (*Storage, error) {
	return &Storage{}, nil
}
