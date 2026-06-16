package product_units

type CreateProductUnitRequest struct {
	UnitName         string   `json:"unit_name" binding:"required"`
	ConversionFactor int      `json:"conversion_factor" binding:"required"`
	PriceSell        float64  `json:"price_sell"`
	PriceCost        *float64 `json:"price_cost"`
	IsBaseUnit       *bool    `json:"is_base_unit"`
}

type UpdateProductUnitRequest struct {
	UnitName         *string  `json:"unit_name"`
	ConversionFactor *int     `json:"conversion_factor"`
	PriceSell        *float64 `json:"price_sell"`
	PriceCost        *float64 `json:"price_cost"`
	IsBaseUnit       *bool    `json:"is_base_unit"`
}
