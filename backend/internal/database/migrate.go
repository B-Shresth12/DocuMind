package database

import (
	"fmt"

	"github.com/B-Shresth12/DocuMind/backend/constants"
	"github.com/B-Shresth12/DocuMind/backend/internal/auth"
	"github.com/B-Shresth12/DocuMind/backend/internal/users"
	"github.com/B-Shresth12/DocuMind/backend/internal/workspaces"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	fmt.Printf("%sMigrating Tables...%s\n", constants.ColorYellow, constants.ColorReset)

	if err := db.AutoMigrate(
		&users.User{},
		&workspaces.Workspace{},
		&workspaces.WorkspaceMember{},
		&auth.Session{},
	); err != nil {
		fmt.Printf("%sMigration failed%s\n", constants.ColorRed, constants.ColorReset)
		return fmt.Errorf("run database migrations: %w", err)
	}

	fmt.Printf("%sMigrating Tables Completed%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}
