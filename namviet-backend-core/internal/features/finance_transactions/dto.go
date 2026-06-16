package finance_transactions

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

type CreateFinanceTransactionRequest struct {
	Code              string       `json:"code" binding:"required"`
	TransactionDate   time.Time    `json:"transaction_date"`
	Flow              string       `json:"flow" binding:"required"`
	BusinessType      *string      `json:"business_type"`
	CategoryID        *int64       `json:"category_id"`
	Amount            float64      `json:"amount" binding:"required"`
	FundAccountID     int64        `json:"fund_account_id" binding:"required"`
	PartnerType       *string      `json:"partner_type"`
	PartnerID         *string      `json:"partner_id"`
	PartnerNameCache  *string      `json:"partner_name_cache"`
	RefType           *string      `json:"ref_type"`
	RefID             *string      `json:"ref_id"`
	Description       *string      `json:"description"`
	EvidenceURL       *string      `json:"evidence_url"`
	CashTally         *roles.JSONB `json:"cash_tally"`
	RefAdvanceID      *int64       `json:"ref_advance_id"`
	TargetBankInfo    *roles.JSONB `json:"target_bank_info"`
	BankReferenceID   *string      `json:"bank_reference_id"`
	BookType          *string      `json:"book_type"`
	CreatedBy         *string      `json:"created_by"`
}

type UpdateFinanceTransactionRequest struct {
	TransactionDate   *time.Time   `json:"transaction_date"`
	Flow              *string      `json:"flow"`
	BusinessType      *string      `json:"business_type"`
	CategoryID        *int64       `json:"category_id"`
	Amount            *float64     `json:"amount"`
	FundAccountID     *int64       `json:"fund_account_id"`
	PartnerType       *string      `json:"partner_type"`
	PartnerID         *string      `json:"partner_id"`
	PartnerNameCache  *string      `json:"partner_name_cache"`
	RefType           *string      `json:"ref_type"`
	RefID             *string      `json:"ref_id"`
	Description       *string      `json:"description"`
	EvidenceURL       *string      `json:"evidence_url"`
	Status            *string      `json:"status"`
	CashTally         *roles.JSONB `json:"cash_tally"`
	RefAdvanceID      *int64       `json:"ref_advance_id"`
	TargetBankInfo    *roles.JSONB `json:"target_bank_info"`
	BankReferenceID   *string      `json:"bank_reference_id"`
	BookType          *string      `json:"book_type"`
	IsPosted          *bool        `json:"is_posted"`
}
