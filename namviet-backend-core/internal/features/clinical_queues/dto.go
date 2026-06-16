package clinical_queues

type CreateClinicalQueueRequest struct {
	AppointmentID *string `json:"appointment_id"`
	CustomerID    int64   `json:"customer_id" binding:"required"`
	DoctorID      *string `json:"doctor_id"`
	QueueNumber   int     `json:"queue_number" binding:"required"`
	PriorityLevel *string `json:"priority_level"`
}

type UpdateClinicalQueueRequest struct {
	DoctorID      *string `json:"doctor_id"`
	Status        *string `json:"status"`
	PriorityLevel *string `json:"priority_level"`
}
