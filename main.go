package main

import (
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
	PORT = ":7777"
	api  = "/api/v1"
)

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

	authRoutes := router.Group(fmt.Sprintf("%s/auth", api))
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
		authRoutes.GET("/verify-email/:code", authController.VerifyEmail)
		authRoutes.GET("/forgot-password", authController.ForgotPassword)
	}

	userRoutes := router.Group(fmt.Sprintf("%s/accounts", api)).Use(middleware.AuthorizeJWT(jwtService, logger))
	{
		userRoutes.GET("/profile", userController.GetProfile)
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
		presRoutes.GET("/get-all", presController.GetAllPresentations)
		presRoutes.POST("/create", presController.CreatePresentation)
		presRoutes.PUT("/:id/edit", presController.UpdatePresentation)
		presRoutes.DELETE("/delete/:id", presController.DeletePresentation)
		presRoutes.GET("/:id/slides/get-all", slideController.GetAllSlides)
		presRoutes.POST("/:id/slide/create", slideController.CreateSlide)
		presRoutes.PUT("/:id/slide/:slide_id/edit", slideController.UpdateSlide)
		presRoutes.DELETE("/:id/slide/delete/:slide_id", slideController.DeleteSlide)
		presRoutes.POST("/:id/vote/:content_id/submit", slideController.UpdateOptionVote)
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// start websocket server
	hub := websocket.NewHub()
	go hub.Run()
	router.GET("/ws", func(c *gin.Context) {
		logger.Warn("hello websocket")
		websocket.ServeWs(hub, c.Writer, c.Request)
	})

	// start http server
	srv := &http.Server{
		Addr:    PORT,
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

//router.GET("/socket.io/*any", gin.WrapH(sServer))
//router.POST("/socket.io/*any", gin.WrapH(sServer))
//router.StaticFS("/public", http.Dir("templates"))

//sServer := socketio.NewServer(&engineio.Options{
//	Transports: []transport.Transport{
//		&polling.Transport{
//			Client: &http.Client{
//				Timeout: time.Minute,
//			},
//		},
//	},
//})

//sServer.OnConnect("/", func(s socketio.Conn) error {
//	s.SetContext("")
//	logger.Info(fmt.Sprintf("[CONNECTED]: %s", s.ID()))
//	s.Join("bcast")
//	return nil
//})
//
//sServer.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
//	logger.Info(fmt.Sprintf("[NOTICE]: %s", msg))
//	s.Emit("reply", "[MESSAGE]: "+msg)
//})
//
//sServer.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
//	msg = fmt.Sprintf("User #%s: %s", s.ID(), msg)
//	s.SetContext(msg)
//	sServer.BroadcastToRoom("", "bcast", "reply", msg)
//	return "receive " + msg
//})
//
//sServer.OnEvent("/", "bye", func(s socketio.Conn) string {
//	last := s.Context().(string)
//	s.Emit("bye", last)
//	s.Close()
//	return last
//})
//
//sServer.OnError("/", func(s socketio.Conn, e error) {
//	logger.Error(fmt.Sprintf("[ERROR]: %s", e.Error()))
//})
//
//sServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
//	logger.Warn(fmt.Sprintf("[CLOSED]: %s", reason))
//})

//go func() {
//	if err := sServer.Serve(); err != nil {
//		log.Fatalf("socketio listen error: %s\n", err)
//	}
//}()
//defer sServer.Close()
