package promotions

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllPromotions() ([]Promotion, error) {
	var results []Promotion
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetPromotionByID(id string) (*Promotion, error) {
	var result Promotion
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promotion not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreatePromotion(data *Promotion) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdatePromotion(data *Promotion) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeletePromotion(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&Promotion{}).Error
}
