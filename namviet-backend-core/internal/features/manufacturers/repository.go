package manufacturers

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllManufacturers() ([]Manufacturer, error) {
	var results []Manufacturer
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetManufacturerByID(id int64) (*Manufacturer, error) {
	var result Manufacturer
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("manufacturer not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateManufacturer(data *Manufacturer) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateManufacturer(data *Manufacturer) error {
	db := supabase.DB
	return db.Save(data).Error
}
