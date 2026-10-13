package users

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `grom:"type:uuid;default:gen_random_uuid();primaryKey"`
	Email        string    `grom:"type:varchar(320);uniqueIndex;not null"`
	PasswordHash string    `grom:"type:text;not null"`
	DisplayName  string    `grom:"type:varchar(120);not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DisabledAt   *time.Time
}
