package clinical_queues

import "time"

func GetAllClinicalQueuesService() ([]ClinicalQueue, error) {
	return GetAllClinicalQueues()
}

func GetClinicalQueueByIDService(id string) (*ClinicalQueue, error) {
	return GetClinicalQueueByID(id)
}

func CreateClinicalQueueService(req CreateClinicalQueueRequest) (*ClinicalQueue, error) {
	priorityLevel := "normal"
	if req.PriorityLevel != nil {
		priorityLevel = *req.PriorityLevel
	}

	queue := &ClinicalQueue{
		AppointmentID: req.AppointmentID,
		CustomerID:    req.CustomerID,
		DoctorID:      req.DoctorID,
		QueueNumber:   req.QueueNumber,
		Status:        "waiting",
		PriorityLevel: priorityLevel,
	}

	if err := CreateClinicalQueue(queue); err != nil {
		return nil, err
	}
	return queue, nil
}

func UpdateClinicalQueueService(id string, req UpdateClinicalQueueRequest) (*ClinicalQueue, error) {
	queue, err := GetClinicalQueueByID(id)
	if err != nil {
		return nil, err
	}

	if req.DoctorID != nil {
		queue.DoctorID = req.DoctorID
	}
	if req.Status != nil {
		queue.Status = *req.Status
	}
	if req.PriorityLevel != nil {
		queue.PriorityLevel = *req.PriorityLevel
	}
	
	now := time.Now()
	queue.UpdatedAt = &now

	if err := UpdateClinicalQueue(queue); err != nil {
		return nil, err
	}
	return queue, nil
}

func DeleteClinicalQueueService(id string) error {
	return DeleteClinicalQueue(id)
}
