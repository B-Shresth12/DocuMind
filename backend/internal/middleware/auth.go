package middleware

import (
	"net/http"

	"github.com/B-Shresth12/DocuMind/backend/response"
	"github.com/gin-gonic/gin"
)

type SessionValidator interface {
	FindSessionUserID(token string) (string, error)
	GetSessionCookieName() string
}

type AuthMiddleware struct {
	Sessions SessionValidator
}

func NewAuthMiddleware(sessions SessionValidator) *AuthMiddleware {
	return &AuthMiddleware{
		Sessions: sessions,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(m.Sessions.GetSessionCookieName())
		if err != nil || token == "" {
			response.Message(
				c,
				http.StatusUnauthorized,
				"not authenticated",
			)

			c.Abort()
			return
		}

		userID, err := m.Sessions.FindSessionUserID(token)
		if err != nil {
			response.Message(
				c,
				http.StatusUnauthorized,
				"invalid or expired session",
			)
			c.Abort()

			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}
