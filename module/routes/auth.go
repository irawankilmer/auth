package routes

import (
	"auth_service/internal/handler"
	"auth_service/internal/middleware"
	"auth_service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gogaruda/valigo"
)

func AuthRoutesRegister(rg *gin.RouterGroup, authService service.AuthService) {
	validate := valigo.NewValigo()

	authHandler := handler.NewAuthHandler(authService, validate)

	rg.POST("/register", authHandler.Register)
	rg.POST("/login", authHandler.Login)

	auth := rg.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/logout", authHandler.Logout)

}
