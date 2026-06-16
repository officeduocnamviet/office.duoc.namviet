package medical_visits

import "time"

func GetAllMedicalVisitsService() ([]MedicalVisit, error) {
	return GetAllMedicalVisits()
}

func GetMedicalVisitByIDService(id string) (*MedicalVisit, error) {
	return GetMedicalVisitByID(id)
}

func CreateMedicalVisitService(req CreateMedicalVisitRequest) (*MedicalVisit, error) {
	visit := &MedicalVisit{
		AppointmentID: req.AppointmentID,
		CustomerID:    req.CustomerID,
		DoctorID:      req.DoctorID,
		Symptoms:      req.Symptoms,
		Status:        "in_progress",
	}

	if err := CreateMedicalVisit(visit); err != nil {
		return nil, err
	}
	return visit, nil
}

func UpdateMedicalVisitService(id string, req UpdateMedicalVisitRequest) (*MedicalVisit, error) {
	visit, err := GetMedicalVisitByID(id)
	if err != nil {
		return nil, err
	}

	if req.Temperature != nil {
		visit.Temperature = req.Temperature
	}
	if req.Pulse != nil {
		visit.Pulse = req.Pulse
	}
	if req.SpO2 != nil {
		visit.SpO2 = req.SpO2
	}
	if req.BPSystolic != nil {
		visit.BPSystolic = req.BPSystolic
	}
	if req.BPDiastolic != nil {
		visit.BPDiastolic = req.BPDiastolic
	}
	if req.Weight != nil {
		visit.Weight = req.Weight
	}
	if req.Height != nil {
		visit.Height = req.Height
	}
	if req.Symptoms != nil {
		visit.Symptoms = req.Symptoms
	}
	if req.ExaminationSummary != nil {
		visit.ExaminationSummary = req.ExaminationSummary
	}
	if req.Diagnosis != nil {
		visit.Diagnosis = req.Diagnosis
	}
	if req.ICDCode != nil {
		visit.ICDCode = req.ICDCode
	}
	if req.DoctorNotes != nil {
		visit.DoctorNotes = req.DoctorNotes
	}
	if req.RedFlags != nil {
		visit.RedFlags = *req.RedFlags
	}
	if req.Status != nil {
		visit.Status = *req.Status
	}
	
	now := time.Now()
	visit.UpdatedAt = &now

	if err := UpdateMedicalVisit(visit); err != nil {
		return nil, err
	}
	return visit, nil
}

func DeleteMedicalVisitService(id string) error {
	return DeleteMedicalVisit(id)
}
