package approvals

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// --- Approval Requests API ---

// GetAllApprovalRequestsHandler
// @Summary Get All Approval Requests
// @Description Retrieve a list of all approval requests
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ApprovalRequest
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-requests [get]
func GetAllApprovalRequestsHandler(c *gin.Context) {
	items, err := GetAllApprovalRequestsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetApprovalRequestHandler
// @Summary Get Approval Request by ID
// @Description Retrieve a specific approval request
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {object} ApprovalRequest
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /approval-requests/{id} [get]
func GetApprovalRequestHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetApprovalRequestByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateApprovalRequestHandler
// @Summary Create a new Approval Request
// @Description Create a new approval request
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateApprovalRequestDto true "Request Details"
// @Success 201 {object} ApprovalRequest
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-requests [post]
func CreateApprovalRequestHandler(c *gin.Context) {
	var req CreateApprovalRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateApprovalRequestService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateApprovalRequestHandler
// @Summary Update Approval Request
// @Description Update approval request details
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Param request body UpdateApprovalRequestDto true "Update Details"
// @Success 200 {object} ApprovalRequest
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-requests/{id} [put]
func UpdateApprovalRequestHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateApprovalRequestDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateApprovalRequestService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteApprovalRequestHandler
// @Summary Delete Approval Request
// @Description Delete an approval request
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-requests/{id} [delete]
func DeleteApprovalRequestHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteApprovalRequestService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}


// --- Approval Steps API ---

// GetStepsByRequestIDHandler
// @Summary Get Approval Steps by Request ID
// @Description Retrieve steps for a request
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Request ID"
// @Success 200 {array} ApprovalStep
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-requests/{id}/steps [get]
func GetStepsByRequestIDHandler(c *gin.Context) {
	requestID := c.Param("id")
	if requestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID format"})
		return
	}
	items, err := GetStepsByRequestIDService(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// CreateApprovalStepHandler
// @Summary Create a new Approval Step
// @Description Create a new approval step
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateApprovalStepDto true "Step Details"
// @Success 201 {object} ApprovalStep
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-steps [post]
func CreateApprovalStepHandler(c *gin.Context) {
	var req CreateApprovalStepDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateApprovalStepService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateApprovalStepHandler
// @Summary Update Approval Step
// @Description Update approval step details (e.g. status)
// @Tags Approvals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Step ID"
// @Param request body UpdateApprovalStepDto true "Update Details"
// @Success 200 {object} ApprovalStep
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /approval-steps/{id} [put]
func UpdateApprovalStepHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateApprovalStepDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateApprovalStepService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}
