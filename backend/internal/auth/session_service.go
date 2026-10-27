package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/B-Shresth12/DocuMind/backend/internal/users"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const SessionCookieName = "documind_session"
const SessionDuration = 7 * 24 * time.Hour

var ErrInvalidCredentials = errors.New("invalid email or password")

type SessionService struct {
	DB *gorm.DB
}

func ipPtr(raw string) *string {
	if net.ParseIP(raw) == nil {
		return nil
	}
	return &raw
}

func (s *SessionService) GetSessionCookieName() string {
	return SessionCookieName
}

func (s *SessionService) Login(
	email, password, userAgent string, ip string,
) (*users.User, string, error) {
	var user users.User

	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", ErrInvalidCredentials
		}
		return nil, "", fmt.Errorf("find user: %w", err)
	}

	if user.DisabledAt != nil {
		return nil, "", ErrInvalidCredentials
	}

	valid, err := VerifyPassword(password, user.PasswordHash)
	if err != nil {
		return nil, "", fmt.Errorf("verify password hash: %w", err)
	}

	if !valid {
		return nil, "", ErrInvalidCredentials
	}

	//Generating a session token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, "", fmt.Errorf("genereate session token: %w", err)
	}

	token := base64.RawURLEncoding.EncodeToString(tokenBytes)
	tokenHash := sha256.Sum256([]byte(token))

	session := Session{
		UserID:     user.ID,
		TokenHash:  tokenHash[:],
		ExpiresAt:  time.Now().Add(SessionDuration),
		UserAgent:  userAgent,
		LastSeenAt: time.Now(),
		IPAddress:  ipPtr(ip),
	}

	if err := s.DB.Create(&session).Error; err != nil {
		return nil, "", fmt.Errorf("create session: %w", err)
	}
	return &user, token, nil
}

func (s *SessionService) FindSessionUser(token string) (*users.User, error) {
	tokenHash := sha256.Sum256([]byte(token))

	var session Session
	err := s.DB.Where(
		"token_hash = ? AND revoked_at IS NULL AND expires_at > ?",
		tokenHash[:],
		time.Now(),
	).First(&session).Error

	if err != nil {
		return nil, err
	}

	var user users.User

	err = s.DB.Where(
		"id = ? AND disabled_at IS NULL",
		session.UserID,
	).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (s *SessionService) Logout(token string) error {
	tokenHash := sha256.Sum256([]byte(token))
	now := time.Now()

	result := s.DB.Model(&Session{}).
		Where("token_hash = ? AND revoked_at IS NULL", tokenHash[:]).
		Update("revoked_at", now)

	if result.Error != nil {
		return fmt.Errorf("revoke session: %w", result.Error)
	}

	return nil
}

func (s *SessionService) FindSessionUserID(token string) (string, error) {
	user, err := s.FindSessionUser(token)
	if err != nil {
		return "", err
	}

	return fmt.Sprint(user.ID), nil
}

var _ uuid.UUID
