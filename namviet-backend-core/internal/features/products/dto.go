package products

// CreateProductRequest
type CreateProductRequest struct {
	Name             string  `json:"name" binding:"required"`
	SKU              *string `json:"sku"`
	Barcode          *string `json:"barcode"`
	Description      *string `json:"description"`
	CategoryID       *int64  `json:"category_id"`
	ManufacturerID   *int64  `json:"manufacturer_id"`
	CategoryName     *string `json:"category_name"`
	ManufacturerName *string `json:"manufacturer_name"`
	ActualCost       float64 `json:"actual_cost"`
	WholesaleUnit    *string `json:"wholesale_unit"`
	RetailUnit       *string `json:"retail_unit"`
	ConversionFactor *int    `json:"conversion_factor"`
	PriceCost        *float64 `json:"price_cost"`
	PriceSell        *float64 `json:"price_sell"`
}

// UpdateProductRequest
type UpdateProductRequest struct {
	Name             *string  `json:"name"`
	SKU              *string  `json:"sku"`
	Barcode          *string  `json:"barcode"`
	Description      *string  `json:"description"`
	CategoryID       *int64   `json:"category_id"`
	ManufacturerID   *int64   `json:"manufacturer_id"`
	CategoryName     *string  `json:"category_name"`
	ManufacturerName *string  `json:"manufacturer_name"`
	ActualCost       *float64 `json:"actual_cost"`
	Status           *string  `json:"status"`
}
