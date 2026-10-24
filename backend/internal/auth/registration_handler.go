package auth

import (
	"errors"
	"net/http"

	"github.com/B-Shresth12/DocuMind/backend/httpx"
	"github.com/B-Shresth12/DocuMind/backend/response"
	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	Service *RegistrationService
}

type RegisterDTO struct {
	User RegisterResponse `json:"user"`
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Data(
			c,
			http.StatusUnprocessableEntity,
			"Validation Error",
			gin.H{
				"error": httpx.FormatValidationErrors(err),
			},
		)
		return
	}

	responseData, err := h.Service.Register(req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {

			response.Message(
				c,
				http.StatusConflict,
				"Email already registered",
			)

			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create account",
		})
		response.Message(
			c,
			http.StatusInternalServerError,
			"Something went wrong while registering",
		)
		return
	}

	response.Data(
		c,
		http.StatusCreated,
		"Registration successfull",
		RegisterDTO{
			User: RegisterResponse{
				ID:            responseData.ID,
				DisplayName:   responseData.DisplayName,
				Email:         responseData.Email,
				WorkspaceID:   responseData.WorkspaceID,
				WorkspaceName: responseData.WorkspaceName,
			},
		},
	)
}
