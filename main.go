package main

import (
	"advanced-webapp-project/controller"
	"advanced-webapp-project/db"
	"advanced-webapp-project/helper"
	"advanced-webapp-project/middleware"
	"advanced-webapp-project/repository"
	"advanced-webapp-project/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

const PORT = ":7777"

var (
	sqlDB  = db.NewSQLDB()
	logger = helper.NewLogger()

	userRepo = repository.NewUserRepo(sqlDB)

	jwtService  = service.NewJWTService(logger)
	authService = service.NewAuthService(userRepo)

	authController = controller.NewAuthHandler(logger, jwtService, authService)
)

func main() {
	defer db.Close(sqlDB)
	r := gin.Default()

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	groupRoutes := r.Group("/groups", middleware.AuthorizeJWT(jwtService, logger))
	{
		groupRoutes.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, map[string]any{"message": "success"})
		})
	}

	_ = r.Run(PORT)
}
