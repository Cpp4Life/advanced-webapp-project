package main

import (
	"advanced-webapp-project/controller"
	"advanced-webapp-project/db"
	"advanced-webapp-project/helper"
	"advanced-webapp-project/middleware"
	"advanced-webapp-project/repository"
	"advanced-webapp-project/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const PORT = ":7777"

var (
	sqlDB  = db.NewSQLDB()
	logger = helper.NewLogger()

	userRepo  = repository.NewUserRepo(sqlDB)
	groupRepo = repository.NewGroupRepo(sqlDB)

	jwtService   = service.NewJWTService(logger)
	authService  = service.NewAuthService(userRepo)
	userService  = service.NewUserService(userRepo)
	groupService = service.NewGroupService(groupRepo)

	authController  = controller.NewAuthHandler(logger, jwtService, authService)
	userController  = controller.NewUserController(logger, jwtService, userService, groupService)
	groupController = controller.NewGroupController(logger, jwtService, groupService, userService)
)

func main() {
	defer db.Close(sqlDB)
	r := gin.Default()
	r.Use(cors.New(middleware.InitCors()))

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

	userRoutes := r.Group("accounts").Use(middleware.AuthorizeJWT(jwtService, logger))
	{
		userRoutes.GET("/profile", userController.GetProfile)
		userRoutes.POST("/edit", userController.UpdateProfile)
		userRoutes.POST("/create-group", groupController.CreateGroup)
		userRoutes.GET("/manage-groups", groupController.GetCreatedGroupsByUserId)
		userRoutes.GET("/joined-groups", groupController.GetJoinedGroupsByUserId)
	}

	groupRoutes := r.Group("/group")
	{
		groupRoutes.GET("/get-all", groupController.GetAllGroups)
		groupRoutes.GET("/:id/details", groupController.GetGroupMemberDetailsByGroupId)
	}

	_ = r.Run(PORT)
}
