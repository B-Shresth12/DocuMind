package workspaces

import (
	"time"

	"github.com/google/uuid"
)

type Workspace struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"type:varchar(120);not null"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
