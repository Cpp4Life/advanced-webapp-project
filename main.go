package main

import (
	"advanced-webapp-project/config"
	"advanced-webapp-project/controller"
	"advanced-webapp-project/db"
	_ "advanced-webapp-project/docs"
	"advanced-webapp-project/middleware"
	"advanced-webapp-project/repository"
	"advanced-webapp-project/service"
	"advanced-webapp-project/utils"
	"advanced-webapp-project/websocket"
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	api = "/api/v1"
)

var (
	appConfig, _ = config.LoadConfig(".")
	sqlDB        = db.NewSQLDB()
	logger       = utils.NewLogger()

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
	oauthController = controller.NewOauthController(logger, jwtService, authService)
	userController  = controller.NewUserController(logger, jwtService, userService, groupService)
	groupController = controller.NewGroupController(logger, jwtService, groupService, userService, authService, mailService)
	presController  = controller.NewPresController(logger, jwtService, presService, userService)
	slideController = controller.NewSlideController(logger, slideService)
)

// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {
	defer db.Close(sqlDB)
	router := gin.Default()
	router.Use(cors.New(middleware.InitCors()))

	router.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Version": "0.0.1"})
	})

	authRoutes := router.Group(fmt.Sprintf("%s/auth", api))
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.GET("/verify-email/:code", authController.VerifyEmail)
		authRoutes.PUT("/forgot-password", authController.ForgotPassword)
	}

	oauthRoutes := router.Group(fmt.Sprintf("%s/oauth", api))
	{
		oauthRoutes.GET("/google/login", oauthController.GoogleOauthLogin)
		oauthRoutes.GET("/google/callback", oauthController.GoogleOauthCallback)
	}

	userRoutes := router.Group(fmt.Sprintf("%s/accounts", api)).Use(middleware.AuthorizeJWT(jwtService, logger))
	{
		userRoutes.GET("/profile", userController.GetProfile)
		userRoutes.PUT("/change-password", userController.ChangePassword)
		userRoutes.PUT("/edit", userController.UpdateProfile)
		userRoutes.POST("/create-group", groupController.CreateGroup)
		userRoutes.GET("/manage-groups", groupController.GetCreatedGroupsByUserId)
		userRoutes.GET("/joined-groups", groupController.GetJoinedGroupsByUserId)
	}

	groupRoutes := router.Group(fmt.Sprintf("%s/group", api))
	{
		groupRoutes.GET("/get-all", groupController.GetAllGroups)
		groupRoutes.GET("/:id/general", middleware.AuthorizeJWT(jwtService, logger), groupController.GetGroupById)
		groupRoutes.GET("/:id/details", middleware.AuthorizeJWT(jwtService, logger), groupController.GetGroupMemberDetailsByGroupId)
		groupRoutes.POST("/:id/edit", middleware.AuthorizeJWT(jwtService, logger), groupController.UpdateUserRole)
		groupRoutes.DELETE("/delete/:id", middleware.AuthorizeJWT(jwtService, logger), groupController.DeleteGroup)
		groupRoutes.POST("/:id/add-member", middleware.AuthorizeJWT(jwtService, logger), groupController.AddMemberToGroup)
		groupRoutes.DELETE("/:id/delete-member", middleware.AuthorizeJWT(jwtService, logger), groupController.DeleteMember)
		groupRoutes.POST("/:id/invite-member", middleware.AuthorizeJWT(jwtService, logger), groupController.InviteMember)
	}

	presRoutes := router.Group(fmt.Sprintf("%s/presentation", api)).Use(middleware.AuthorizeJWT(jwtService, logger))
	{
		presRoutes.GET("/:id/general", presController.GetPresentationById)
		presRoutes.GET("/get-all", presController.GetAllPresentationsByUserId)
		presRoutes.POST("/create", presController.CreatePresentation)
		presRoutes.PUT("/:id/edit", presController.UpdatePresentation)
		presRoutes.DELETE("/delete/:id", presController.DeletePresentation)
		presRoutes.POST("/:id/group-presentation", presController.PresentGroup)
		presRoutes.GET("/:id/slides/get-all", slideController.GetAllSlides)
		presRoutes.POST("/:id/slide/create", slideController.CreateSlide)
		presRoutes.PUT("/:id/slide/:slide_id/edit", slideController.UpdateSlide)
		presRoutes.DELETE("/:id/slide/delete/:slide_id", slideController.DeleteSlide)
		presRoutes.POST("/:id/vote/:content_id/submit", slideController.UpdateOptionVote)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.StaticFS("/public", http.Dir("templates"))

	// start websocket server
	hub := websocket.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		roomId := c.Query("roomId")
		logger.Info("room id: ", roomId)
		websocket.ServeWs(slideService, roomId, hub, c.Writer, c.Request)
	})

	logger.Info("Listening and serving HTTP on :", appConfig.Port)
	// start http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", appConfig.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Warn("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	logger.Warn("Server exiting")
}
