package approvals

import (
	"time"

	"github.com/namviet/backend-core/internal/features/roles"
)

// ApprovalRequest represents the approval_requests table
type ApprovalRequest struct {
	ID          string      `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RequestType string      `gorm:"type:text;not null" json:"request_type"`
	RefID       *string     `gorm:"type:text" json:"ref_id,omitempty"`
	RequesterID *string     `gorm:"type:uuid" json:"requester_id,omitempty"`
	Status      string      `gorm:"type:text;default:'pending'" json:"status"`
	CurrentStep int         `gorm:"type:integer;default:1" json:"current_step"`
	Payload     roles.JSONB `gorm:"type:jsonb" json:"payload,omitempty"`
	CreatedAt   *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
	UpdatedAt   *time.Time  `gorm:"type:timestamp with time zone;default:now()" json:"updated_at,omitempty"`
}

// ApprovalStep represents the approval_steps table
type ApprovalStep struct {
	ID           string     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RequestID    *string    `gorm:"type:uuid;index" json:"request_id,omitempty"`
	StepOrder    int        `gorm:"type:integer;not null" json:"step_order"`
	ApproverID   *string    `gorm:"type:uuid" json:"approver_id,omitempty"`
	ApproverRole *string    `gorm:"type:text" json:"approver_role,omitempty"`
	Status       string     `gorm:"type:text;default:'pending'" json:"status"`
	Comments     *string    `gorm:"type:text" json:"comments,omitempty"`
	ActionAt     *time.Time `gorm:"type:timestamp with time zone" json:"action_at,omitempty"`
	CreatedAt    *time.Time `gorm:"type:timestamp with time zone;default:now()" json:"created_at,omitempty"`
}
