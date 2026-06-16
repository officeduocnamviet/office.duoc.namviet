package orders

type CreateOrderRequest struct {
	Code          string                  `json:"code" binding:"required"`
	CustomerID    *int64                  `json:"customer_id"`
	CreatorID     *string                 `json:"creator_id"`
	OrderType     string                  `json:"order_type"` // B2B or B2C
	Note          *string                 `json:"note"`
	Items         []CreateOrderItemRequest `json:"items" binding:"required"`
}

type CreateOrderItemRequest struct {
	ProductID        int64    `json:"product_id" binding:"required"`
	Quantity         int      `json:"quantity" binding:"required"`
	UOM              string   `json:"uom" binding:"required"`
	ConversionFactor *int     `json:"conversion_factor"`
	UnitPrice        float64  `json:"unit_price" binding:"required"`
	Discount         *float64 `json:"discount"`
	IsGift           *bool    `json:"is_gift"`
	Note             *string  `json:"note"`
	BatchNo          *string  `json:"batch_no"`
	ExpiryDate       *string  `json:"expiry_date"`
}

type UpdateOrderRequest struct {
	Status        *string `json:"status"`
	PaymentStatus *string `json:"payment_status"`
	Note          *string `json:"note"`
}
