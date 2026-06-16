package audit_logs

import (
	"errors"

	"github.com/namviet/backend-core/internal/platform/supabase"
	"gorm.io/gorm"
)

func GetAllAuditLogs() ([]SystemAuditLog, error) {
	var results []SystemAuditLog
	db := supabase.DB
	err := db.Order("created_at DESC").Limit(100).Find(&results).Error
	return results, err
}

func GetAuditLogByID(id string) (*SystemAuditLog, error) {
	var result SystemAuditLog
	db := supabase.DB
	err := db.First(&result, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("audit log not found")
		}
		return nil, err
	}
	return &result, nil
}
