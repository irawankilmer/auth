package routes

import (
	"auth_service/internal/handler"
	"auth_service/internal/middleware"
	"auth_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/valigo"
)

func ProfileRoutesRegister(rg *gin.RouterGroup, profileService service.ProfileService) {
	validate := valigo.NewValigo()

	profileHandler := handler.NewProfileHandler(profileService, validate)

	auth := rg.Group("/")
	auth.Use(middleware.AuthMiddleware())

	auth.POST("/setting", profileHandler.Setting)

	pCompleted := auth.Group("/")
	pCompleted.Use(middleware.RequireCompletedProfile(profileService))
	pCompleted.GET("/me", profileHandler.Me)
}
