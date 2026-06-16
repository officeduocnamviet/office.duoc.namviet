package promotions

func GetAllPromotionsService() ([]Promotion, error) {
	return GetAllPromotions()
}

func GetPromotionByIDService(id string) (*Promotion, error) {
	return GetPromotionByID(id)
}

func CreatePromotionService(req CreatePromotionRequest) (*Promotion, error) {
	status := "active"
	promotion := &Promotion{
		Code:      req.Code,
		Name:      req.Name,
		Rules:     req.Rules,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Status:    &status,
	}

	if err := CreatePromotion(promotion); err != nil {
		return nil, err
	}
	return promotion, nil
}

func UpdatePromotionService(id string, req UpdatePromotionRequest) (*Promotion, error) {
	promotion, err := GetPromotionByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		promotion.Name = *req.Name
	}
	if req.Rules != nil {
		promotion.Rules = *req.Rules
	}
	if req.StartDate != nil {
		promotion.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		promotion.EndDate = *req.EndDate
	}
	if req.Status != nil {
		promotion.Status = req.Status
	}

	if err := UpdatePromotion(promotion); err != nil {
		return nil, err
	}
	return promotion, nil
}

func DeletePromotionService(id string) error {
	return DeletePromotion(id)
}
