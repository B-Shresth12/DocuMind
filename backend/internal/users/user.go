package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"type:varchar(320);uniqueIndex:idx_users_email;not null"`
	PasswordHash string    `json:"-" gorm:"type:text;not null"`
	DisplayName  string    `gorm:"type:varchar(120);not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
	DisabledAt   *time.Time
}
