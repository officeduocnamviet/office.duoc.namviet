package approvals

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

// Approval Requests
func GetAllApprovalRequests() ([]ApprovalRequest, error) {
	var results []ApprovalRequest
	db := supabase.DB
	err := db.Find(&results).Error
	return results, err
}

func GetApprovalRequestByID(id string) (*ApprovalRequest, error) {
	var result ApprovalRequest
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("approval request not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateApprovalRequest(data *ApprovalRequest) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateApprovalRequest(data *ApprovalRequest) error {
	db := supabase.DB
	return db.Save(data).Error
}

func DeleteApprovalRequest(id string) error {
	db := supabase.DB
	return db.Where("id = ?", id).Delete(&ApprovalRequest{}).Error
}

// Approval Steps
func GetStepsByRequestID(requestID string) ([]ApprovalStep, error) {
	var results []ApprovalStep
	db := supabase.DB
	err := db.Where("request_id = ?", requestID).Order("step_order ASC").Find(&results).Error
	return results, err
}

func GetApprovalStepByID(id string) (*ApprovalStep, error) {
	var result ApprovalStep
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("approval step not found")
		}
		return nil, err
	}
	return &result, nil
}

func CreateApprovalStep(data *ApprovalStep) error {
	db := supabase.DB
	return db.Create(data).Error
}

func UpdateApprovalStep(data *ApprovalStep) error {
	db := supabase.DB
	return db.Save(data).Error
}
