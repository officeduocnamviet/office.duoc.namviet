package warehouses

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllWarehouses() ([]Warehouse, error) {
	var results []Warehouse
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetWarehouseByID(id int64) (*Warehouse, error) {
	var result Warehouse
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("warehouse not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateWarehouse(data *Warehouse) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateWarehouse(data *Warehouse) error {
	db := supabase.DB
	return db.Save(data).Error
}
