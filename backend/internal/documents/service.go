package documents

import (
	"context"
	"errors"

	"github.com/B-Shresth12/DocuMind/backend/internal/storage"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrWorkspaceAccessDenied = errors.New("workspace access denied")

type DocumentService struct {
	DB      *gorm.DB
	Storage storage.Storage
}

func (s *DocumentService) WorkspaceDocuments(
	ctx context.Context,
	workspaceID uuid.UUID,
	userID uuid.UUID,
) ([]Document, error) {
	// Check whether the user belongs to this workspace.
	var membershipCount int64

	err := s.DB.WithContext(ctx).
		Table("workspace_members").
		Where("workspace_id = ? AND user_id = ?", workspaceID, userID).
		Count(&membershipCount).Error

	if err != nil {
		return nil, err
	}

	if membershipCount == 0 {
		return nil, ErrWorkspaceAccessDenied
	}

	// Retrieve documents belonging to the workspace.
	var documents []Document

	err = s.DB.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("created_at DESC").
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}
