package workspaces

import (
	"fmt"

	"gorm.io/gorm"
)

type WorkspaceService struct {
	DB *gorm.DB
}

type AuthWorkspace struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}
type AuthWorkspaceList struct {
	Workspaces []AuthWorkspace `json:"workspaces"`
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
