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
	presRepo  = repository.NewPresRepo(sqlDB)
	slideRepo = repository.NewSlideRepo(sqlDB)

	jwtService   = service.NewJWTService(logger)
	mailService  = service.NewMailerService(logger)
	authService  = service.NewAuthService(userRepo)
	userService  = service.NewUserService(userRepo)
	groupService = service.NewGroupService(groupRepo)
	presService  = service.NewPresService(presRepo)
	slideService = service.NewSlideService(slideRepo)

	authController  = controller.NewAuthHandler(logger, jwtService, authService, mailService)
	userController  = controller.NewUserController(logger, jwtService, userService, groupService)
	groupController = controller.NewGroupController(logger, jwtService, groupService, userService)
	presController  = controller.NewPresController(logger, jwtService, presService, userService)
	slideController = controller.NewSlideController(logger, slideService)
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
		userRoutes.PUT("/edit", userController.UpdateProfile)
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

	presRoutes := r.Group("/presentation").Use(middleware.AuthorizeJWT(jwtService, logger))
	{
		presRoutes.GET("/:id/general", presController.GetPresentationById)
		presRoutes.GET("/get-all", presController.GetAllPresentations)
		presRoutes.POST("/create", presController.CreatePresentation)
		presRoutes.PUT("/:id/edit", presController.UpdatePresentation)
		presRoutes.DELETE("/delete/:id", presController.DeletePresentation)
		presRoutes.GET("/:id/slides/get-all", slideController.GetAllSlides)
		presRoutes.POST("/:id/slide/create", slideController.CreateSlide)
		presRoutes.PUT("/:id/slide/:slide_id/edit", slideController.UpdateSlide)
		presRoutes.DELETE("/:id/slide/delete/:slide_id", slideController.DeleteSlide)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	_ = r.Run(PORT)
}
