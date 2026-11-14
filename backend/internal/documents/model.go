package documents

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Document struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	WorkspaceID uuid.UUID      `gorm:"type:uuid;not null;index" json:"workspace_id"`
	UploadedBy  uuid.UUID      `gorm:"type:uuid;not null;index" json:"uploaded_by"`
	Name        string         `gorm:"type:varchar(255);not null;index" json:"name"`
	StoragePath string         `gorm:"type:varchar(255);not null;index" json:"storage_path"`
	MimeType    string         `gorm:"type:varchar(100);not null;index" json:"mime_type"`
	Size        int64          `gorm:"not null" json:"size"`
	Checksum    string         `gorm:"type:varchar(64); not null" json:"checksum"`
	Status      string         `gorm:"type:varchar(30);not null; default:'uploaded'; index" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (d *Document) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}

	return nil
}
