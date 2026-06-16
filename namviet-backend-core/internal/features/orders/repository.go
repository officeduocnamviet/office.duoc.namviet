package orders

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllOrders() ([]Order, error) {
	var results []Order
	db := supabase.DB
	err := db.Preload("Items").Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetOrderByID(id string) (*Order, error) {
	var result Order
	db := supabase.DB
	err := db.Preload("Items").First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateOrderWithItems(order *Order) error {
	db := supabase.DB
	// Use transaction
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		return nil
	})
}

func UpdateOrder(data *Order) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteOrder(id string) error {
	db := supabase.DB
	return db.Model(&Order{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
