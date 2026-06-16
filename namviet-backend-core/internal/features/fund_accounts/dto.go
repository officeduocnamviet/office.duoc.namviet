package fund_accounts

type CreateFundAccountRequest struct {
	Name           string       `json:"name" binding:"required"`
	Type           string       `json:"type" binding:"required"`
	Location       *string      `json:"location"`
	AccountNumber  *string      `json:"account_number"`
	BankID         *int64       `json:"bank_id"`
	InitialBalance float64      `json:"initial_balance"`
	Balance        float64      `json:"balance"`
	Currency       *string      `json:"currency"`
	BankInfo       *JSONMap     `json:"bank_info"`
	Description    *string      `json:"description"`
	AccountID      *string      `json:"account_id"`
}

type UpdateFundAccountRequest struct {
	Name          *string      `json:"name"`
	Type          *string      `json:"type"`
	Location      *string      `json:"location"`
	AccountNumber *string      `json:"account_number"`
	BankID        *int64       `json:"bank_id"`
	Balance       *float64     `json:"balance"`
	Currency      *string      `json:"currency"`
	Status        *string      `json:"status"`
	BankInfo      *JSONMap     `json:"bank_info"`
	Description   *string      `json:"description"`
	AccountID     *string      `json:"account_id"`
}
