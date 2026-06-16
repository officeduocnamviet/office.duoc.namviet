package payrolls

import "time"

func GetAllPayrollsService() ([]Payroll, error) {
	return GetAllPayrolls()
}

func GetPayrollByIDService(id string) (*Payroll, error) {
	return GetPayrollByID(id)
}

func CreatePayrollService(req CreatePayrollRequest) (*Payroll, error) {
	payroll := &Payroll{
		EmployeeID:      req.EmployeeID,
		PeriodMonth:     req.PeriodMonth,
		PeriodYear:      req.PeriodYear,
		BaseSalary:      req.BaseSalary,
		Allowances:      req.Allowances,
		OvertimePay:     req.OvertimePay,
		Bonuses:         req.Bonuses,
		Deductions:      req.Deductions,
		TaxAmount:       req.TaxAmount,
		InsuranceAmount: req.InsuranceAmount,
		NetSalary:       req.NetSalary,
		Status:          "draft",
	}

	if err := CreatePayroll(payroll); err != nil {
		return nil, err
	}
	return payroll, nil
}

func UpdatePayrollService(id string, req UpdatePayrollRequest) (*Payroll, error) {
	payroll, err := GetPayrollByID(id)
	if err != nil {
		return nil, err
	}

	if req.BaseSalary != nil {
		payroll.BaseSalary = *req.BaseSalary
	}
	if req.Allowances != nil {
		payroll.Allowances = req.Allowances
	}
	if req.OvertimePay != nil {
		payroll.OvertimePay = req.OvertimePay
	}
	if req.Bonuses != nil {
		payroll.Bonuses = req.Bonuses
	}
	if req.Deductions != nil {
		payroll.Deductions = req.Deductions
	}
	if req.TaxAmount != nil {
		payroll.TaxAmount = req.TaxAmount
	}
	if req.InsuranceAmount != nil {
		payroll.InsuranceAmount = req.InsuranceAmount
	}
	if req.NetSalary != nil {
		payroll.NetSalary = *req.NetSalary
	}
	if req.Status != nil {
		payroll.Status = *req.Status
	}
	if req.PaymentDate != nil {
		payroll.PaymentDate = req.PaymentDate
	}
	
	now := time.Now()
	payroll.UpdatedAt = &now

	if err := UpdatePayroll(payroll); err != nil {
		return nil, err
	}
	return payroll, nil
}

func DeletePayrollService(id string) error {
	return DeletePayroll(id)
}
