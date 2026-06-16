package accounting_journals

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAllAccountingJournalsHandler
// @Summary Get All Accounting Journals
// @Description Retrieve a list of all accounting journals
// @Tags Accounting Journals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {array} AccountingJournal
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounting-journals [get]
func GetAllAccountingJournalsHandler(c *gin.Context) {
	ajs, err := GetAllAccountingJournalsService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ajs)
}

// GetAccountingJournalHandler
// @Summary Get Accounting Journal by ID
// @Description Retrieve a specific accounting journal
// @Tags Accounting Journals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Accounting Journal ID"
// @Success 200 {object} AccountingJournal
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /accounting-journals/{id} [get]
func GetAccountingJournalHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	aj, err := GetAccountingJournalByIDService(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aj)
}

// CreateAccountingJournalHandler
// @Summary Create a new Accounting Journal
// @Description Create a new accounting journal
// @Tags Accounting Journals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateAccountingJournalRequest true "Accounting Journal Details"
// @Success 201 {object} AccountingJournal
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounting-journals [post]
func CreateAccountingJournalHandler(c *gin.Context) {
	var req CreateAccountingJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aj, err := CreateAccountingJournalService(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, aj)
}

// UpdateAccountingJournalHandler
// @Summary Update Accounting Journal
// @Description Update accounting journal details
// @Tags Accounting Journals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Accounting Journal ID"
// @Param request body UpdateAccountingJournalRequest true "Update Details"
// @Success 200 {object} AccountingJournal
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounting-journals/{id} [put]
func UpdateAccountingJournalHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	var req UpdateAccountingJournalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	aj, err := UpdateAccountingJournalService(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, aj)
}

// DeleteAccountingJournalHandler
// @Summary Delete Accounting Journal
// @Description Delete an accounting journal
// @Tags Accounting Journals
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Accounting Journal ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /accounting-journals/{id} [delete]
func DeleteAccountingJournalHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}
	
	if err := DeleteAccountingJournalService(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
