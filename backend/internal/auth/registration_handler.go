package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegistrationHandler struct {
	Service *RegistrationService
}

func (h *RegistrationHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid registration details",
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
