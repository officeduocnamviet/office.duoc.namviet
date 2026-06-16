package employees

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreateEmployeeRequest struct {
	UserID       *string      `json:"user_id"`
	Code         *string      `json:"code"`
	DepartmentID *string      `json:"department_id"`
	PositionID   *string      `json:"position_id"`
	Type         *string      `json:"type"`
	HireDate     *time.Time   `json:"hire_date"`
	BaseSalary   *float64     `json:"base_salary"`
	SalaryType   *string      `json:"salary_type"`
	InsuranceNo  *string      `json:"insurance_no"`
	TaxCode      *string      `json:"tax_code"`
	BankAccount  *roles.JSONB `json:"bank_account"`
}

type UpdateEmployeeRequest struct {
	Code            *string      `json:"code"`
	DepartmentID    *string      `json:"department_id"`
	PositionID      *string      `json:"position_id"`
	Type            *string      `json:"type"`
	Status          *string      `json:"status"`
	HireDate        *time.Time   `json:"hire_date"`
	TerminationDate *time.Time   `json:"termination_date"`
	BaseSalary      *float64     `json:"base_salary"`
	SalaryType      *string      `json:"salary_type"`
	InsuranceNo     *string      `json:"insurance_no"`
	TaxCode         *string      `json:"tax_code"`
	BankAccount     *roles.JSONB `json:"bank_account"`
}
