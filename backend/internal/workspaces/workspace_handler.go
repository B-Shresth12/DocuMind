package workspaces

import "github.com/gin-gonic/gin"

type WorkspaceHandler struct {
	Service *WorkspaceService
}

func (h *WorkspaceHandler) GetAuthWorkspaces(c *gin.Context) {

}
