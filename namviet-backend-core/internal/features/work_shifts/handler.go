package work_shifts

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// --- Work Shifts API ---

// GetAllWorkShiftsHandler
// @Summary Get All Work Shifts
// @Description Retrieve a list of all work shifts
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} WorkShift
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /work-shifts [get]
func GetAllWorkShiftsHandler(c *gin.Context) {
	items, err := GetAllWorkShiftsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetWorkShiftHandler
// @Summary Get Work Shift by ID
// @Description Retrieve a specific work shift
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Shift ID"
// @Success 200 {object} WorkShift
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /work-shifts/{id} [get]
func GetWorkShiftHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetWorkShiftByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateWorkShiftHandler
// @Summary Create a new Work Shift
// @Description Create a new work shift
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateWorkShiftRequest true "Shift Details"
// @Success 201 {object} WorkShift
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /work-shifts [post]
func CreateWorkShiftHandler(c *gin.Context) {
	var req CreateWorkShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateWorkShiftService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateWorkShiftHandler
// @Summary Update Work Shift
// @Description Update work shift details
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Shift ID"
// @Param request body UpdateWorkShiftRequest true "Update Details"
// @Success 200 {object} WorkShift
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /work-shifts/{id} [put]
func UpdateWorkShiftHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateWorkShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateWorkShiftService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteWorkShiftHandler
// @Summary Delete Work Shift
// @Description Delete a work shift
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Shift ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /work-shifts/{id} [delete]
func DeleteWorkShiftHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteWorkShiftService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Shift Assignments API ---

// GetAllShiftAssignmentsHandler
// @Summary Get All Shift Assignments
// @Description Retrieve a list of all shift assignments
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ShiftAssignment
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-assignments [get]
func GetAllShiftAssignmentsHandler(c *gin.Context) {
	items, err := GetAllShiftAssignmentsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetShiftAssignmentHandler
// @Summary Get Shift Assignment by ID
// @Description Retrieve a specific shift assignment
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Assignment ID"
// @Success 200 {object} ShiftAssignment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /shift-assignments/{id} [get]
func GetShiftAssignmentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetShiftAssignmentByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateShiftAssignmentHandler
// @Summary Create a new Shift Assignment
// @Description Create a new shift assignment
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateShiftAssignmentRequest true "Assignment Details"
// @Success 201 {object} ShiftAssignment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-assignments [post]
func CreateShiftAssignmentHandler(c *gin.Context) {
	var req CreateShiftAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateShiftAssignmentService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateShiftAssignmentHandler
// @Summary Update Shift Assignment
// @Description Update shift assignment details
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Assignment ID"
// @Param request body UpdateShiftAssignmentRequest true "Update Details"
// @Success 200 {object} ShiftAssignment
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-assignments/{id} [put]
func UpdateShiftAssignmentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateShiftAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateShiftAssignmentService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteShiftAssignmentHandler
// @Summary Delete Shift Assignment
// @Description Delete a shift assignment
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path integer true "Assignment ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-assignments/{id} [delete]
func DeleteShiftAssignmentHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteShiftAssignmentService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// --- Shift Handovers API ---

// GetAllShiftHandoversHandler
// @Summary Get All Shift Handovers
// @Description Retrieve a list of all shift handovers
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} ShiftHandover
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-handovers [get]
func GetAllShiftHandoversHandler(c *gin.Context) {
	items, err := GetAllShiftHandoversService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetShiftHandoverHandler
// @Summary Get Shift Handover by ID
// @Description Retrieve a specific shift handover
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Handover ID"
// @Success 200 {object} ShiftHandover
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /shift-handovers/{id} [get]
func GetShiftHandoverHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	item, err := GetShiftHandoverByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateShiftHandoverHandler
// @Summary Create a new Shift Handover
// @Description Create a new shift handover
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateShiftHandoverRequest true "Handover Details"
// @Success 201 {object} ShiftHandover
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-handovers [post]
func CreateShiftHandoverHandler(c *gin.Context) {
	var req CreateShiftHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := CreateShiftHandoverService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateShiftHandoverHandler
// @Summary Update Shift Handover
// @Description Update shift handover details
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Handover ID"
// @Param request body UpdateShiftHandoverRequest true "Update Details"
// @Success 200 {object} ShiftHandover
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-handovers/{id} [put]
func UpdateShiftHandoverHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateShiftHandoverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := UpdateShiftHandoverService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteShiftHandoverHandler
// @Summary Delete Shift Handover
// @Description Delete a shift handover
// @Tags Advanced HR
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Handover ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /shift-handovers/{id} [delete]
func DeleteShiftHandoverHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteShiftHandoverService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
