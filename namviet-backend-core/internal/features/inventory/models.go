package inventory

import "time"

// InventoryBatch represents the inventory_batches table
type InventoryBatch struct {
	ID          int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   int64      `gorm:"type:bigint;not null" json:"product_id"`
	BatchID     int64      `gorm:"type:bigint;not null" json:"batch_id"`
	WarehouseID int64      `gorm:"type:bigint;not null" json:"warehouse_id"`
	Quantity    int        `gorm:"type:integer;default:0;not null" json:"quantity"`
	CreatedAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt   *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}

// InventoryTransaction represents the inventory_transactions table
type InventoryTransaction struct {
	ID            string     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	WarehouseID   int64      `gorm:"type:bigint;not null" json:"warehouse_id"`
	ProductID     int64      `gorm:"type:bigint;not null" json:"product_id"`
	BatchID       *int64     `gorm:"type:bigint" json:"batch_id,omitempty"`
	Type          string     `gorm:"type:text;not null" json:"type"` // e.g. IN, OUT
	ActionGroup   *string    `gorm:"type:text" json:"action_group,omitempty"`
	Quantity      int        `gorm:"type:integer;not null" json:"quantity"`
	UnitPrice     *float64   `gorm:"type:numeric" json:"unit_price,omitempty"`
	RefID         *string    `gorm:"type:text" json:"ref_id,omitempty"`
	Description   *string    `gorm:"type:text" json:"description,omitempty"`
	PartnerID     *int64     `gorm:"type:bigint" json:"partner_id,omitempty"`
	CreatedBy     *string    `gorm:"type:uuid" json:"created_by,omitempty"`
	CreatedAt     *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
