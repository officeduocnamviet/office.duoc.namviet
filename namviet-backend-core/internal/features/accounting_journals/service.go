package accounting_journals

func GetAllAccountingJournalsService() ([]AccountingJournal, error) {
	return GetAllAccountingJournals()
}

func GetAccountingJournalByIDService(id string) (*AccountingJournal, error) {
	return GetAccountingJournalByID(id)
}

func CreateAccountingJournalService(req CreateAccountingJournalRequest) (*AccountingJournal, error) {
	aj := &AccountingJournal{
		EntryDate:     req.EntryDate,
		DocType:       req.DocType,
		SourceRefID:   req.SourceRefID,
		Description:   req.Description,
		AccountDebit:  req.AccountDebit,
		AccountCredit: req.AccountCredit,
		Amount:        req.Amount,
		PostedBy:      req.PostedBy,
	}

	if err := CreateAccountingJournal(aj); err != nil {
		return nil, err
	}
	return aj, nil
}

func UpdateAccountingJournalService(id string, req UpdateAccountingJournalRequest) (*AccountingJournal, error) {
	aj, err := GetAccountingJournalByID(id)
	if err != nil {
		return nil, err
	}

	if req.EntryDate != nil {
		aj.EntryDate = *req.EntryDate
	}
	if req.DocType != nil {
		aj.DocType = *req.DocType
	}
	if req.SourceRefID != nil {
		aj.SourceRefID = req.SourceRefID
	}
	if req.Description != nil {
		aj.Description = req.Description
	}
	if req.AccountDebit != nil {
		aj.AccountDebit = *req.AccountDebit
	}
	if req.AccountCredit != nil {
		aj.AccountCredit = *req.AccountCredit
	}
	if req.Amount != nil {
		aj.Amount = *req.Amount
	}
	if req.PostedBy != nil {
		aj.PostedBy = req.PostedBy
	}

	if err := UpdateAccountingJournal(aj); err != nil {
		return nil, err
	}
	return aj, nil
}

func DeleteAccountingJournalService(id string) error {
	return DeleteAccountingJournal(id)
}
