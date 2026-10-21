package auth

import (
	"errors"
	"net/http"

	"github.com/B-Shresth12/DocuMind/backend/httpx"
	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	Service *RegistrationService
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid registration details",
			"error":   httpx.FormatValidationErrors(err),
		})
		return
	}

	response, err := h.Service.Register(req)
	if err != nil {
		if errors.Is(err, ErrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "email already registered",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not create account",
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}
