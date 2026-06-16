package marketing_campaigns

import "time"

func GetAllMarketingCampaignsService() ([]MarketingCampaign, error) {
	return GetAllMarketingCampaigns()
}

func GetMarketingCampaignByIDService(id string) (*MarketingCampaign, error) {
	return GetMarketingCampaignByID(id)
}

func CreateMarketingCampaignService(req CreateMarketingCampaignRequest) (*MarketingCampaign, error) {
	campaign := &MarketingCampaign{
		Name:          req.Name,
		Description:   req.Description,
		TargetSegment: req.TargetSegment,
		Budget:        req.Budget,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Status:        "draft",
	}

	if err := CreateMarketingCampaign(campaign); err != nil {
		return nil, err
	}
	return campaign, nil
}

func UpdateMarketingCampaignService(id string, req UpdateMarketingCampaignRequest) (*MarketingCampaign, error) {
	campaign, err := GetMarketingCampaignByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		campaign.Name = *req.Name
	}
	if req.Description != nil {
		campaign.Description = req.Description
	}
	if req.TargetSegment != nil {
		campaign.TargetSegment = req.TargetSegment
	}
	if req.Budget != nil {
		campaign.Budget = *req.Budget
	}
	if req.StartDate != nil {
		campaign.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		campaign.EndDate = req.EndDate
	}
	if req.Status != nil {
		campaign.Status = *req.Status
	}

	now := time.Now()
	campaign.UpdatedAt = &now

	if err := UpdateMarketingCampaign(campaign); err != nil {
		return nil, err
	}
	return campaign, nil
}

func DeleteMarketingCampaignService(id string) error {
	return DeleteMarketingCampaign(id)
}
