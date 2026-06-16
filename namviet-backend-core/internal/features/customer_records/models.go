package customer_records

import (
	"time"
)

// CustomerVaccinationRecord represents the customer_vaccination_records table
type CustomerVaccinationRecord struct {
	ID               string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CustomerID       string     `gorm:"type:uuid;not null" json:"customer_id"`
	VaccineName      string     `gorm:"type:text;not null" json:"vaccine_name"`
	DoseNumber       int        `gorm:"type:integer;default:1;not null" json:"dose_number"`
	VaccinationDate  time.Time  `gorm:"type:date;not null" json:"vaccination_date"`
	NextDueDate      *time.Time `gorm:"type:date" json:"next_due_date,omitempty"`
	AdministeredBy   *string    `gorm:"type:text" json:"administered_by,omitempty"`
	Notes            *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt        *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}

// CustomerVoucher represents the customer_vouchers table
type CustomerVoucher struct {
	ID          string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	CustomerID  string     `gorm:"type:uuid;not null" json:"customer_id"`
	PromotionID int64      `gorm:"type:bigint;not null" json:"promotion_id"`
	VoucherCode string     `gorm:"type:text;not null" json:"voucher_code"`
	IsUsed      bool       `gorm:"type:boolean;default:false" json:"is_used"`
	UsedAt      *time.Time `gorm:"type:timestamp with time zone" json:"used_at,omitempty"`
	CreatedAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
