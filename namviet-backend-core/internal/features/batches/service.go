package batches

import "time"

func GetBatchesByProductIDService(productID int64) ([]Batch, error) {
	return GetBatchesByProductID(productID)
}

func GetBatchByIDService(id int64) (*Batch, error) {
	return GetBatchByID(id)
}

func CreateBatchService(req CreateBatchRequest) (*Batch, error) {
	batch := &Batch{
		ProductID:         req.ProductID,
		BatchCode:         req.BatchCode,
		ExpiryDate:        req.ExpiryDate,
		ManufacturingDate: req.ManufacturingDate,
		InboundPrice:      req.InboundPrice,
	}

	if err := CreateBatch(batch); err != nil {
		return nil, err
	}
	return batch, nil
}

func UpdateBatchService(id int64, req UpdateBatchRequest) (*Batch, error) {
	batch, err := GetBatchByID(id)
	if err != nil {
		return nil, err
	}

	if req.BatchCode != nil {
		batch.BatchCode = *req.BatchCode
	}
	if req.ExpiryDate != nil {
		batch.ExpiryDate = *req.ExpiryDate
	}
	if req.ManufacturingDate != nil {
		batch.ManufacturingDate = req.ManufacturingDate
	}
	if req.InboundPrice != nil {
		batch.InboundPrice = req.InboundPrice
	}
	
	now := time.Now()
	batch.UpdatedAt = &now

	if err := UpdateBatch(batch); err != nil {
		return nil, err
	}
	return batch, nil
}
