package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/B-Shresth12/DocuMind/backend/internal/users"
	"github.com/B-Shresth12/DocuMind/backend/internal/workspaces"
	"gorm.io/gorm"
)

var ErrEmailAlreadyExists = errors.New("email already registered")

type RegistrationService struct {
	DB *gorm.DB
}

func (s *RegistrationService) Register(
	req RegisterRequest,
) (*RegisterResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	displayName := strings.TrimSpace(req.DisplayName)
	workspaceName := strings.TrimSpace(req.WorkspaceName)

	passwordHash, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	var response RegisterResponse

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		user := users.User{
			Email:        email,
			PasswordHash: passwordHash,
			DisplayName:  displayName,
		}

		if err := tx.Create(&user).Error; err != nil {
			// PostgreSQL unique-constraint violations use SQLSTATE 23505.
			// We will refine this mapping to inspect the exact constraint.
			if strings.Contains(strings.ToLower(err.Error()), "unique") ||
				strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
				return ErrEmailAlreadyExists
			}
			return fmt.Errorf("create user: %w", err)
		}

		workspace := workspaces.Workspace{
			Name:      workspaceName,
			CreatedBy: user.ID,
		}

		if err := tx.Create(&workspace).Error; err != nil {
			return fmt.Errorf("create workspace: %w", err)
		}

		member := workspaces.WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      user.ID,
			Role:        "owner",
		}

		if err := tx.Create(&member).Error; err != nil {
			return fmt.Errorf("create workspace membership: %w", err)
		}

		response = RegisterResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DisplayName:   user.DisplayName,
			WorkspaceID:   workspace.ID.String(),
			WorkspaceName: workspace.Name,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}
