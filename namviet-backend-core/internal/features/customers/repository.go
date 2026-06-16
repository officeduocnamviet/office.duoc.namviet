package customers

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllCustomers() ([]Customer, error) {
	var results []Customer
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetCustomerByID(id int64) (*Customer, error) {
	var result Customer
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateCustomer(data *Customer) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateCustomer(data *Customer) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteCustomer(id int64) error {
	db := supabase.DB
	return db.Model(&Customer{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
