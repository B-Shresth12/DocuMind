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

		workspaceReq := workspaces.WorkspaceRequest{
			Name: workspaceName,
		}

		workspaceService := workspaces.WorkspaceService{
			DB: tx,
		}

		workspace, err := workspaceService.AddAuthWorkspace(workspaceReq, user.ID.String())
		if err != nil {
			return fmt.Errorf("%v", err.Error())
		}

		response = RegisterResponse{
			ID:            user.ID.String(),
			Email:         user.Email,
			DisplayName:   user.DisplayName,
			WorkspaceID:   workspace.ID,
			WorkspaceName: workspace.Name,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &response, nil
}
