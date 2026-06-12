package models

import (
	"time"
)

type Warehouse struct {
	ID        int64      `json:"id"`
	CompanyID *string    `json:"company_id"` // UUID Nullable
	Key       string     `json:"key"`
	Code      *string    `json:"code"`
	Name      string     `json:"name"`
	Type      string     `json:"type"`
	Address   *string    `json:"address"`
	Manager   *string    `json:"manager"`
	Phone     *string    `json:"phone"`
	Latitude  *float64   `json:"latitude"`
	Longitude *float64   `json:"longitude"`
	Status     string     `json:"status"`
	OutletType *string    `json:"outlet_type"`
	CreatedAt  *time.Time `json:"created_at"`
}
