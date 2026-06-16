package employees

import "time"

func GetAllEmployeesService() ([]Employee, error) {
	return GetAllEmployees()
}

func GetEmployeeByIDService(id string) (*Employee, error) {
	return GetEmployeeByID(id)
}

func CreateEmployeeService(req CreateEmployeeRequest) (*Employee, error) {
	empType := "full_time"
	if req.Type != nil {
		empType = *req.Type
	}

	salaryType := "monthly"
	if req.SalaryType != nil {
		salaryType = *req.SalaryType
	}

	emp := &Employee{
		UserID:       req.UserID,
		Code:         req.Code,
		DepartmentID: req.DepartmentID,
		PositionID:   req.PositionID,
		Type:         empType,
		Status:       "active",
		HireDate:     req.HireDate,
		BaseSalary:   req.BaseSalary,
		SalaryType:   salaryType,
		InsuranceNo:  req.InsuranceNo,
		TaxCode:      req.TaxCode,
	}

	if req.BankAccount != nil {
		emp.BankAccount = *req.BankAccount
	}

	if err := CreateEmployee(emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func UpdateEmployeeService(id string, req UpdateEmployeeRequest) (*Employee, error) {
	emp, err := GetEmployeeByID(id)
	if err != nil {
		return nil, err
	}

	if req.Code != nil {
		emp.Code = req.Code
	}
	if req.DepartmentID != nil {
		emp.DepartmentID = req.DepartmentID
	}
	if req.PositionID != nil {
		emp.PositionID = req.PositionID
	}
	if req.Type != nil {
		emp.Type = *req.Type
	}
	if req.Status != nil {
		emp.Status = *req.Status
	}
	if req.HireDate != nil {
		emp.HireDate = req.HireDate
	}
	if req.TerminationDate != nil {
		emp.TerminationDate = req.TerminationDate
	}
	if req.BaseSalary != nil {
		emp.BaseSalary = req.BaseSalary
	}
	if req.SalaryType != nil {
		emp.SalaryType = *req.SalaryType
	}
	if req.InsuranceNo != nil {
		emp.InsuranceNo = req.InsuranceNo
	}
	if req.TaxCode != nil {
		emp.TaxCode = req.TaxCode
	}
	if req.BankAccount != nil {
		emp.BankAccount = *req.BankAccount
	}
	
	now := time.Now()
	emp.UpdatedAt = &now

	if err := UpdateEmployee(emp); err != nil {
		return nil, err
	}
	return emp, nil
}

func DeleteEmployeeService(id string) error {
	return DeleteEmployee(id)
}
