package warehouses

// CreateWarehouseRequest
type CreateWarehouseRequest struct {
	Name      string  `json:"name" binding:"required"`
	Type      string  `json:"type"`
	Address   *string `json:"address"`
	Manager   *string `json:"manager"`
	Status    string  `json:"status"`
}

// UpdateWarehouseRequest
type UpdateWarehouseRequest struct {
	Name      *string `json:"name"`
	Type      *string `json:"type"`
	Address   *string `json:"address"`
	Manager   *string `json:"manager"`
	Status    *string `json:"status"`
}
