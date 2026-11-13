package storage

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, path string, file io.Reader) error

	Get(ctx context.Context, path string) (io.ReadCloser, error)

	Delete(ctx context.Context, path string) error
}
