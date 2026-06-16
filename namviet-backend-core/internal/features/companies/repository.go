package companies

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Companies
func GetAllCompanies() ([]Company, error) {
	var results []Company
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetCompanyByID(id string) (*Company, error) {
	var result Company
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("company not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateCompany(data *Company) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateCompany(data *Company) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteCompany(id string) error {
	db := supabase.DB
	return db.Model(&Company{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}

// Branches
func GetAllBranches() ([]Branch, error) {
	var results []Branch
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetBranchByID(id string) (*Branch, error) {
	var result Branch
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("branch not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateBranch(data *Branch) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateBranch(data *Branch) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteBranch(id string) error {
	db := supabase.DB
	return db.Model(&Branch{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
