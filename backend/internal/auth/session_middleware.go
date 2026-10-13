package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *SessionHandler) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(SessionCookieName)
		if err != nil || token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authentication required",
			})
			return
		}

		user, err := h.Service.FindSessionUser(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired session",
			})

			return
		}

		c.Set("user_id", user.ID)
		c.Next()
	}
}
