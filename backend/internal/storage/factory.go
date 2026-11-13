package storage

import (
	"fmt"

	"github.com/B-Shresth12/DocuMind/backend/internal/config"
)

func NewStorage(cfg config.Config) (Storage, error) {
	switch cfg.StorageDriver {
	case "local":
		return NewLocalStorage(cfg.StorageLocalPath)

	// case "s3":
	// 	return NewS3Storage(cfg)

	// case "minio":
	// 	return NewMinioStorage(cfg)

	default:
		return nil, fmt.Errorf(
			"unsupported storage driver: %s",
			cfg.StorageDriver,
		)
	}
}
