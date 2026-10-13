package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const SessionCookieName = "documind_session"

type SessionHandler struct {
	Service *SessionService
	Env     string
	DB      *gorm.DB
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email,max=320"`
	Password string `json:"password" binding:"required,email,min=1,max=128"`
}

func (h *SessionHandler) setSessionCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,              //name
		token,                          //value
		int(SessionDuration.Seconds()), //maxAge
		"/",                            //path
		"",                             //domain
		h.Env == "production",          //secure
		true,                           //httpOnly
	)
}

func (h *SessionHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid login details",
		})
		return
	}

	user, token, err := h.Service.Login(
		strings.ToLower(strings.TrimSpace(req.Email)),
		req.Password,
		c.Request.UserAgent(),
	)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid email or password",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not log in",
		})
		return
	}

	h.setSessionCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"email":        user.Email,
		"display_name": user.DisplayName,
	})
}

func (h *SessionHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})

		return
	}

	var user struct {
		ID          string `json:"id"`
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}

	if err := h.DB.
		Table("users").
		Select("id, email, display_name").
		Where("id = ? AND disabled_at IS NULL", userID).
		Take(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
	}

	c.JSON(http.StatusOK, user)
}

func (h *SessionHandler) Logout(c *gin.Context) {
	token, err := c.Cookie(SessionCookieName)
	if err == nil && token != "" {
		if err := h.Service.Logout(token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not logout",
			})
			return
		}
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		SessionCookieName,     //name
		"",                    //value
		-1,                    //maxAge
		"/",                   //path
		"",                    //domain
		h.Env == "production", //secure
		true,                  //httpOnly
	)

	c.Status(http.StatusNoContent)
}

var _ = time.Second
