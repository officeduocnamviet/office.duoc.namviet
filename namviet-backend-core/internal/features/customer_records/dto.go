package customer_records

import "time"

// Vaccination Record DTOs
type CreateVaccinationRecordRequest struct {
	CustomerID       string     `json:"customer_id" binding:"required"`
	VaccineName      string     `json:"vaccine_name" binding:"required"`
	DoseNumber       int        `json:"dose_number"`
	VaccinationDate  time.Time  `json:"vaccination_date" binding:"required"`
	NextDueDate      *time.Time `json:"next_due_date"`
	AdministeredBy   *string    `json:"administered_by"`
	Notes            *string    `json:"notes"`
}

type UpdateVaccinationRecordRequest struct {
	VaccineName      *string    `json:"vaccine_name"`
	DoseNumber       *int       `json:"dose_number"`
	VaccinationDate  *time.Time `json:"vaccination_date"`
	NextDueDate      *time.Time `json:"next_due_date"`
	AdministeredBy   *string    `json:"administered_by"`
	Notes            *string    `json:"notes"`
}

// Voucher DTOs
type CreateCustomerVoucherRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	PromotionID int64  `json:"promotion_id" binding:"required"`
	VoucherCode string `json:"voucher_code" binding:"required"`
}

type UpdateCustomerVoucherRequest struct {
	IsUsed *bool      `json:"is_used"`
	UsedAt *time.Time `json:"used_at"`
}
