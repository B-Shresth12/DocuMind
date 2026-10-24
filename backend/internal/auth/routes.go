package auth

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(api *gin.RouterGroup, db *gorm.DB, appEnv string) {
	registrationService := &RegistrationService{
		DB: db,
	}

	registrationHandler := &RegistrationHandler{
		Service: registrationService,
	}

	sessionService := &SessionService{
		DB: db,
	}

	sessionHandler := &SessionHandler{
		Service: sessionService,
		Env:     appEnv,
		DB:      db,
	}

	v1 := api.Group("/v1")
	authRoutes := v1.Group("/auth")
	authRoutes.POST("/register", registrationHandler.Register)
	authRoutes.POST("/login", sessionHandler.Login)

	protected := v1.Group("")
	protected.Use(sessionHandler.RequireAuth())

	protected.GET("/me", sessionHandler.Me)
	protected.POST("/auth/logout", sessionHandler.Logout)
}
