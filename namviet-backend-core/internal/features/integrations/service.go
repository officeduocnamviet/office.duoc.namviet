package integrations

import "time"

// Connections
func GetAllConnectionsService() ([]ThirdPartyConnection, error) {
	return GetAllConnections()
}

func GetConnectionByIDService(id string) (*ThirdPartyConnection, error) {
	return GetConnectionByID(id)
}

func CreateConnectionService(req CreateConnectionRequest) (*ThirdPartyConnection, error) {
	conn := &ThirdPartyConnection{
		PartnerName: req.PartnerName,
		APIKey:      req.APIKey,
		SecretKey:   req.SecretKey,
		WebhookURL:  req.WebhookURL,
		Status:      "active",
	}

	if err := CreateConnection(conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func UpdateConnectionService(id string, req UpdateConnectionRequest) (*ThirdPartyConnection, error) {
	conn, err := GetConnectionByID(id)
	if err != nil {
		return nil, err
	}

	if req.PartnerName != nil {
		conn.PartnerName = *req.PartnerName
	}
	if req.APIKey != nil {
		conn.APIKey = req.APIKey
	}
	if req.SecretKey != nil {
		conn.SecretKey = req.SecretKey
	}
	if req.WebhookURL != nil {
		conn.WebhookURL = req.WebhookURL
	}
	if req.Status != nil {
		conn.Status = *req.Status
	}

	now := time.Now()
	conn.UpdatedAt = &now

	if err := UpdateConnection(conn); err != nil {
		return nil, err
	}
	return conn, nil
}

func DeleteConnectionService(id string) error {
	return DeleteConnection(id)
}

// Webhook Logs
func GetAllWebhookLogsService() ([]WebhookLog, error) {
	return GetAllWebhookLogs()
}

func GetWebhookLogByIDService(id string) (*WebhookLog, error) {
	return GetWebhookLogByID(id)
}

func CreateWebhookLogService(req CreateWebhookLogRequest) (*WebhookLog, error) {
	log := &WebhookLog{
		PartnerID:      req.PartnerID,
		EventType:      req.EventType,
		Payload:        req.Payload,
		ResponseStatus: req.ResponseStatus,
		ResponseBody:   req.ResponseBody,
	}

	if err := CreateWebhookLog(log); err != nil {
		return nil, err
	}
	return log, nil
}
