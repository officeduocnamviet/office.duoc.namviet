package product_units

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetUnitsByProductID(productID int64) ([]ProductUnit, error) {
	var results []ProductUnit
	db := supabase.DB
	err := db.Where("product_id = ? AND deleted_at IS NULL", productID).Find(&results).Error
	return results, err
}

func GetProductUnitByID(id int64) (*ProductUnit, error) {
	var result ProductUnit
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product unit not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateProductUnit(data *ProductUnit) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateProductUnit(data *ProductUnit) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteProductUnit(id int64) error {
	db := supabase.DB
	// Soft delete
	return db.Model(&ProductUnit{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
