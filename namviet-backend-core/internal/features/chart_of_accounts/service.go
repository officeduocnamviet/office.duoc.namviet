package chart_of_accounts

import "time"

func GetAllChartOfAccountsService() ([]ChartOfAccount, error) {
	return GetAllChartOfAccounts()
}

func GetChartOfAccountByIDService(id string) (*ChartOfAccount, error) {
	return GetChartOfAccountByID(id)
}

func CreateChartOfAccountService(req CreateChartOfAccountRequest) (*ChartOfAccount, error) {
	allowPosting := true
	if req.AllowPosting != nil {
		allowPosting = *req.AllowPosting
	}

	coa := &ChartOfAccount{
		AccountCode:  req.AccountCode,
		Name:         req.Name,
		ParentID:     req.ParentID,
		Type:         req.Type,
		BalanceType:  req.BalanceType,
		Status:       "active",
		AllowPosting: allowPosting,
	}

	if err := CreateChartOfAccount(coa); err != nil {
		return nil, err
	}
	return coa, nil
}

func UpdateChartOfAccountService(id string, req UpdateChartOfAccountRequest) (*ChartOfAccount, error) {
	coa, err := GetChartOfAccountByID(id)
	if err != nil {
		return nil, err
	}

	if req.AccountCode != nil {
		coa.AccountCode = *req.AccountCode
	}
	if req.Name != nil {
		coa.Name = *req.Name
	}
	if req.ParentID != nil {
		coa.ParentID = req.ParentID
	}
	if req.Type != nil {
		coa.Type = *req.Type
	}
	if req.BalanceType != nil {
		coa.BalanceType = *req.BalanceType
	}
	if req.Status != nil {
		coa.Status = *req.Status
	}
	if req.AllowPosting != nil {
		coa.AllowPosting = *req.AllowPosting
	}
	
	now := time.Now()
	coa.UpdatedAt = &now

	if err := UpdateChartOfAccount(coa); err != nil {
		return nil, err
	}
	return coa, nil
}

func DeleteChartOfAccountService(id string) error {
	return DeleteChartOfAccount(id)
}
