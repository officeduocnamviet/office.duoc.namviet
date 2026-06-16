package fund_accounts

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

// FundAccount represents the fund_accounts table
type FundAccount struct {
	ID             int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string      `gorm:"type:text;not null" json:"name"`
	Type           string      `gorm:"type:text;not null" json:"type"`
	Location       *string     `gorm:"type:text" json:"location,omitempty"`
	AccountNumber  *string     `gorm:"type:text" json:"account_number,omitempty"`
	BankID         *int64      `gorm:"type:bigint" json:"bank_id,omitempty"`
	InitialBalance float64     `gorm:"type:numeric;default:0" json:"initial_balance"`
	Balance        float64     `gorm:"type:numeric;default:0" json:"balance"`
	Currency       string      `gorm:"type:text;default:'VND'" json:"currency"`
	Status         string      `gorm:"type:text;default:'active'" json:"status"`
	BankInfo       JSONMap     `gorm:"type:jsonb" json:"bank_info,omitempty"`
	Description    *string     `gorm:"type:text" json:"description,omitempty"`
	AccountID      *string     `gorm:"type:text" json:"account_id,omitempty"`
	CreatedAt      *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt      *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt      *time.Time  `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
