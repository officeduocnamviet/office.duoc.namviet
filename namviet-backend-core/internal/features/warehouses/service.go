package warehouses

import "time"

func GetAllWarehousesService() ([]Warehouse, error) {
	return GetAllWarehouses()
}

func GetWarehouseByIDService(id int64) (*Warehouse, error) {
	return GetWarehouseByID(id)
}

func CreateWarehouseService(req CreateWarehouseRequest) (*Warehouse, error) {
	status := "active"
	if req.Status != "" {
		status = req.Status
	}
	
	whType := "main"
	if req.Type != "" {
		whType = req.Type
	}

	warehouse := &Warehouse{
		Name:      req.Name,
		Type:      whType,
		Address:   req.Address,
		Manager:   req.Manager,
		Status:    status,
	}

	if err := CreateWarehouse(warehouse); err != nil {
		return nil, err
	}
	return warehouse, nil
}

func UpdateWarehouseService(id int64, req UpdateWarehouseRequest) (*Warehouse, error) {
	warehouse, err := GetWarehouseByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		warehouse.Name = *req.Name
	}
	if req.Type != nil {
		warehouse.Type = *req.Type
	}
	if req.Address != nil {
		warehouse.Address = req.Address
	}
	if req.Manager != nil {
		warehouse.Manager = req.Manager
	}
	if req.Status != nil {
		warehouse.Status = *req.Status
	}
	
	now := time.Now()
	warehouse.UpdatedAt = &now

	if err := UpdateWarehouse(warehouse); err != nil {
		return nil, err
	}
	return warehouse, nil
}
