package approvals

import "time"

// Approval Requests
func GetAllApprovalRequestsService() ([]ApprovalRequest, error) {
	return GetAllApprovalRequests()
}

func GetApprovalRequestByIDService(id string) (*ApprovalRequest, error) {
	return GetApprovalRequestByID(id)
}

func CreateApprovalRequestService(req CreateApprovalRequestDto) (*ApprovalRequest, error) {
	ar := &ApprovalRequest{
		RequestType: req.RequestType,
		RefID:       req.RefID,
		RequesterID: req.RequesterID,
		Status:      "pending",
		CurrentStep: 1,
	}

	if req.Payload != nil {
		ar.Payload = *req.Payload
	}

	if err := CreateApprovalRequest(ar); err != nil {
		return nil, err
	}
	return ar, nil
}

func UpdateApprovalRequestService(id string, req UpdateApprovalRequestDto) (*ApprovalRequest, error) {
	ar, err := GetApprovalRequestByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		ar.Status = *req.Status
	}
	if req.CurrentStep != nil {
		ar.CurrentStep = *req.CurrentStep
	}
	if req.Payload != nil {
		ar.Payload = *req.Payload
	}

	now := time.Now()
	ar.UpdatedAt = &now

	if err := UpdateApprovalRequest(ar); err != nil {
		return nil, err
	}
	return ar, nil
}

func DeleteApprovalRequestService(id string) error {
	return DeleteApprovalRequest(id)
}

// Approval Steps
func GetStepsByRequestIDService(requestID string) ([]ApprovalStep, error) {
	return GetStepsByRequestID(requestID)
}

func CreateApprovalStepService(req CreateApprovalStepDto) (*ApprovalStep, error) {
	step := &ApprovalStep{
		RequestID:    req.RequestID,
		StepOrder:    req.StepOrder,
		ApproverID:   req.ApproverID,
		ApproverRole: req.ApproverRole,
		Status:       "pending",
	}

	if err := CreateApprovalStep(step); err != nil {
		return nil, err
	}
	return step, nil
}

func UpdateApprovalStepService(id string, req UpdateApprovalStepDto) (*ApprovalStep, error) {
	step, err := GetApprovalStepByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		step.Status = *req.Status
		now := time.Now()
		step.ActionAt = &now
	}
	if req.Comments != nil {
		step.Comments = req.Comments
	}

	if err := UpdateApprovalStep(step); err != nil {
		return nil, err
	}
	return step, nil
}
