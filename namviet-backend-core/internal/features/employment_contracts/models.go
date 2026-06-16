package employment_contracts

import (
	"time"
)

// EmploymentContract represents the employment_contracts table
type EmploymentContract struct {
	ID                       string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID                   string     `gorm:"type:uuid;not null" json:"user_id"`
	ContractCode             string     `gorm:"type:text;not null" json:"contract_code"`
	BaseSalary               float64    `gorm:"type:numeric;default:0;not null" json:"base_salary"`
	StandardWorkingDays      int        `gorm:"type:integer;default:26;not null" json:"standard_working_days"`
	KPIConversionRate        float64    `gorm:"type:numeric;default:0" json:"kpi_conversion_rate"`
	CommissionRatePercent    float64    `gorm:"type:numeric;default:0" json:"commission_rate_percent"`
	TaxDeductionAmount       float64    `gorm:"type:numeric;default:0" json:"tax_deduction_amount"`
	InsuranceDeductionAmount float64    `gorm:"type:numeric;default:0" json:"insurance_deduction_amount"`
	ValidFrom                time.Time  `gorm:"type:date;not null" json:"valid_from"`
	ValidTo                  *time.Time `gorm:"type:date" json:"valid_to,omitempty"`
	Status                   string     `gorm:"type:text;default:'active';not null" json:"status"`
	CreatedAt                *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt                *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}
