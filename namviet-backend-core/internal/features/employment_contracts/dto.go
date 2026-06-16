package employment_contracts

import "time"

type CreateEmploymentContractRequest struct {
	UserID                   string     `json:"user_id" binding:"required"`
	ContractCode             string     `json:"contract_code" binding:"required"`
	BaseSalary               float64    `json:"base_salary"`
	StandardWorkingDays      int        `json:"standard_working_days"`
	KPIConversionRate        float64    `json:"kpi_conversion_rate"`
	CommissionRatePercent    float64    `json:"commission_rate_percent"`
	TaxDeductionAmount       float64    `json:"tax_deduction_amount"`
	InsuranceDeductionAmount float64    `json:"insurance_deduction_amount"`
	ValidFrom                time.Time  `json:"valid_from" binding:"required"`
	ValidTo                  *time.Time `json:"valid_to"`
}

type UpdateEmploymentContractRequest struct {
	ContractCode             *string    `json:"contract_code"`
	BaseSalary               *float64   `json:"base_salary"`
	StandardWorkingDays      *int       `json:"standard_working_days"`
	KPIConversionRate        *float64   `json:"kpi_conversion_rate"`
	CommissionRatePercent    *float64   `json:"commission_rate_percent"`
	TaxDeductionAmount       *float64   `json:"tax_deduction_amount"`
	InsuranceDeductionAmount *float64   `json:"insurance_deduction_amount"`
	ValidFrom                *time.Time `json:"valid_from"`
	ValidTo                  *time.Time `json:"valid_to"`
	Status                   *string    `json:"status"`
}
