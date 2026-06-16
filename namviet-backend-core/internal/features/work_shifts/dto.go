package work_shifts

import "time"

// Work Shift DTOs
type CreateWorkShiftRequest struct {
	BranchID  int64  `json:"branch_id" binding:"required"`
	Name      string `json:"name" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	IsActive  *bool  `json:"is_active"`
}

type UpdateWorkShiftRequest struct {
	BranchID  *int64  `json:"branch_id"`
	Name      *string `json:"name"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	IsActive  *bool   `json:"is_active"`
}

// Shift Assignment DTOs
type CreateShiftAssignmentRequest struct {
	ShiftID    int64     `json:"shift_id" binding:"required"`
	UserID     string    `json:"user_id" binding:"required"`
	WorkDate   time.Time `json:"work_date" binding:"required"`
	Status     *string   `json:"status"`
	IsOvertime *bool     `json:"is_overtime"`
}

type UpdateShiftAssignmentRequest struct {
	Status     *string `json:"status"`
	IsOvertime *bool   `json:"is_overtime"`
}

// Shift Handover DTOs
type CreateShiftHandoverRequest struct {
	AssignmentID        *int64  `json:"assignment_id"`
	UserID              string  `json:"user_id" binding:"required"`
	BranchID            int64   `json:"branch_id" binding:"required"`
	SystemCashAmount    float64 `json:"system_cash_amount"`
	SystemCODAmount     float64 `json:"system_cod_amount"`
	ActualCashSubmitted float64 `json:"actual_cash_submitted" binding:"required"`
}

type UpdateShiftHandoverRequest struct {
	Status               *string `json:"status"`
	FinanceTransactionID *int64  `json:"finance_transaction_id"`
}
