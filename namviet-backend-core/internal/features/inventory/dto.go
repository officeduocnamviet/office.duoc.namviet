package inventory

type CreateTransactionRequest struct {
	WarehouseID   int64   `json:"warehouse_id" binding:"required"`
	ProductID     int64   `json:"product_id" binding:"required"`
	BatchID       *int64  `json:"batch_id"`
	Type          string  `json:"type" binding:"required"` // IN, OUT
	ActionGroup *string  `json:"action_group"`
	Quantity    int      `json:"quantity" binding:"required"`
	UnitPrice   *float64 `json:"unit_price"`
	RefID       *string  `json:"ref_id"`
	Description *string  `json:"description"`
	PartnerID   *int64   `json:"partner_id"`
}
