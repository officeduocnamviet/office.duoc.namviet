package manufacturers

// CreateManufacturerRequest
type CreateManufacturerRequest struct {
	Name    string  `json:"name" binding:"required"`
	Country *string `json:"country"`
	Status  string  `json:"status"`
}

// UpdateManufacturerRequest
type UpdateManufacturerRequest struct {
	Name    *string `json:"name"`
	Country *string `json:"country"`
	Status  *string `json:"status"`
}
