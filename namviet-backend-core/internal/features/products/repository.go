package products

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllProducts() ([]Product, error) {
	var results []Product
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetProductByID(id int64) (*Product, error) {
	var result Product
	db := supabase.DB
	// Later we will preload units here
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateProduct(data *Product) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateProduct(data *Product) error {
	db := supabase.DB
	return db.Save(data).Error
}
