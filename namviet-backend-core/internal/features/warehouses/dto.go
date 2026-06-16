package warehouses

// CreateWarehouseRequest
type CreateWarehouseRequest struct {
	Key       string   `json:"key" binding:"required"`
	Name      string   `json:"name" binding:"required"`
	Unit      string   `json:"unit"`
	Address   *string  `json:"address"`
	Type      string   `json:"type"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Code      *string  `json:"code"`
	Manager   *string  `json:"manager"`
	Phone     *string  `json:"phone"`
}

// UpdateWarehouseRequest
type UpdateWarehouseRequest struct {
	Name      *string  `json:"name"`
	Unit      *string  `json:"unit"`
	Address   *string  `json:"address"`
	Type      *string  `json:"type"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
	Code      *string  `json:"code"`
	Manager   *string  `json:"manager"`
	Phone     *string  `json:"phone"`
}
