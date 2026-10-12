package auth

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null;index"`
	TokenHash  []byte    `gorm:"type:bytea;uniqueIndex;not null"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime"`
	LastSeenAt time.Time `gorm:"not null;autoCreateTime"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	RevokedAt  *time.Time
	UserAgent  string `gorm:"type:text"`
	IPAddress  string `gorm:"type:inet"`
}
