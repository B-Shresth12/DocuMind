package documents

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxDocumentSize   = 20 << 20 //20MB
	MaxFilenameLength = 255
)

var AllowedExtensions = map[string]bool{
	".pdf":  true,
	".docx": true,
	".txt":  true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
}

type UploadedDocumentRequest struct {
	WorkspaceID uuid.UUID
	Filename    string
	Size        int64
}

func (r *UploadedDocumentRequest) Validate() error {
	if r.WorkspaceID == uuid.Nil {
		return errors.New("workspace ID is required")
	}

	r.Filename = strings.TrimSpace(r.Filename)

	if r.Filename == "" {
		return errors.New("Filename is required")
	}

	if len(r.Filename) > MaxFilenameLength {
		return errors.New(fmt.Sprintf("filename must not exceed %v characters", MaxFilenameLength))
	}

	if r.Size <= 0 {
		return errors.New("file cannot be empty")
	}

	ext := strings.ToLower(filepath.Ext(r.Filename))

	if !AllowedExtensions[ext] {
		return errors.New("unsupported file format")
	}

	return nil
}
