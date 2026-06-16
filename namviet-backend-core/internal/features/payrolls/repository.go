package payrolls

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllPayrolls() ([]Payroll, error) {
	var results []Payroll
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetPayrollByID(id string) (*Payroll, error) {
	var result Payroll
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("payroll not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreatePayroll(data *Payroll) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdatePayroll(data *Payroll) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeletePayroll(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&Payroll{}).Error
}
