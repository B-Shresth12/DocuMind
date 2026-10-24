package router

import (
	"context"
	"net/http"
	"time"

	"github.com/B-Shresth12/DocuMind/backend/internal/auth"
	"github.com/B-Shresth12/DocuMind/backend/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RouterProvider struct {
	DB     *gorm.DB
	AppEnv string
}

func (p *RouterProvider) Register() (*gin.Engine, error) {
	sqlDB, err := p.DB.DB()

	if err != nil {
		return nil, err
	}

	engine := gin.New()
	engine.Use(gin.Logger(), gin.Recovery())

	engine.GET("/health/live", func(c *gin.Context) {
		response.Data(
			c,                        // gin context
			http.StatusOK,            //status
			"fetched",                //message
			gin.H{"status": "alive"}, //data
		)
	})

	engine.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(
			c.Request.Context(),
			2*time.Second,
		)
		defer cancel()

		if err := sqlDB.PingContext(ctx); err != nil {
			response.Data(
				c,
				http.StatusServiceUnavailable,
				"NOT READY",
				gin.H{
					"status": "not ready",
				},
			)
		}

		response.Data(
			c,
			http.StatusOK,
			"READY",
			gin.H{
				"status": "ready",
			},
		)
	})

	api := engine.Group("api/")
	auth.RegisterRoutes(api, p.DB, p.AppEnv)

	return engine, nil
}
