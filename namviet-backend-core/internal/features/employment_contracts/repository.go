package employment_contracts

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllEmploymentContracts() ([]EmploymentContract, error) {
	var results []EmploymentContract
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetEmploymentContractByID(id string) (*EmploymentContract, error) {
	var result EmploymentContract
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("employment contract not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateEmploymentContract(data *EmploymentContract) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateEmploymentContract(data *EmploymentContract) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteEmploymentContract(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&EmploymentContract{}).Error
}
