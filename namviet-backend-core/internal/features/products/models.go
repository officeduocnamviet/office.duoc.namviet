package products

import "time"

type Product struct {
	ID                    int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name                  string     `gorm:"type:text;not null" json:"name"`
	SKU                   *string    `gorm:"type:text" json:"sku,omitempty"`
	Barcode               *string    `gorm:"type:text" json:"barcode,omitempty"`
	Description           *string    `gorm:"type:text" json:"description,omitempty"`
	ActiveIngredient      *string    `gorm:"type:text" json:"active_ingredient,omitempty"`
	ImageURL              *string    `gorm:"type:text" json:"image_url,omitempty"`
	Status                string     `gorm:"type:text;default:'active'" json:"status"`
	CategoryID            *int64     `gorm:"type:bigint" json:"category_id,omitempty"`
	ManufacturerID        *int64     `gorm:"type:bigint" json:"manufacturer_id,omitempty"`
	CategoryName          *string    `gorm:"type:text" json:"category_name,omitempty"`
	ManufacturerName      *string    `gorm:"type:text" json:"manufacturer_name,omitempty"`
	DistributorID         *int64     `gorm:"type:bigint" json:"distributor_id,omitempty"`
	InvoicePrice          *float64   `gorm:"type:numeric;default:0" json:"invoice_price,omitempty"`
	ActualCost            float64    `gorm:"type:numeric;default:0;not null" json:"actual_cost"`
	WholesaleUnit         *string    `gorm:"type:text;default:'Hộp'" json:"wholesale_unit,omitempty"`
	RetailUnit            *string    `gorm:"type:text;default:'Vỉ'" json:"retail_unit,omitempty"`
	ConversionFactor      *int       `gorm:"type:integer;default:1" json:"conversion_factor,omitempty"`
	WholesaleMarginValue  *float64   `gorm:"type:numeric;default:0" json:"wholesale_margin_value,omitempty"`
	WholesaleMarginType   *string    `gorm:"type:text;default:'%'" json:"wholesale_margin_type,omitempty"`
	RetailMarginValue     *float64   `gorm:"type:numeric;default:0" json:"retail_margin_value,omitempty"`
	RetailMarginType      *string    `gorm:"type:text;default:'%'" json:"retail_margin_type,omitempty"`
	ItemsPerCarton        *int       `gorm:"type:integer;default:1" json:"items_per_carton,omitempty"`
	CartonWeight          *float64   `gorm:"type:numeric;default:0" json:"carton_weight,omitempty"`
	CartonDimensions      *string    `gorm:"type:text" json:"carton_dimensions,omitempty"`
	PurchasingPolicy      *string    `gorm:"type:text;default:'ALLOW_LOOSE'" json:"purchasing_policy,omitempty"`
	RegistrationNumber    *string    `gorm:"type:text" json:"registration_number,omitempty"`
	CreatedAt             *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt             *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt             *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
