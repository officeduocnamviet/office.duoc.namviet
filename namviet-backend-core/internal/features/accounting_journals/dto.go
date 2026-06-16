package accounting_journals

import "time"

type CreateAccountingJournalRequest struct {
	EntryDate     time.Time `json:"entry_date" binding:"required"`
	DocType       string    `json:"doc_type" binding:"required"`
	SourceRefID   *string   `json:"source_ref_id"`
	Description   *string   `json:"description"`
	AccountDebit  string    `json:"account_debit" binding:"required"`
	AccountCredit string    `json:"account_credit" binding:"required"`
	Amount        float64   `json:"amount" binding:"required"`
	PostedBy      *string   `json:"posted_by"`
}

type UpdateAccountingJournalRequest struct {
	EntryDate     *time.Time `json:"entry_date"`
	DocType       *string    `json:"doc_type"`
	SourceRefID   *string    `json:"source_ref_id"`
	Description   *string    `json:"description"`
	AccountDebit  *string    `json:"account_debit"`
	AccountCredit *string    `json:"account_credit"`
	Amount        *float64   `json:"amount"`
	PostedBy      *string    `json:"posted_by"`
}
