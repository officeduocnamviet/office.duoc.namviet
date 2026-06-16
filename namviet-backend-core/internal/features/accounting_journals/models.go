package accounting_journals

import (
	"time"
)

// AccountingJournal represents the accounting_journals table
type AccountingJournal struct {
	ID            string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	EntryDate     time.Time  `gorm:"type:date;not null" json:"entry_date"`
	DocType       string     `gorm:"type:text;not null" json:"doc_type"`
	SourceRefID   *string    `gorm:"type:text" json:"source_ref_id,omitempty"`
	Description   *string    `gorm:"type:text" json:"description,omitempty"`
	AccountDebit  string     `gorm:"type:text;not null" json:"account_debit"`
	AccountCredit string     `gorm:"type:text;not null" json:"account_credit"`
	Amount        float64    `gorm:"type:numeric;not null" json:"amount"`
	PostedBy      *string    `gorm:"type:uuid" json:"posted_by,omitempty"`
	CreatedAt     *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
