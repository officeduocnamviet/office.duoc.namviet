package batches

import "time"

// Batch represents the batches table
type Batch struct {
	ID                int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID         int64      `gorm:"type:bigint;not null" json:"product_id"`
	BatchCode         string     `gorm:"type:text;not null" json:"batch_code"`
	ExpiryDate        string     `gorm:"type:date;not null" json:"expiry_date"`
	ManufacturingDate *string    `gorm:"type:date" json:"manufacturing_date,omitempty"`
	InboundPrice      *float64   `gorm:"type:numeric;default:0" json:"inbound_price,omitempty"`
	CreatedAt         *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt         *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt         *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
