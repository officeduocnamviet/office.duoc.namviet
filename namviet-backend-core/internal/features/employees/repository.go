package employees

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllEmployees() ([]Employee, error) {
	var results []Employee
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetEmployeeByID(id string) (*Employee, error) {
	var result Employee
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employee not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateEmployee(data *Employee) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateEmployee(data *Employee) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteEmployee(id string) error {
	db := supabase.DB
	return db.Model(&Employee{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
