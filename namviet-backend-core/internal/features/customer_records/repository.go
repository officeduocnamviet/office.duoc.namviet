package customer_records

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Vaccination Records
func GetAllVaccinationRecords() ([]CustomerVaccinationRecord, error) {
	var results []CustomerVaccinationRecord
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetVaccinationRecordByID(id string) (*CustomerVaccinationRecord, error) {
	var result CustomerVaccinationRecord
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("vaccination record not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateVaccinationRecord(data *CustomerVaccinationRecord) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateVaccinationRecord(data *CustomerVaccinationRecord) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteVaccinationRecord(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&CustomerVaccinationRecord{}).Error
}

// Vouchers
func GetAllCustomerVouchers() ([]CustomerVoucher, error) {
	var results []CustomerVoucher
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetCustomerVoucherByID(id string) (*CustomerVoucher, error) {
	var result CustomerVoucher
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("customer voucher not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateCustomerVoucher(data *CustomerVoucher) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateCustomerVoucher(data *CustomerVoucher) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteCustomerVoucher(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&CustomerVoucher{}).Error
}
