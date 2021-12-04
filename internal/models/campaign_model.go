package models

type Campaign struct {
	ID string `json:"id"`

	UserID string `json:"user_id"`

	CampaignName      string   `json:"campaign_name"`
	CampaignStartDate string   `json:"campaign_start_date"`
	DailyBudget       float32  `json:"daily_budget"`
	MatchType         string   `json:"match_type"`
	Bid               float32  `json:"bid"`
	SKU               string   `json:"sku"`
	TotalKeywords     uint     `json:"total_keywords"`
	Keywords          []string `json:"keywords"`

	NegativeMatchType string   `json:"negative_match_type"`
	NegativeKeywords  []string `json:"negative_keywords"`
}
