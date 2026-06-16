package inventory

type CreateTransactionRequest struct {
	WarehouseID   int64   `json:"warehouse_id" binding:"required"`
	ProductID     int64   `json:"product_id" binding:"required"`
	BatchID       *int64  `json:"batch_id"`
	Type          string  `json:"type" binding:"required"` // IN, OUT
	Quantity      int     `json:"quantity" binding:"required"`
	ReferenceID   *string `json:"reference_id"`
	ReferenceType *string `json:"reference_type"`
}
