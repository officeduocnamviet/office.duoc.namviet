package chart_of_accounts

type CreateChartOfAccountRequest struct {
	AccountCode  string  `json:"account_code" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	ParentID     *string `json:"parent_id"`
	Type         string  `json:"type" binding:"required"`
	BalanceType  string  `json:"balance_type" binding:"required"`
	AllowPosting *bool   `json:"allow_posting"`
}

type UpdateChartOfAccountRequest struct {
	AccountCode  *string `json:"account_code"`
	Name         *string `json:"name"`
	ParentID     *string `json:"parent_id"`
	Type         *string `json:"type"`
	BalanceType  *string `json:"balance_type"`
	Status       *string `json:"status"`
	AllowPosting *bool   `json:"allow_posting"`
}
