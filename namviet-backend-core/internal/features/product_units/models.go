package product_units

import "time"

type ProductUnit struct {
	ID               int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID        int64      `gorm:"type:bigint;not null" json:"product_id"`
	UnitName         string     `gorm:"type:text;not null" json:"unit_name"`
	ConversionFactor int        `gorm:"type:integer;default:1;not null" json:"conversion_factor"`
	PriceSell        float64    `gorm:"type:numeric;default:0;not null" json:"price_sell"`
	PriceCost        *float64   `gorm:"type:numeric;default:0" json:"price_cost,omitempty"`
	IsBaseUnit       *bool      `gorm:"type:boolean;default:false" json:"is_base_unit,omitempty"`
	CreatedAt        *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt        *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
	DeletedAt        *time.Time `gorm:"type:timestamp with time zone" json:"deleted_at,omitempty"`
}
