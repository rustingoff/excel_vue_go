package services

import (
	"github.com/rustingoff/excel_vue_go/internal/models"
	"github.com/rustingoff/excel_vue_go/internal/repositories"
	"time"
)

type CampaignService interface {
	CreateCampaign(campaign models.Campaign) error
	DeleteCamapign(string) error

	GetAllCampaigns() ([]models.Campaign, error)
	GetCampaignById(string) (models.Campaign, error)
}

type campaignService struct {
	repo repositories.CampaignRepository
}

func GetCampaignService(repo repositories.CampaignRepository) CampaignService {
	return &campaignService{repo: repo}
}

func (service *campaignService) CreateCampaign(campaign models.Campaign) error {
	today, _ := time.LoadLocation("America/Los_Angeles")
	campaign.CampaignStartDate = time.Now().In(today).Format("01/02/2006")

	return service.repo.CreateCampaign(campaign)
}

func (service *campaignService) DeleteCamapign(id string) error {
	return service.repo.DeleteCamapign(id)
}

func (service *campaignService) GetAllCampaigns() ([]models.Campaign, error) {
	return service.repo.GetAllCampaigns()
}

func (service *campaignService) GetCampaignById(id string) (models.Campaign, error) {
	return service.repo.GetCampaignById(id)
}
