package accounting_journals

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllAccountingJournals() ([]AccountingJournal, error) {
	var results []AccountingJournal
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetAccountingJournalByID(id string) (*AccountingJournal, error) {
	var result AccountingJournal
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("accounting journal not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateAccountingJournal(data *AccountingJournal) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateAccountingJournal(data *AccountingJournal) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteAccountingJournal(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&AccountingJournal{}).Error
}
