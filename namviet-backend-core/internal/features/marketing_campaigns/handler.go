package marketing_campaigns

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllMarketingCampaignsHandler
// @Summary Get All Marketing Campaigns
// @Description Retrieve a list of all marketing campaigns
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} MarketingCampaign
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /marketing-campaigns [get]
func GetAllMarketingCampaignsHandler(c *gin.Context) {
	items, err := GetAllMarketingCampaignsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetMarketingCampaignHandler
// @Summary Get Marketing Campaign by ID
// @Description Retrieve a specific marketing campaign
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Campaign ID"
// @Success 200 {object} MarketingCampaign
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /marketing-campaigns/{id} [get]
func GetMarketingCampaignHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetMarketingCampaignByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateMarketingCampaignHandler
// @Summary Create a new Marketing Campaign
// @Description Create a new marketing campaign
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateMarketingCampaignRequest true "Campaign Details"
// @Success 201 {object} MarketingCampaign
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /marketing-campaigns [post]
func CreateMarketingCampaignHandler(c *gin.Context) {
	var req CreateMarketingCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateMarketingCampaignService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateMarketingCampaignHandler
// @Summary Update Marketing Campaign
// @Description Update marketing campaign details
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Campaign ID"
// @Param request body UpdateMarketingCampaignRequest true "Update Details"
// @Success 200 {object} MarketingCampaign
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /marketing-campaigns/{id} [put]
func UpdateMarketingCampaignHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateMarketingCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateMarketingCampaignService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteMarketingCampaignHandler
// @Summary Delete Marketing Campaign
// @Description Delete a marketing campaign
// @Tags CRM & Marketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Campaign ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /marketing-campaigns/{id} [delete]
func DeleteMarketingCampaignHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteMarketingCampaignService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
