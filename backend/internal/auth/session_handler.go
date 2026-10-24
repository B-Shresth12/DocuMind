package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/B-Shresth12/DocuMind/backend/httpx"
	"github.com/B-Shresth12/DocuMind/backend/response"
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
	Password string `json:"password" binding:"required,min=1,max=128"`
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

type UserDTO struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
}

type LoginData struct {
	User UserDTO `json:"user"`
}

func (h *SessionHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Data(
			c,
			http.StatusUnprocessableEntity,
			"invalid login details",
			gin.H{
				"errors": httpx.FormatValidationErrors(err),
			},
		)
		return
	}

	user, token, err := h.Service.Login(
		strings.ToLower(strings.TrimSpace(req.Email)),
		req.Password,
		c.Request.UserAgent(),
		c.ClientIP(),
	)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Message(
				c,
				http.StatusUnauthorized,
				"Invalid email or password",
			)
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not log in",
		})

		response.Message(
			c,
			http.StatusInternalServerError,
			"Somthing went wrong while logging in",
		)
		return
	}

	h.setSessionCookie(c, token)

	response.Data(c, http.StatusOK, "Logged in successfully", LoginData{
		User: UserDTO{
			ID:          user.ID.String(),
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	})
}

func (h *SessionHandler) Me(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {

		response.Message(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)

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
		response.Message(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)
	}

	response.Data(c, http.StatusOK, "Logged in successfully", LoginData{
		User: UserDTO{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	})
}

func (h *SessionHandler) Logout(c *gin.Context) {
	token, err := c.Cookie(SessionCookieName)
	if err == nil && token != "" {
		if err := h.Service.Logout(token); err != nil {

			response.Message(
				c,
				http.StatusInternalServerError,
				"something went wrong while logging out",
			)
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
