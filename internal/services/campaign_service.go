package services

import (
	"fmt"
	"github.com/rustingoff/excel_vue_go/internal/models"
	"github.com/rustingoff/excel_vue_go/internal/repositories"
	"github.com/xuri/excelize/v2"
	"log"
	"time"
)

type CampaignService interface {
	CreateCampaign(campaign models.Campaign) error
	DeleteCamapign(string) error

	GetAllCampaigns() ([]models.Campaign, error)
	GetCampaignById(string) (models.Campaign, error)

	ExportCampaignExcel(string) (*excelize.File, error)
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

func (service *campaignService) ExportCampaignExcel(id string) (*excelize.File, error) {
	campaign, err := service.repo.GetCampaignById(id)
	if err != nil {
		return nil, err
	}

	f, err := excelize.OpenFile("./static/template.xlsx")
	if err != nil {
		log.Println("[ERR]: failed to open excel file, ", err.Error())
		return nil, err
	}

	var matchType string
	campaignsCount := len(campaign.Keywords) / int(campaign.TotalKeywords)
	count := 0
	restKeyCount := 0

	for j := 0; j < campaignsCount; j++ {
		if j > 0 {
			matchType = " - " + campaign.MatchType + fmt.Sprint(j)
		} else {
			matchType = " - " + campaign.MatchType
		}
		c, err := service.writeExportCampaign(f, (j*6)+count, campaign, matchType, campaign.Keywords[int(campaign.TotalKeywords)*j:int(campaign.TotalKeywords)*(j+1)])
		if err != nil {
			return nil, err
		}
		count += int(campaign.TotalKeywords) + c
		if j == campaignsCount-1 && campaignsCount > 1 {
			restKeyCount = len(campaign.Keywords) % int(campaign.TotalKeywords)
			if restKeyCount > 0 {
				_, err = service.writeExportCampaign(f, ((j+1)*6)+count, campaign, matchType, campaign.Keywords[int(campaign.TotalKeywords)*(j+1):len(campaign.Keywords)])
				if err != nil {
					return nil, err
				}
			}
			count += restKeyCount + 6
		}
	}

	err = f.SaveAs("./static/exports/" + campaign.ID + ".xlsx")
	if err != nil {
		log.Println("[ERR]: failed to save file, ", err.Error())
		return nil, err
	}

	return f, nil
}

func (service *campaignService) writeExportCampaign(f *excelize.File, count int, campaign models.Campaign, nameExact string, keywords []string) (int, error) {
	err := f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", count+2), "Campaign")
	if err != nil {
		return 0, err
	}

	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", count+2), campaign.CampaignName+nameExact)
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("E%d", count+2), campaign.DailyBudget)
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("G%d", count+2), campaign.CampaignStartDate)
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("I%d", count+2), "Manual")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("P%d", count+2), "enabled")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("Z%d", count+2), "Dynamic bidding (down only)")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("AA%d", count+2), "All")
	if err != nil {
		return 0, err
	}

	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("AA%d", count+3), "Top of search (page 1)")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("AB%d", count+3), "0%")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", count+3), "Campaign By Placement")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", count+3), campaign.CampaignName+nameExact)
	if err != nil {
		return 0, err
	}

	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("AA%d", count+4), "Rest of search")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", count+4), "Campaign By Placement")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", count+4), campaign.CampaignName+nameExact)
	if err != nil {
		return 0, err
	}

	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("AA%d", count+5), "Product pages")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("AB%d", count+5), "0%")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", count+5), "Campaign By Placement")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", count+5), campaign.CampaignName+nameExact)
	if err != nil {
		return 0, err
	}

	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", count+6), "Ad Group")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("K%d", count+6), campaign.Bid)
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("P%d", count+6), "enabled")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("Q%d", count+6), "enabled")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("J%d", count+6), "Ad Group 1")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", count+6), campaign.CampaignName+nameExact)
	if err != nil {
		return 0, err
	}

	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", count+7), "Ad")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("J%d", count+7), "Ad Group 1")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", count+7), campaign.CampaignName+nameExact)
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("O%d", count+7), campaign.SKU)
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("P%d", count+7), "enabled")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("Q%d", count+7), "enabled")
	if err != nil {
		return 0, err
	}
	err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("R%d", count+7), "enabled")
	if err != nil {
		return 0, err
	}

	for j, keyword := range keywords {
		err := f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", j+8+count), "Keyword")
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("J%d", j+8+count), "Ad Group 1")
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", j+8+count), campaign.CampaignName+nameExact)
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("L%d", j+8+count), keyword)
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("P%d", j+8+count), "enabled")
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("Q%d", j+8+count), "enabled")
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("R%d", j+8+count), "enabled")
		if err != nil {
			return 0, err
		}
		err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("N%d", j+8+count), campaign.MatchType)
		if err != nil {
			return 0, err
		}

		if campaign.NegativeMatchType == "campaign negative phrase" || campaign.NegativeMatchType == "campaign negative exact" {
			for k, negativeKeyword := range campaign.NegativeKeywords {
				err := f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("B%d", j+k+9+count), "Keyword")
				if err != nil {
					return 0, err
				}
				err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("D%d", j+k+9+count), campaign.CampaignName+nameExact)
				if err != nil {
					return 0, err
				}
				err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("L%d", j+k+9+count), negativeKeyword)
				if err != nil {
					return 0, err
				}
				err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("N%d", j+k+9+count), campaign.NegativeMatchType)
				if err != nil {
					return 0, err
				}
				err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("P%d", j+k+9+count), "enabled")
				if err != nil {
					return 0, err
				}
				err = f.SetCellValue("Sponsored Products Campaigns", fmt.Sprintf("R%d", j+k+9+count), "enabled")
				if err != nil {
					return 0, err
				}
			}
		}
	}
	return len(campaign.NegativeKeywords), nil
}
