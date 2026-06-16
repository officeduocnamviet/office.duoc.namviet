package fund_accounts

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllFundAccounts() ([]FundAccount, error) {
	var results []FundAccount
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetFundAccountByID(id string) (*FundAccount, error) {
	var result FundAccount
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("fund account not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateFundAccount(data *FundAccount) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateFundAccount(data *FundAccount) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteFundAccount(id string) error {
	db := supabase.DB
	return db.Model(&FundAccount{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
