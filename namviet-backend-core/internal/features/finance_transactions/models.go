package finance_transactions

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// JSONMap map for GORM
type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return "{}", nil
	}
	b, err := json.Marshal(j)
	return string(b), err
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return errors.New("type assertion to []byte or string failed")
		}
		bytes = []byte(str)
	}
	return json.Unmarshal(bytes, &j)
}

// FinanceTransaction represents the finance_transactions table
type FinanceTransaction struct {
	ID                int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Code              string      `gorm:"type:text;not null" json:"code"`
	TransactionDate   time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"transaction_date"`
	Flow              string      `gorm:"type:text;not null" json:"flow"`
	BusinessType      string      `gorm:"type:text;default:'other'" json:"business_type"`
	CategoryID        *int64      `gorm:"type:bigint" json:"category_id,omitempty"`
	Amount            float64     `gorm:"type:numeric;not null" json:"amount"`
	FundAccountID     int64       `gorm:"type:bigint;not null" json:"fund_account_id"`
	PartnerType       *string     `gorm:"type:text" json:"partner_type,omitempty"`
	PartnerID         *string     `gorm:"type:text" json:"partner_id,omitempty"`
	PartnerNameCache  *string     `gorm:"type:text" json:"partner_name_cache,omitempty"`
	RefType           *string     `gorm:"type:text" json:"ref_type,omitempty"`
	RefID             *string     `gorm:"type:text" json:"ref_id,omitempty"`
	Description       *string     `gorm:"type:text" json:"description,omitempty"`
	EvidenceURL       *string     `gorm:"type:text" json:"evidence_url,omitempty"`
	Status            string      `gorm:"type:text;default:'pending'" json:"status"`
	CashTally         JSONMap     `gorm:"type:jsonb" json:"cash_tally,omitempty"`
	RefAdvanceID      *int64      `gorm:"type:bigint" json:"ref_advance_id,omitempty"`
	TargetBankInfo    JSONMap     `gorm:"type:jsonb" json:"target_bank_info,omitempty"`
	BankReferenceID   *string     `gorm:"type:text" json:"bank_reference_id,omitempty"`
	BookType          string      `gorm:"type:text;default:'BOTH'" json:"book_type"`
	IsPosted          bool        `gorm:"type:boolean;default:false" json:"is_posted"`
	CreatedBy         *string     `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt         *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt         *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt         *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
