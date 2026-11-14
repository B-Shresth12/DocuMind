package documents

import (
	"errors"
	"net/http"

	"github.com/B-Shresth12/DocuMind/backend/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	Service *DocumentService
}

func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	workspaceID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid workspace ID",
		})
		return
	}

	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthenticated",
		})
		return
	}

	documents, err := h.Service.WorkspaceDocuments(c, workspaceID, userID)
	if err != nil {
		if errors.Is(err, ErrWorkspaceAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "You do not have access to this workspace",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve documents",
		})
		return
	}

	response.Data(
		c,
		http.StatusOK,
		"Document Fetched",
		gin.H{
			"documents": documents,
		},
	)
}
