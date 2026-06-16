package fund_accounts

import "time"

func GetAllFundAccountsService() ([]FundAccount, error) {
	return GetAllFundAccounts()
}

func GetFundAccountByIDService(id string) (*FundAccount, error) {
	return GetFundAccountByID(id)
}

func CreateFundAccountService(req CreateFundAccountRequest) (*FundAccount, error) {
	currency := "VND"
	if req.Currency != nil {
		currency = *req.Currency
	}

	fa := &FundAccount{
		Name:           req.Name,
		Type:           req.Type,
		Location:       req.Location,
		AccountNumber:  req.AccountNumber,
		BankID:         req.BankID,
		InitialBalance: req.InitialBalance,
		Balance:        req.Balance,
		Currency:       currency,
		Status:         "active",
		Description:    req.Description,
		AccountID:      req.AccountID,
	}

	if req.BankInfo != nil {
		fa.BankInfo = *req.BankInfo
	}

	if err := CreateFundAccount(fa); err != nil {
		return nil, err
	}
	return fa, nil
}

func UpdateFundAccountService(id string, req UpdateFundAccountRequest) (*FundAccount, error) {
	fa, err := GetFundAccountByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		fa.Name = *req.Name
	}
	if req.Type != nil {
		fa.Type = *req.Type
	}
	if req.Location != nil {
		fa.Location = req.Location
	}
	if req.AccountNumber != nil {
		fa.AccountNumber = req.AccountNumber
	}
	if req.BankID != nil {
		fa.BankID = req.BankID
	}
	if req.Balance != nil {
		fa.Balance = *req.Balance
	}
	if req.Currency != nil {
		fa.Currency = *req.Currency
	}
	if req.Status != nil {
		fa.Status = *req.Status
	}
	if req.BankInfo != nil {
		fa.BankInfo = *req.BankInfo
	}
	if req.Description != nil {
		fa.Description = req.Description
	}
	if req.AccountID != nil {
		fa.AccountID = req.AccountID
	}
	
	now := time.Now()
	fa.UpdatedAt = &now

	if err := UpdateFundAccount(fa); err != nil {
		return nil, err
	}
	return fa, nil
}

func DeleteFundAccountService(id string) error {
	return DeleteFundAccount(id)
}
