package orders

import (
	"time"
)

// Order represents the orders table
type Order struct {
	ID            string       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Code          string       `gorm:"type:text;not null" json:"code"`
	CustomerID    *int64       `gorm:"type:bigint" json:"customer_id,omitempty"`
	CreatorID     *string      `gorm:"type:uuid" json:"creator_id,omitempty"`
	Status        string       `gorm:"type:text;default:'PENDING'" json:"status"`
	OrderType     string       `gorm:"type:text;default:'B2C'" json:"order_type"`
	TotalAmount   *float64     `gorm:"type:numeric;default:0" json:"total_amount,omitempty"`
	FinalAmount   *float64     `gorm:"type:numeric;default:0" json:"final_amount,omitempty"`
	PaymentStatus *string      `gorm:"type:text;default:'unpaid'" json:"payment_status,omitempty"`
	Note          *string      `gorm:"type:text" json:"note,omitempty"`
	Items         []OrderItem  `gorm:"foreignKey:OrderID" json:"items,omitempty"`
	CreatedAt     *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt     *time.Time   `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt     *time.Time   `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}

// OrderItem represents the order_items table
type OrderItem struct {
	ID               string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID          string     `gorm:"type:uuid;not null" json:"order_id"`
	ProductID        int64      `gorm:"type:bigint;not null" json:"product_id"`
	Quantity         int        `gorm:"type:integer;not null" json:"quantity"`
	UOM              string     `gorm:"type:text;not null" json:"uom"` // Unit of Measure
	ConversionFactor *int       `gorm:"type:integer" json:"conversion_factor,omitempty"`
	BaseQuantity     *int       `gorm:"type:integer" json:"base_quantity,omitempty"`
	UnitPrice        float64    `gorm:"type:numeric;not null" json:"unit_price"`
	Discount         *float64   `gorm:"type:numeric;default:0" json:"discount,omitempty"`
	IsGift           *bool      `gorm:"type:boolean;default:false" json:"is_gift,omitempty"`
	Note             *string    `gorm:"type:text" json:"note,omitempty"`
	BatchNo          *string    `gorm:"type:text" json:"batch_no,omitempty"`
	ExpiryDate       *string    `gorm:"type:date" json:"expiry_date,omitempty"`
	TotalLine        *float64   `gorm:"type:numeric" json:"total_line,omitempty"`
	QuantityPicked   *int       `gorm:"type:integer;default:0" json:"quantity_picked,omitempty"`
	QuantityReturned *int       `gorm:"type:integer;default:0" json:"quantity_returned,omitempty"`
	CreatedAt        *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	DeletedAt        *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
