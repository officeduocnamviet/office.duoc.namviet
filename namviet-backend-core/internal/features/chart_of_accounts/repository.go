package chart_of_accounts

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllChartOfAccounts() ([]ChartOfAccount, error) {
	var results []ChartOfAccount
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetChartOfAccountByID(id string) (*ChartOfAccount, error) {
	var result ChartOfAccount
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chart of account not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateChartOfAccount(data *ChartOfAccount) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateChartOfAccount(data *ChartOfAccount) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteChartOfAccount(id string) error {
	db := supabase.DB
	return db.Model(&ChartOfAccount{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
