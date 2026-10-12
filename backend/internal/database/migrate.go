package database

import (
	"fmt"

	"github.com/B-Shresth12/DocuMind/backend/internal/auth"
	"github.com/B-Shresth12/DocuMind/backend/internal/users"
	"github.com/B-Shresth12/DocuMind/backend/internal/workspaces"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&users.User{},
		&workspaces.Workspace{},
		&workspaces.WorkspaceMember{},
		&auth.Session{},
	); err != nil {
		return fmt.Errorf("run database migrations: %w", err)
	}

	return nil
}
