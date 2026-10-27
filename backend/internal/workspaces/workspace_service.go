package workspaces

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkspaceService struct {
	DB *gorm.DB
}

func (s *WorkspaceService) AuthUserWorkspace(userID string) (*AuthWorkspaceList, error) {
	var workspaces []AuthWorkspace

	//fetching results and checking any issues
	if err := s.DB.
		Table("workspace_members AS wm").
		Select("w.id, w.name, wm.role").
		Joins("JOIN workspaces AS w ON w.id = wm.workspace_id").
		Where("wm.user_id = ?", userID).
		Order("w.created_at ASC").
		Scan(&workspaces).Error; err != nil { //if issue found

		//return nil value and logging error
		return nil, fmt.Errorf("get user workspace: %w", err)
	}

	//If workspace is nil or no workspace assigned to the user then send empty list
	if workspaces == nil {
		workspaces = []AuthWorkspace{}
	}

	//return value if every thing is green
	return &AuthWorkspaceList{
		Workspaces: workspaces,
	}, nil
}

func (s *WorkspaceService) AddAuthWorkspace(
	request WorkspaceRequest,
	userID string,
) (*AuthWorkspace, error) {
	var count int64

	createdBy, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	err = s.DB.
		Table("workspaces AS w").
		Joins("JOIN workspace_members AS wm ON wm.workspace_id = w.id").
		Where("wm.user_id = ? AND LOWER(TRIM(w.name)) = LOWER(TRIM(?))",
			createdBy, request.Name).
		Count(&count).Error

	if err != nil {
		return nil, fmt.Errorf("check workspace name: %w", err)
	}

	if count > 0 {
		return nil, fmt.Errorf("workspace name already exists")
	}

	var result AuthWorkspace

	err = s.DB.Transaction(func(tx *gorm.DB) error {
		workspace := &Workspace{
			Name:      request.Name,
			CreatedBy: createdBy,
		}

		if err := tx.Create(workspace).Error; err != nil {
			return err
		}

		member := WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      createdBy,
			Role:        "owner",
		}

		if err := tx.Create(&member).Error; err != nil {
			return err
		}

		result = AuthWorkspace{
			ID:   workspace.ID.String(),
			Name: workspace.Name,
			Role: member.Role,
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("create workspace: %w", err)
	}

	return &result, nil
}

// func (s *WorkspaceService) GetDetailsByID(id string, userID string) {

// }
