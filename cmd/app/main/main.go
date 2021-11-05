package main

import (
	"github.com/gin-gonic/gin"
	server "github.com/rustingoff/excel_vue_go"
	"github.com/rustingoff/excel_vue_go/internal/controllers"
	"github.com/rustingoff/excel_vue_go/internal/repositories"
	"github.com/rustingoff/excel_vue_go/internal/services"
	"github.com/rustingoff/excel_vue_go/packages/elastic_pkg"
)

var (
	elasticClient = elastic_pkg.NewElasticSearchConnection()

	campaignRepository = repositories.GetCampaignRepository(elasticClient)
	campaignService    = services.GetCampaignService(campaignRepository)
	campaignController = controllers.GetCampaignController(campaignService)
)

func main() {

	router := gin.Default()

	campaignRouter := router.Group("/campaign")
	{
		campaignRouter.GET("/", campaignController.GetAllCampaigns)
		campaignRouter.GET("/:id", campaignController.GetCampaignById)
		campaignRouter.POST("/", campaignController.CreateCampaign)
		campaignRouter.DELETE("/:id", campaignController.DeleteCamapign)
	}

	srv := new(server.Server)
	if err := srv.Run(":8080", router); err != nil {
		panic(err)
	}
}
