package work_shifts

// Work Shifts
func GetAllWorkShiftsService() ([]WorkShift, error) {
	return GetAllWorkShifts()
}

func GetWorkShiftByIDService(id int64) (*WorkShift, error) {
	return GetWorkShiftByID(id)
}

func CreateWorkShiftService(req CreateWorkShiftRequest) (*WorkShift, error) {
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	shift := &WorkShift{
		BranchID:  req.BranchID,
		Name:      req.Name,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsActive:  isActive,
	}

	if err := CreateWorkShift(shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func UpdateWorkShiftService(id int64, req UpdateWorkShiftRequest) (*WorkShift, error) {
	shift, err := GetWorkShiftByID(id)
	if err != nil {
		return nil, err
	}

	if req.BranchID != nil {
		shift.BranchID = *req.BranchID
	}
	if req.Name != nil {
		shift.Name = *req.Name
	}
	if req.StartTime != nil {
		shift.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		shift.EndTime = *req.EndTime
	}
	if req.IsActive != nil {
		shift.IsActive = *req.IsActive
	}

	if err := UpdateWorkShift(shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func DeleteWorkShiftService(id int64) error {
	return DeleteWorkShift(id)
}

// Shift Assignments
func GetAllShiftAssignmentsService() ([]ShiftAssignment, error) {
	return GetAllShiftAssignments()
}

func GetShiftAssignmentByIDService(id int64) (*ShiftAssignment, error) {
	return GetShiftAssignmentByID(id)
}

func CreateShiftAssignmentService(req CreateShiftAssignmentRequest) (*ShiftAssignment, error) {
	status := "scheduled"
	if req.Status != nil {
		status = *req.Status
	}
	isOvertime := false
	if req.IsOvertime != nil {
		isOvertime = *req.IsOvertime
	}

	assignment := &ShiftAssignment{
		ShiftID:    req.ShiftID,
		UserID:     req.UserID,
		WorkDate:   req.WorkDate,
		Status:     status,
		IsOvertime: isOvertime,
	}

	if err := CreateShiftAssignment(assignment); err != nil {
		return nil, err
	}
	return assignment, nil
}

func UpdateShiftAssignmentService(id int64, req UpdateShiftAssignmentRequest) (*ShiftAssignment, error) {
	assignment, err := GetShiftAssignmentByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		assignment.Status = *req.Status
	}
	if req.IsOvertime != nil {
		assignment.IsOvertime = *req.IsOvertime
	}

	if err := UpdateShiftAssignment(assignment); err != nil {
		return nil, err
	}
	return assignment, nil
}

func DeleteShiftAssignmentService(id int64) error {
	return DeleteShiftAssignment(id)
}

// Shift Handovers
func GetAllShiftHandoversService() ([]ShiftHandover, error) {
	return GetAllShiftHandovers()
}

func GetShiftHandoverByIDService(id string) (*ShiftHandover, error) {
	return GetShiftHandoverByID(id)
}

func CreateShiftHandoverService(req CreateShiftHandoverRequest) (*ShiftHandover, error) {
	handover := &ShiftHandover{
		AssignmentID:        req.AssignmentID,
		UserID:              req.UserID,
		BranchID:            req.BranchID,
		SystemCashAmount:    req.SystemCashAmount,
		SystemCODAmount:     req.SystemCODAmount,
		ActualCashSubmitted: req.ActualCashSubmitted,
		Status:              "pending_finance",
	}

	if err := CreateShiftHandover(handover); err != nil {
		return nil, err
	}
	return handover, nil
}

func UpdateShiftHandoverService(id string, req UpdateShiftHandoverRequest) (*ShiftHandover, error) {
	handover, err := GetShiftHandoverByID(id)
	if err != nil {
		return nil, err
	}

	if req.Status != nil {
		handover.Status = *req.Status
	}
	if req.FinanceTransactionID != nil {
		handover.FinanceTransactionID = req.FinanceTransactionID
	}

	if err := UpdateShiftHandover(handover); err != nil {
		return nil, err
	}
	return handover, nil
}

func DeleteShiftHandoverService(id string) error {
	return DeleteShiftHandover(id)
}
