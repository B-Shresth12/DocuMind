package documents

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	db *gorm.DB,
	requireAuth gin.HandlerFunc,
) {
	documentService := &DocumentService{
		DB: db,
	}

	documentHandler := &DocumentHandler{
		Service: documentService,
	}

	v1 := api.Group("/v1")
	v1.Use(requireAuth)

	workspacesRoutes := v1.Group("/workspaces")
	workspacesRoutes.GET("/:id/documents", documentHandler.ListDocuments)

}
