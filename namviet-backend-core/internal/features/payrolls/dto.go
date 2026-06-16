package payrolls

import "time"

type CreatePayrollRequest struct {
	EmployeeID      string     `json:"employee_id" binding:"required"`
	PeriodMonth     int        `json:"period_month" binding:"required"`
	PeriodYear      int        `json:"period_year" binding:"required"`
	BaseSalary      float64    `json:"base_salary" binding:"required"`
	Allowances      *float64   `json:"allowances"`
	OvertimePay     *float64   `json:"overtime_pay"`
	Bonuses         *float64   `json:"bonuses"`
	Deductions      *float64   `json:"deductions"`
	TaxAmount       *float64   `json:"tax_amount"`
	InsuranceAmount *float64   `json:"insurance_amount"`
	NetSalary       float64    `json:"net_salary" binding:"required"`
}

type UpdatePayrollRequest struct {
	BaseSalary      *float64   `json:"base_salary"`
	Allowances      *float64   `json:"allowances"`
	OvertimePay     *float64   `json:"overtime_pay"`
	Bonuses         *float64   `json:"bonuses"`
	Deductions      *float64   `json:"deductions"`
	TaxAmount       *float64   `json:"tax_amount"`
	InsuranceAmount *float64   `json:"insurance_amount"`
	NetSalary       *float64   `json:"net_salary"`
	Status          *string    `json:"status"`
	PaymentDate     *time.Time `json:"payment_date"`
}
