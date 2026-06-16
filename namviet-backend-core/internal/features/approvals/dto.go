package approvals

import "github.com/namviet/backend-core/internal/features/roles"

type CreateApprovalRequestDto struct {
	RequestType string       `json:"request_type" binding:"required"`
	RefID       *string      `json:"ref_id"`
	RequesterID *string      `json:"requester_id"`
	Payload     *roles.JSONB `json:"payload"`
}

type UpdateApprovalRequestDto struct {
	Status      *string      `json:"status"`
	CurrentStep *int         `json:"current_step"`
	Payload     *roles.JSONB `json:"payload"`
}

type CreateApprovalStepDto struct {
	RequestID    *string `json:"request_id"`
	StepOrder    int     `json:"step_order" binding:"required"`
	ApproverID   *string `json:"approver_id"`
	ApproverRole *string `json:"approver_role"`
}

type UpdateApprovalStepDto struct {
	Status     *string `json:"status"`
	Comments   *string `json:"comments"`
}
