package attendance_logs

func GetAllAttendanceLogsService() ([]AttendanceLog, error) {
	return GetAllAttendanceLogs()
}

func GetAttendanceLogByIDService(id string) (*AttendanceLog, error) {
	return GetAttendanceLogByID(id)
}

func CreateAttendanceLogService(req CreateAttendanceLogRequest) (*AttendanceLog, error) {
	log := &AttendanceLog{
		UserID:      req.UserID,
		BranchID:    req.BranchID,
		CheckInTime: req.CheckInTime,
		CheckInIP:   req.CheckInIP,
		CheckInLat:  req.CheckInLat,
		CheckInLng:  req.CheckInLng,
		Status:      "present",
	}

	if err := CreateAttendanceLog(log); err != nil {
		return nil, err
	}
	return log, nil
}

func UpdateAttendanceLogService(id string, req UpdateAttendanceLogRequest) (*AttendanceLog, error) {
	log, err := GetAttendanceLogByID(id)
	if err != nil {
		return nil, err
	}

	if req.CheckOutTime != nil {
		log.CheckOutTime = req.CheckOutTime
	}
	if req.CheckOutIP != nil {
		log.CheckOutIP = req.CheckOutIP
	}
	if req.CheckOutLat != nil {
		log.CheckOutLat = req.CheckOutLat
	}
	if req.CheckOutLng != nil {
		log.CheckOutLng = req.CheckOutLng
	}
	if req.Status != nil {
		log.Status = *req.Status
	}

	if err := UpdateAttendanceLog(log); err != nil {
		return nil, err
	}
	return log, nil
}

func DeleteAttendanceLogService(id string) error {
	return DeleteAttendanceLog(id)
}
