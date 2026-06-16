package employees

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// Employee represents the employees table
type Employee struct {
	ID              string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID          *string     `gorm:"type:uuid" json:"user_id,omitempty"`
	Code            *string     `gorm:"type:text" json:"code,omitempty"`
	DepartmentID    *string     `gorm:"type:uuid" json:"department_id,omitempty"`
	PositionID      *string     `gorm:"type:uuid" json:"position_id,omitempty"`
	Type            string      `gorm:"type:text;default:'full_time'" json:"type"`
	Status          string      `gorm:"type:text;default:'active'" json:"status"`
	HireDate        *time.Time  `gorm:"type:date" json:"hire_date,omitempty"`
	TerminationDate *time.Time  `gorm:"type:date" json:"termination_date,omitempty"`
	BaseSalary      *float64    `gorm:"type:numeric" json:"base_salary,omitempty"`
	SalaryType      string      `gorm:"type:text;default:'monthly'" json:"salary_type"`
	InsuranceNo     *string     `gorm:"type:text" json:"insurance_no,omitempty"`
	TaxCode         *string     `gorm:"type:text" json:"tax_code,omitempty"`
	BankAccount     roles.JSONB `gorm:"type:jsonb;default:'{}'::jsonb" json:"bank_account,omitempty"`
	CreatedAt       *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt       *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt       *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
