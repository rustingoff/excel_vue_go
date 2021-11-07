package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/olivere/elastic/v7"
	"github.com/rustingoff/excel_vue_go/internal/models"
	"log"
)

type CampaignRepository interface {
	CreateCampaign(campaign models.Campaign) error
	DeleteCampaign(string) error

	GetAllCampaigns(userID string) ([]models.Campaign, error)
	GetCampaignById(string) (models.Campaign, error)
}

const _CampaignIndex = "campaigns"

type campaignRepository struct {
	elasticClient *elastic.Client
}

func GetCampaignRepository(e *elastic.Client) CampaignRepository {
	return &campaignRepository{e}
}

func (repo *campaignRepository) CreateCampaign(campaign models.Campaign) error {

	_, err := repo.elasticClient.Index().Index(_CampaignIndex).BodyJson(campaign).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to create campaign, ", err.Error())
		return err
	}

	return nil
}

func (repo *campaignRepository) DeleteCampaign(id string) error {

	_, err := repo.elasticClient.Delete().Index(_CampaignIndex).Id(id).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to delete campaign, ", err.Error())
		return err
	}

	return nil
}

func (repo *campaignRepository) GetAllCampaigns(userID string) ([]models.Campaign, error) {
	query := elastic.NewMatchQuery("user_id", userID)

	res, err := repo.elasticClient.Search().Index(_CampaignIndex).Query(query).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to get all campaigns")
		return nil, err
	}

	var campaigns = make([]models.Campaign, 0)

	for i := 0; i < int(res.TotalHits()); i++ {
		var campaign models.Campaign

		err = json.Unmarshal(res.Hits.Hits[i].Source, &campaign)
		if err != nil {
			log.Println("[ERR]: failed to unmarshal source")
			return nil, err
		}

		campaign.ID = res.Hits.Hits[i].Id
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

func (repo *campaignRepository) GetCampaignById(id string) (models.Campaign, error) {
	query := elastic.NewMatchQuery("_id", id)

	res, err := repo.elasticClient.Search(_CampaignIndex).Query(query).Do(context.TODO())
	if err != nil {
		log.Println("[ERR]: failed to get campaign by id")
		return models.Campaign{}, err
	}

	if int(res.TotalHits()) > 0 {
		var campaign models.Campaign

		err = json.Unmarshal(res.Hits.Hits[0].Source, &campaign)
		if err != nil {
			log.Println("[ERR]: failed to unmarshal source")
			return models.Campaign{}, err
		}
		campaign.ID = res.Hits.Hits[0].Id
		return campaign, nil
	}

	return models.Campaign{}, errors.New("campaign not found")
}
