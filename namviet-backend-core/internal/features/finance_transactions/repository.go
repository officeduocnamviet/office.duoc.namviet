package finance_transactions

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllFinanceTransactions() ([]FinanceTransaction, error) {
	var results []FinanceTransaction
	db := supabase.DB
	err := db.Where("deleted_at IS NULL").Find(&results).Error
	return results, err
}

func GetFinanceTransactionByID(id string) (*FinanceTransaction, error) {
	var result FinanceTransaction
	db := supabase.DB
	err := db.First(&result, "id = ? AND deleted_at IS NULL", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("finance transaction not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateFinanceTransaction(data *FinanceTransaction) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateFinanceTransaction(data *FinanceTransaction) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteFinanceTransaction(id string) error {
	db := supabase.DB
	return db.Model(&FinanceTransaction{}).Where("id = ?", id).Update("deleted_at", gorm.Expr("now()")).Error
}
