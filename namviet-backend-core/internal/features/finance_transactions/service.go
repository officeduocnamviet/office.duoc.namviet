package finance_transactions

import "time"

func GetAllFinanceTransactionsService() ([]FinanceTransaction, error) {
	return GetAllFinanceTransactions()
}

func GetFinanceTransactionByIDService(id string) (*FinanceTransaction, error) {
	return GetFinanceTransactionByID(id)
}

func CreateFinanceTransactionService(req CreateFinanceTransactionRequest) (*FinanceTransaction, error) {
	businessType := "other"
	if req.BusinessType != nil {
		businessType = *req.BusinessType
	}
	
	bookType := "BOTH"
	if req.BookType != nil {
		bookType = *req.BookType
	}

	ft := &FinanceTransaction{
		Code:             req.Code,
		Flow:             req.Flow,
		BusinessType:     businessType,
		CategoryID:       req.CategoryID,
		Amount:           req.Amount,
		FundAccountID:    req.FundAccountID,
		PartnerType:      req.PartnerType,
		PartnerID:        req.PartnerID,
		PartnerNameCache: req.PartnerNameCache,
		RefType:          req.RefType,
		RefID:            req.RefID,
		Description:      req.Description,
		EvidenceURL:      req.EvidenceURL,
		Status:           "pending",
		RefAdvanceID:     req.RefAdvanceID,
		BankReferenceID:  req.BankReferenceID,
		BookType:         bookType,
		CreatedBy:        req.CreatedBy,
	}

	if !req.TransactionDate.IsZero() {
		ft.TransactionDate = req.TransactionDate
	} else {
		ft.TransactionDate = time.Now()
	}

	if req.CashTally != nil {
		ft.CashTally = *req.CashTally
	}
	if req.TargetBankInfo != nil {
		ft.TargetBankInfo = *req.TargetBankInfo
	}

	if err := CreateFinanceTransaction(ft); err != nil {
		return nil, err
	}
	return ft, nil
}

func UpdateFinanceTransactionService(id string, req UpdateFinanceTransactionRequest) (*FinanceTransaction, error) {
	ft, err := GetFinanceTransactionByID(id)
	if err != nil {
		return nil, err
	}

	if req.TransactionDate != nil {
		ft.TransactionDate = *req.TransactionDate
	}
	if req.Flow != nil {
		ft.Flow = *req.Flow
	}
	if req.BusinessType != nil {
		ft.BusinessType = *req.BusinessType
	}
	if req.CategoryID != nil {
		ft.CategoryID = req.CategoryID
	}
	if req.Amount != nil {
		ft.Amount = *req.Amount
	}
	if req.FundAccountID != nil {
		ft.FundAccountID = *req.FundAccountID
	}
	if req.PartnerType != nil {
		ft.PartnerType = req.PartnerType
	}
	if req.PartnerID != nil {
		ft.PartnerID = req.PartnerID
	}
	if req.PartnerNameCache != nil {
		ft.PartnerNameCache = req.PartnerNameCache
	}
	if req.RefType != nil {
		ft.RefType = req.RefType
	}
	if req.RefID != nil {
		ft.RefID = req.RefID
	}
	if req.Description != nil {
		ft.Description = req.Description
	}
	if req.EvidenceURL != nil {
		ft.EvidenceURL = req.EvidenceURL
	}
	if req.Status != nil {
		ft.Status = *req.Status
	}
	if req.CashTally != nil {
		ft.CashTally = *req.CashTally
	}
	if req.RefAdvanceID != nil {
		ft.RefAdvanceID = req.RefAdvanceID
	}
	if req.TargetBankInfo != nil {
		ft.TargetBankInfo = *req.TargetBankInfo
	}
	if req.BankReferenceID != nil {
		ft.BankReferenceID = req.BankReferenceID
	}
	if req.BookType != nil {
		ft.BookType = *req.BookType
	}
	if req.IsPosted != nil {
		ft.IsPosted = *req.IsPosted
	}
	
	now := time.Now()
	ft.UpdatedAt = &now

	if err := UpdateFinanceTransaction(ft); err != nil {
		return nil, err
	}
	return ft, nil
}

func DeleteFinanceTransactionService(id string) error {
	return DeleteFinanceTransaction(id)
}
