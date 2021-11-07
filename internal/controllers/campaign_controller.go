package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rustingoff/excel_vue_go/internal/models"
	"github.com/rustingoff/excel_vue_go/internal/services"
	"log"
	"net/http"
)

type CampaignController interface {
	CreateCampaign(c *gin.Context)
	DeleteCampaign(c *gin.Context)

	GetAllCampaigns(c *gin.Context)
	GetCampaignById(c *gin.Context)

	ExportCampaignExcel(c *gin.Context)
}

type campaignController struct {
	campaignService services.CampaignService
}

func GetCampaignController(service services.CampaignService) CampaignController {
	return &campaignController{campaignService: service}
}

func (controller *campaignController) CreateCampaign(c *gin.Context) {
	userByToken, _ := c.Get("currentUser")

	var campaign models.Campaign
	campaign.UserID = userByToken.(models.User).ID

	if err := c.ShouldBindJSON(&campaign); err != nil {
		log.Println("[ERR]: failed binding json to structure, ", err.Error())
		c.AbortWithStatusJSON(http.StatusUnprocessableEntity, "invalid json structure")
		return
	}

	err := controller.campaignService.CreateCampaign(campaign)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to create")
		return
	}

	log.Println("[INF]: campaign was created successfully")
	c.JSON(http.StatusCreated, "OK")
}

func (controller *campaignController) DeleteCampaign(c *gin.Context) {
	var campaignID = c.Param("id")

	err := controller.campaignService.DeleteCampaign(campaignID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to delete")
		return
	}

	log.Println("[INF]: campaign was deleted successfully")
	c.JSON(http.StatusNoContent, "deleted")
}

func (controller *campaignController) GetAllCampaigns(c *gin.Context) {
	userByToken, _ := c.Get("currentUser")
	userID := userByToken.(models.User).ID

	campaigns, err := controller.campaignService.GetAllCampaigns(userID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get all")
		return
	}

	log.Println("[INF]: successfully got campaigns")
	c.JSON(http.StatusOK, campaigns)
}

func (controller *campaignController) GetCampaignById(c *gin.Context) {
	campaignID := c.Param("id")

	campaign, err := controller.campaignService.GetCampaignById(campaignID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to get one")
		return
	}

	log.Println("[INF]: successfully got campaign by id")
	c.JSON(http.StatusOK, campaign)
}

func (controller *campaignController) ExportCampaignExcel(c *gin.Context) {
	id := c.Param("id")

	err := controller.campaignService.ExportCampaignExcel(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, "failed to export")
		return
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.File("./static/exports/" + id + ".xlsx")

	log.Println("[INF]: excel was exported successfully")
}
