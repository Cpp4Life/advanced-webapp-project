package main

import (
	"advanced-webapp-project/controller"
	"advanced-webapp-project/db"
	"advanced-webapp-project/helper"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
)

const PORT = ":7777"

var (
	sqlDB  = db.NewSQLDB()
	logger = helper.NewLogger()

	jwtService = service.NewJWTService(logger)

	authController = controller.NewAuthHandler(logger, jwtService)
)

func main() {
	defer sqlDB.Close()
	r := gin.Default()

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	_ = r.Run(PORT)
}
