package employment_contracts

import "time"

func GetAllEmploymentContractsService() ([]EmploymentContract, error) {
	return GetAllEmploymentContracts()
}

func GetEmploymentContractByIDService(id string) (*EmploymentContract, error) {
	return GetEmploymentContractByID(id)
}

func CreateEmploymentContractService(req CreateEmploymentContractRequest) (*EmploymentContract, error) {
	contract := &EmploymentContract{
		UserID:                   req.UserID,
		ContractCode:             req.ContractCode,
		BaseSalary:               req.BaseSalary,
		StandardWorkingDays:      req.StandardWorkingDays,
		KPIConversionRate:        req.KPIConversionRate,
		CommissionRatePercent:    req.CommissionRatePercent,
		TaxDeductionAmount:       req.TaxDeductionAmount,
		InsuranceDeductionAmount: req.InsuranceDeductionAmount,
		ValidFrom:                req.ValidFrom,
		ValidTo:                  req.ValidTo,
		Status:                   "active",
	}

	if err := CreateEmploymentContract(contract); err != nil {
		return nil, err
	}
	return contract, nil
}

func UpdateEmploymentContractService(id string, req UpdateEmploymentContractRequest) (*EmploymentContract, error) {
	contract, err := GetEmploymentContractByID(id)
	if err != nil {
		return nil, err
	}

	if req.ContractCode != nil {
		contract.ContractCode = *req.ContractCode
	}
	if req.BaseSalary != nil {
		contract.BaseSalary = *req.BaseSalary
	}
	if req.StandardWorkingDays != nil {
		contract.StandardWorkingDays = *req.StandardWorkingDays
	}
	if req.KPIConversionRate != nil {
		contract.KPIConversionRate = *req.KPIConversionRate
	}
	if req.CommissionRatePercent != nil {
		contract.CommissionRatePercent = *req.CommissionRatePercent
	}
	if req.TaxDeductionAmount != nil {
		contract.TaxDeductionAmount = *req.TaxDeductionAmount
	}
	if req.InsuranceDeductionAmount != nil {
		contract.InsuranceDeductionAmount = *req.InsuranceDeductionAmount
	}
	if req.ValidFrom != nil {
		contract.ValidFrom = *req.ValidFrom
	}
	if req.ValidTo != nil {
		contract.ValidTo = req.ValidTo
	}
	if req.Status != nil {
		contract.Status = *req.Status
	}

	now := time.Now()
	contract.UpdatedAt = &now

	if err := UpdateEmploymentContract(contract); err != nil {
		return nil, err
	}
	return contract, nil
}

func DeleteEmploymentContractService(id string) error {
	return DeleteEmploymentContract(id)
}
