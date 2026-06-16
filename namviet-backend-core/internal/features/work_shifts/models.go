package work_shifts

import (
	"time"
)

// WorkShift represents the work_shifts table
type WorkShift struct {
	ID        int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	BranchID  int64      `gorm:"type:bigint;not null" json:"branch_id"`
	Name      string     `gorm:"type:text;not null" json:"name"`
	StartTime string     `gorm:"type:time without time zone;not null" json:"start_time"`
	EndTime   string     `gorm:"type:time without time zone;not null" json:"end_time"`
	IsActive  bool       `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}

// ShiftAssignment represents the shift_assignments table
type ShiftAssignment struct {
	ID         int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ShiftID    int64      `gorm:"type:bigint;not null" json:"shift_id"`
	UserID     string     `gorm:"type:uuid;not null" json:"user_id"`
	WorkDate   time.Time  `gorm:"type:date;not null" json:"work_date"`
	Status     string     `gorm:"type:text;default:'scheduled';not null" json:"status"`
	IsOvertime bool       `gorm:"type:boolean;default:false" json:"is_overtime"`
	CreatedAt  *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}

// ShiftHandover represents the shift_handovers table
type ShiftHandover struct {
	ID                   string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AssignmentID         *int64     `gorm:"type:bigint" json:"assignment_id,omitempty"`
	UserID               string     `gorm:"type:uuid;not null" json:"user_id"`
	BranchID             int64      `gorm:"type:bigint;not null" json:"branch_id"`
	SystemCashAmount     float64    `gorm:"type:numeric;default:0;not null" json:"system_cash_amount"`
	SystemCODAmount      float64    `gorm:"type:numeric;default:0;not null" json:"system_cod_amount"`
	ActualCashSubmitted  float64    `gorm:"type:numeric;not null" json:"actual_cash_submitted"`
	Status               string     `gorm:"type:text;default:'pending_finance';not null" json:"status"`
	FinanceTransactionID *int64     `gorm:"type:bigint" json:"finance_transaction_id,omitempty"`
	CreatedAt            *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
