package time_attendance

import "time"

func GetAllTimeAttendancesService() ([]TimeAttendance, error) {
	return GetAllTimeAttendances()
}

func GetTimeAttendanceByIDService(id string) (*TimeAttendance, error) {
	return GetTimeAttendanceByID(id)
}

func CreateTimeAttendanceService(req CreateTimeAttendanceRequest) (*TimeAttendance, error) {
	status := "present"
	if req.Status != nil {
		status = *req.Status
	}

	shiftType := "morning"
	if req.ShiftType != nil {
		shiftType = *req.ShiftType
	}

	ta := &TimeAttendance{
		EmployeeID:    req.EmployeeID,
		Date:          req.Date,
		CheckIn:       req.CheckIn,
		CheckOut:      req.CheckOut,
		Status:        status,
		ShiftType:     shiftType,
		OvertimeHours: req.OvertimeHours,
		Note:          req.Note,
	}

	if req.Location != nil {
		ta.Location = *req.Location
	}
	if req.DeviceInfo != nil {
		ta.DeviceInfo = *req.DeviceInfo
	}

	if err := CreateTimeAttendance(ta); err != nil {
		return nil, err
	}
	return ta, nil
}

func UpdateTimeAttendanceService(id string, req UpdateTimeAttendanceRequest) (*TimeAttendance, error) {
	ta, err := GetTimeAttendanceByID(id)
	if err != nil {
		return nil, err
	}

	if req.CheckIn != nil {
		ta.CheckIn = req.CheckIn
	}
	if req.CheckOut != nil {
		ta.CheckOut = req.CheckOut
	}
	if req.Status != nil {
		ta.Status = *req.Status
	}
	if req.ShiftType != nil {
		ta.ShiftType = *req.ShiftType
	}
	if req.OvertimeHours != nil {
		ta.OvertimeHours = req.OvertimeHours
	}
	if req.Location != nil {
		ta.Location = *req.Location
	}
	if req.DeviceInfo != nil {
		ta.DeviceInfo = *req.DeviceInfo
	}
	if req.Note != nil {
		ta.Note = req.Note
	}
	
	now := time.Now()
	ta.UpdatedAt = &now

	if err := UpdateTimeAttendance(ta); err != nil {
		return nil, err
	}
	return ta, nil
}

func DeleteTimeAttendanceService(id string) error {
	return DeleteTimeAttendance(id)
}
