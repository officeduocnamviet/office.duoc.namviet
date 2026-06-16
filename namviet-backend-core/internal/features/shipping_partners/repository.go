package shipping_partners

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllShippingPartners() ([]ShippingPartner, error) {
	var results []ShippingPartner
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetShippingPartnerByID(id string) (*ShippingPartner, error) {
	var result ShippingPartner
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("shipping partner not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateShippingPartner(data *ShippingPartner) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateShippingPartner(data *ShippingPartner) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteShippingPartner(id string) error {
	db := supabase.DB
	return db.Model(&ShippingPartner{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
