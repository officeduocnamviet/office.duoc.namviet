package appointments

import "time"

func GetAllAppointmentsService() ([]Appointment, error) {
	return GetAllAppointments()
}

func GetAppointmentByIDService(id string) (*Appointment, error) {
	return GetAppointmentByID(id)
}

func CreateAppointmentService(req CreateAppointmentRequest) (*Appointment, error) {
	appointment := &Appointment{
		CustomerID:      req.CustomerID,
		DoctorID:        req.DoctorID,
		RoomID:          req.RoomID,
		ServiceType:     req.ServiceType,
		AppointmentTime: req.AppointmentTime,
		Symptoms:        req.Symptoms,
		Note:            req.Note,
		Status:          "pending",
	}

	if err := CreateAppointment(appointment); err != nil {
		return nil, err
	}
	return appointment, nil
}

func UpdateAppointmentService(id string, req UpdateAppointmentRequest) (*Appointment, error) {
	appointment, err := GetAppointmentByID(id)
	if err != nil {
		return nil, err
	}

	if req.DoctorID != nil {
		appointment.DoctorID = req.DoctorID
	}
	if req.RoomID != nil {
		appointment.RoomID = req.RoomID
	}
	if req.ServiceType != nil {
		appointment.ServiceType = *req.ServiceType
	}
	if req.AppointmentTime != nil {
		appointment.AppointmentTime = *req.AppointmentTime
	}
	if req.CheckInTime != nil {
		appointment.CheckInTime = req.CheckInTime
	}
	if req.Status != nil {
		appointment.Status = *req.Status
	}
	if req.Symptoms != nil {
		appointment.Symptoms = *req.Symptoms
	}
	if req.Note != nil {
		appointment.Note = req.Note
	}
	
	now := time.Now()
	appointment.UpdatedAt = &now

	if err := UpdateAppointment(appointment); err != nil {
		return nil, err
	}
	return appointment, nil
}

func DeleteAppointmentService(id string) error {
	return DeleteAppointment(id)
}
