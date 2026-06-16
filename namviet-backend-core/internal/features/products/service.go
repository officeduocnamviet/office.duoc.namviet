package products

import "time"

func GetAllProductsService() ([]Product, error) {
	return GetAllProducts()
}

func GetProductByIDService(id int64) (*Product, error) {
	return GetProductByID(id)
}

func CreateProductService(req CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:             req.Name,
		SKU:              req.SKU,
		Barcode:          req.Barcode,
		Description:      req.Description,
		CategoryID:       req.CategoryID,
		ManufacturerID:   req.ManufacturerID,
		CategoryName:     req.CategoryName,
		ManufacturerName: req.ManufacturerName,
		ActualCost:       req.ActualCost,
		WholesaleUnit:    req.WholesaleUnit,
		RetailUnit:       req.RetailUnit,
		ConversionFactor: req.ConversionFactor,
		PriceCost:        req.PriceCost,
		PriceSell:        req.PriceSell,
		Status:           "active",
	}

	if err := CreateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}

func UpdateProductService(id int64, req UpdateProductRequest) (*Product, error) {
	product, err := GetProductByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.SKU != nil {
		product.SKU = req.SKU
	}
	if req.Barcode != nil {
		product.Barcode = req.Barcode
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.ManufacturerID != nil {
		product.ManufacturerID = req.ManufacturerID
	}
	if req.CategoryName != nil {
		product.CategoryName = req.CategoryName
	}
	if req.ManufacturerName != nil {
		product.ManufacturerName = req.ManufacturerName
	}
	if req.ActualCost != nil {
		product.ActualCost = *req.ActualCost
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	
	now := time.Now()
	product.UpdatedAt = &now

	if err := UpdateProduct(product); err != nil {
		return nil, err
	}
	return product, nil
}
