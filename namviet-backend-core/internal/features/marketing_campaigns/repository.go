package marketing_campaigns

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllMarketingCampaigns() ([]MarketingCampaign, error) {
	var results []MarketingCampaign
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetMarketingCampaignByID(id string) (*MarketingCampaign, error) {
	var result MarketingCampaign
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("marketing campaign not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateMarketingCampaign(data *MarketingCampaign) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateMarketingCampaign(data *MarketingCampaign) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteMarketingCampaign(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&MarketingCampaign{}).Error
}
