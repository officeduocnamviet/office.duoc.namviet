package product_units

import "time"

func CreateProductUnitService(productID int64, req CreateProductUnitRequest) (*ProductUnit, error) {
	unit := &ProductUnit{
		ProductID:        productID,
		UnitName:         req.UnitName,
		ConversionFactor: req.ConversionFactor,
		PriceSell:        req.PriceSell,
		PriceCost:        req.PriceCost,
		IsBaseUnit:       req.IsBaseUnit,
	}

	if err := CreateProductUnit(unit); err != nil {
		return nil, err
	}
	return unit, nil
}

func UpdateProductUnitService(id int64, req UpdateProductUnitRequest) (*ProductUnit, error) {
	unit, err := GetProductUnitByID(id)
	if err != nil {
		return nil, err
	}

	if req.UnitName != nil {
		unit.UnitName = *req.UnitName
	}
	if req.ConversionFactor != nil {
		unit.ConversionFactor = *req.ConversionFactor
	}
	if req.PriceSell != nil {
		unit.PriceSell = *req.PriceSell
	}
	if req.PriceCost != nil {
		unit.PriceCost = req.PriceCost
	}
	if req.IsBaseUnit != nil {
		unit.IsBaseUnit = req.IsBaseUnit
	}
	
	now := time.Now()
	unit.UpdatedAt = &now

	if err := UpdateProductUnit(unit); err != nil {
		return nil, err
	}
	return unit, nil
}

func DeleteProductUnitService(id int64) error {
	return DeleteProductUnit(id)
}
