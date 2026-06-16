package audit_logs

func GetAllAuditLogsService() ([]SystemAuditLog, error) {
	return GetAllAuditLogs()
}

func GetAuditLogByIDService(id string) (*SystemAuditLog, error) {
	return GetAuditLogByID(id)
}
