package workspaces

import (
	"net/http"

	"github.com/B-Shresth12/DocuMind/backend/httpx"
	"github.com/B-Shresth12/DocuMind/backend/response"
	"github.com/gin-gonic/gin"
)

type WorkspaceHandler struct {
	Service *WorkspaceService
}

func (h *WorkspaceHandler) GetAuthWorkspaces(c *gin.Context) {
	value, exists := c.Get("user_id")
	if !exists {
		response.Message(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	userId, ok := value.(string)
	if !ok {
		response.Message(
			c,
			http.StatusInternalServerError,
			"Invalid User ID",
		)

		return
	}

	workspaces, err := h.Service.AuthUserWorkspace(userId)
	if err != nil {
		response.Message(c,
			http.StatusInternalServerError,
			"Failed to fetch workspace",
		)
		return
	}

	response.Data(
		c,
		http.StatusOK,
		"Workspace fetched",
		workspaces,
	)

}

func (h *WorkspaceHandler) AddAuthWorkspace(c *gin.Context) {

	userID, exists := c.Get("user_id")
	if !exists {
		response.Message(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.Message(c, http.StatusUnauthorized, "Invalid user identity")
		return
	}

	var request WorkspaceRequest
	if err := c.ShouldBind(&request); err != nil {
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

	workspace, err := h.Service.AddAuthWorkspace(request, userIDStr)
	if err != nil {
		response.Message(
			c,
			http.StatusInternalServerError,
			err.Error(),
		)
		return
	}

	response.Data(
		c,
		200,
		"Workspace Created!!",
		workspace,
	)
}

func (h *WorkspaceHandler) getDetails(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		response.Message(
			c,
			http.StatusUnauthorized,
			"Unauthorized",
		)
		return
	}

	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		response.Message(c, http.StatusUnauthorized, "Invalid user identity")
		return
	}

	// workspace, err := h.Service.GetDetailsByID(id, userIDStr)
}
