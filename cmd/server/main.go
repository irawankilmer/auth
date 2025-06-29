package main

import (
	"auth_service/internal/config"
	"auth_service/module"
	"auth_service/module/routes"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {
	config.LoadENV()
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	db := config.ConnectDB()
	app := module.AuthModule(db)

	r := gin.Default()
	api := r.Group("/api")
	routes.AuthRoutesRegister(api.Group("/auth"), app.AuthService)
	routes.ProfileRoutesRegister(api.Group("/profile"), app.ProfileService)

	_ = r.Run(":" + os.Getenv("APP_PORT"))
}
