package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	server "github.com/rustingoff/excel_vue_go"
	"github.com/rustingoff/excel_vue_go/internal/controllers"
	"github.com/rustingoff/excel_vue_go/internal/middlewares"
	"github.com/rustingoff/excel_vue_go/internal/repositories"
	"github.com/rustingoff/excel_vue_go/internal/services"
	"github.com/rustingoff/excel_vue_go/packages/elastic_pkg"
	"github.com/rustingoff/excel_vue_go/packages/redis_pkg"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	elasticClient = elastic_pkg.NewElasticSearchConnection()
	redisClient   = redis_pkg.NewRedisConnection()

	userRepository = repositories.GetUserRepository(elasticClient, redisClient)
	userService    = services.GetUserService(userRepository)
	userController = controllers.GetUserController(userService)

	campaignRepository = repositories.GetCampaignRepository(elasticClient)
	campaignService    = services.GetCampaignService(campaignRepository)
	campaignController = controllers.GetCampaignController(campaignService)
)

func main() {

	router := gin.Default()
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "DELETE", "GET", "PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authController := router.Group("/auth")
	{
		authController.POST("/sign-up", userController.CreateUser) //only development
		authController.POST("/sign-in", userController.Login)
	}

	campaignRouter := router.Group("/campaign").Use(middlewares.CheckToken(userRepository))
	{
		campaignRouter.GET("/", campaignController.GetAllCampaigns)
		campaignRouter.GET("/:id", campaignController.GetCampaignById)
		campaignRouter.GET("/export/:id", campaignController.ExportCampaignExcel)
		campaignRouter.POST("/", campaignController.CreateCampaign)
		campaignRouter.DELETE("/:id", campaignController.DeleteCampaign)
	}

	srv := new(server.Server)

	go func() {
		// service connections
		if err := srv.Run(":8080", router); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //nolint:govet
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.ShutDown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	if <-ctx.Done(); true {
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
