package workspaces

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	db *gorm.DB,
	requireAuth gin.HandlerFunc,
) {
	workspaceService := &WorkspaceService{
		DB: db,
	}

	workspaceHandler := &WorkspaceHandler{
		Service: workspaceService,
	}

	v1 := api.Group("/v1")
	workspaceRoutes := v1.Group("/workspaces")
	workspaceRoutes.Use(requireAuth)

	workspaceRoutes.GET("/", workspaceHandler.GetAuthWorkspaces)
	workspaceRoutes.POST("/", workspaceHandler.AddAuthWorkspace)
	workspaceRoutes.GET("/show/:id", workspaceHandler.GetDetails)
}
