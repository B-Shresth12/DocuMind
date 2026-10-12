package workspaces

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceMember struct {
	WorkspaceID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey;index"`
	Role        string    `gorm:"type:varchar(20);not null;check:role IN ('owner','admin','editor','viewer')"`
	JoinedAt    time.Time `gorm:"not null;autoCreateTime"`
}
