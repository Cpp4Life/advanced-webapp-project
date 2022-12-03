package main

import (
	"advanced-webapp-project/controller"
	"advanced-webapp-project/db"
	_ "advanced-webapp-project/docs"
	"advanced-webapp-project/middleware"
	"advanced-webapp-project/repository"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const PORT = ":7777"

var (
	sqlDB  = db.NewSQLDB()
	logger = utils.NewLogger()

	userRepo  = repository.NewUserRepo(sqlDB)
	groupRepo = repository.NewGroupRepo(sqlDB)

	jwtService   = service.NewJWTService(logger)
	mailService  = service.NewMailerService(logger)
	authService  = service.NewAuthService(userRepo)
	userService  = service.NewUserService(userRepo)
	groupService = service.NewGroupService(groupRepo)

	authController  = controller.NewAuthHandler(logger, jwtService, authService, mailService)
	userController  = controller.NewUserController(logger, jwtService, userService, groupService)
	groupController = controller.NewGroupController(logger, jwtService, groupService, userService)
)

// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {
	defer db.Close(sqlDB)
	r := gin.Default()
	r.Use(cors.New(middleware.InitCors()))

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.GET("/verify-email/:code", authController.VerifyEmail)
	}

	userRoutes := r.Group("/accounts").Use(middleware.AuthorizeJWT(jwtService, logger))
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
		groupRoutes.GET("/:id/general", middleware.AuthorizeJWT(jwtService, logger), groupController.GetGroupById)
		groupRoutes.GET("/:id/details", middleware.AuthorizeJWT(jwtService, logger), groupController.GetGroupMemberDetailsByGroupId)
		groupRoutes.POST("/:id/edit", middleware.AuthorizeJWT(jwtService, logger), groupController.UpdateUserRole)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = r.Run(PORT)
}
