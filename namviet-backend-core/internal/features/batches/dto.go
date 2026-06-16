package batches

// CreateBatchRequest
type CreateBatchRequest struct {
	ProductID         int64    `json:"product_id" binding:"required"`
	BatchCode         string   `json:"batch_code" binding:"required"`
	ExpiryDate        string   `json:"expiry_date" binding:"required"`
	ManufacturingDate *string  `json:"manufacturing_date"`
	InboundPrice      *float64 `json:"inbound_price"`
}

// UpdateBatchRequest
type UpdateBatchRequest struct {
	BatchCode         *string  `json:"batch_code"`
	ExpiryDate        *string  `json:"expiry_date"`
	ManufacturingDate *string  `json:"manufacturing_date"`
	InboundPrice      *float64 `json:"inbound_price"`
}
