package payrolls

import (
	"time"
)

// Payroll represents the payrolls table
type Payroll struct {
	ID              string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	EmployeeID      string     `gorm:"type:uuid;not null" json:"employee_id"`
	PeriodMonth     int        `gorm:"type:integer;not null" json:"period_month"`
	PeriodYear      int        `gorm:"type:integer;not null" json:"period_year"`
	BaseSalary      float64    `gorm:"type:numeric;not null" json:"base_salary"`
	Allowances      *float64   `gorm:"type:numeric;default:0" json:"allowances,omitempty"`
	OvertimePay     *float64   `gorm:"type:numeric;default:0" json:"overtime_pay,omitempty"`
	Bonuses         *float64   `gorm:"type:numeric;default:0" json:"bonuses,omitempty"`
	Deductions      *float64   `gorm:"type:numeric;default:0" json:"deductions,omitempty"`
	TaxAmount       *float64   `gorm:"type:numeric;default:0" json:"tax_amount,omitempty"`
	InsuranceAmount *float64   `gorm:"type:numeric;default:0" json:"insurance_amount,omitempty"`
	NetSalary       float64    `gorm:"type:numeric;not null" json:"net_salary"`
	Status          string     `gorm:"type:text;default:'draft'" json:"status"`
	PaymentDate     *time.Time `gorm:"type:timestamp with time zone" json:"payment_date,omitempty"`
	CreatedAt       *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt       *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}
