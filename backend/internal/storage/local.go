package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type LocalStorage struct {
	BasePath string
}

// Constructor
func NewLocalStorage(basePath string) (*LocalStorage, error) {
	if strings.TrimSpace(basePath) == "" {
		return nil, errors.New("Storage base path cannot be empty")
	}

	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, fmt.Errorf("resolve storage path: %v", err)
	}

	if err := os.MkdirAll(absPath, 0700); err != nil {
		return nil, fmt.Errorf("create storage directory: %v", err)
	}

	return &LocalStorage{
		BasePath: absPath,
	}, nil
}

// Resolve and validate a storage key
func (s *LocalStorage) resolvePath(key string) (string, error) {
	if strings.TrimSpace(key) == "" {
		return "", errors.New("storage key cannot be empty")
	}

	key = filepath.FromSlash(key)

	if filepath.IsAbs(key) || filepath.VolumeName(key) != "" {
		return "", errors.New("absolute paths are not allowed")
	}

	cleanPath := filepath.Clean(key)

	if cleanPath == "." || cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(os.PathSeparator)) {
		return "", errors.New("invalid storage path")
	}

	fullPath := filepath.Join(s.BasePath, cleanPath)

	rel, err := filepath.Rel(s.BasePath, fullPath)
	if err != nil {
		return "", fmt.Errorf("resolve relative path: %w", err)
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", errors.New("storage path escapes base directory")
	}

	return fullPath, nil
}

// Upload stores a file in the local filesystem
func (s *LocalStorage) Upload(
	ctx context.Context,
	key string,
	file io.Reader,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fullPath, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(fullPath), 0700); err != nil {
		return fmt.Errorf("create document directory: %w", err)
	}

	dst, err := os.OpenFile(
		fullPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	_, copyErr := io.Copy(dst, &contextReader{
		ctx: ctx,
		r:   file,
	})

	closeErr := dst.Close()

	if copyErr != nil {
		return fmt.Errorf("write file: %w", copyErr)
	}

	if closeErr != nil {
		return fmt.Errorf("close file: %w", closeErr)
	}

	return nil
}

// Get retrieves a file from local storage.
func (s *LocalStorage) Get(
	ctx context.Context,
	key string,
) (io.ReadCloser, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	fullPath, err := s.resolvePath(key)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	return file, nil
}

// Delete removes a file from local storage.
func (s *LocalStorage) Delete(
	ctx context.Context,
	key string,
) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fullPath, err := s.resolvePath(key)
	if err != nil {
		return err
	}

	err = os.Remove(fullPath)

	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	return nil
}

// contextReader checks cancellation during file uploads.
type contextReader struct {
	ctx context.Context
	r   io.Reader
}

func (r *contextReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}

	return r.r.Read(p)
}

// Compile-time interface verification.
var _ Storage = (*LocalStorage)(nil)
