package orders

import "time"

func GetAllOrdersService() ([]Order, error) {
	return GetAllOrders()
}

func GetOrderByIDService(id string) (*Order, error) {
	return GetOrderByID(id)
}

func CreateOrderService(req CreateOrderRequest) (*Order, error) {
	orderType := "B2C"
	if req.OrderType != "" {
		orderType = req.OrderType
	}

	var totalAmount float64
	var finalAmount float64

	items := make([]OrderItem, len(req.Items))
	for i, itemReq := range req.Items {
		
		// Simple basic calculation
		lineTotal := float64(itemReq.Quantity) * itemReq.UnitPrice
		if itemReq.Discount != nil {
			lineTotal -= *itemReq.Discount
		}
		if itemReq.IsGift != nil && *itemReq.IsGift {
			lineTotal = 0
		}
		
		totalAmount += float64(itemReq.Quantity) * itemReq.UnitPrice
		finalAmount += lineTotal
		
		baseQty := itemReq.Quantity
		if itemReq.ConversionFactor != nil {
			baseQty = itemReq.Quantity * (*itemReq.ConversionFactor)
		}

		items[i] = OrderItem{
			ProductID:        itemReq.ProductID,
			Quantity:         itemReq.Quantity,
			UOM:              itemReq.UOM,
			ConversionFactor: itemReq.ConversionFactor,
			BaseQuantity:     &baseQty,
			UnitPrice:        itemReq.UnitPrice,
			Discount:         itemReq.Discount,
			IsGift:           itemReq.IsGift,
			Note:             itemReq.Note,
			BatchNo:          itemReq.BatchNo,
			ExpiryDate:       itemReq.ExpiryDate,
			TotalLine:        &lineTotal,
		}
	}
	
	unpaid := "unpaid"

	order := &Order{
		Code:          req.Code,
		CustomerID:    req.CustomerID,
		CreatorID:     req.CreatorID,
		OrderType:     orderType,
		Note:          req.Note,
		Status:        "PENDING",
		PaymentStatus: &unpaid,
		TotalAmount:   &totalAmount,
		FinalAmount:   &finalAmount,
		Items:         items,
	}

	if err := CreateOrderWithItems(order); err != nil {
		return nil, err
	}
	return order, nil
}

func UpdateOrderService(id string, req UpdateOrderRequest) (*Order, error) {
	order, err := GetOrderByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		order.Status = *req.Status
	}
	if req.PaymentStatus != nil {
		order.PaymentStatus = req.PaymentStatus
	}
	if req.Note != nil {
		order.Note = req.Note
	}
	
	now := time.Now()
	order.UpdatedAt = &now

	if err := UpdateOrder(order); err != nil {
		return nil, err
	}
	return order, nil
}

func DeleteOrderService(id string) error {
	return DeleteOrder(id)
}
