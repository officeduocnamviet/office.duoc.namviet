import { apiClient as api } from '@/lib/axios';
import {
  AppointmentsAppointment as Appointment,
  AppointmentsCreateAppointmentRequest as CreateAppointmentRequest,
  AppointmentsUpdateAppointmentRequest as UpdateAppointmentRequest,
  ClinicalQueuesClinicalQueue as ClinicalQueue,
  ClinicalQueuesCreateClinicalQueueRequest as CreateClinicalQueueRequest,
  ClinicalQueuesUpdateClinicalQueueRequest as UpdateClinicalQueueRequest,
  MedicalVisitsMedicalVisit as MedicalVisit,
  MedicalVisitsCreateMedicalVisitRequest as CreateMedicalVisitRequest,
  MedicalVisitsUpdateMedicalVisitRequest as UpdateMedicalVisitRequest
} from '@namviet/shared-types/src/backend.d';

export type {
  Appointment,
  CreateAppointmentRequest,
  UpdateAppointmentRequest,
  ClinicalQueue,
  CreateClinicalQueueRequest,
  UpdateClinicalQueueRequest,
  MedicalVisit,
  CreateMedicalVisitRequest,
  UpdateMedicalVisitRequest
};

export const clinicalApi = {
  // Appointments
  getAppointments: () => api.get<Appointment[]>('/api/appointments').then(res => res.data),
  getAppointment: (id: string) => api.get<Appointment>(`/api/appointments/${id}`).then(res => res.data),
  createAppointment: (data: CreateAppointmentRequest) => api.post<Appointment>('/api/appointments', data).then(res => res.data),
  updateAppointment: (id: string, data: UpdateAppointmentRequest) => api.put<Appointment>(`/api/appointments/${id}`, data).then(res => res.data),
  deleteAppointment: (id: string) => api.delete(`/api/appointments/${id}`).then(res => res.data),

  // Clinical Queues
  getQueues: () => api.get<ClinicalQueue[]>('/api/clinical-queues').then(res => res.data),
  getQueue: (id: string) => api.get<ClinicalQueue>(`/api/clinical-queues/${id}`).then(res => res.data),
  createQueue: (data: CreateClinicalQueueRequest) => api.post<ClinicalQueue>('/api/clinical-queues', data).then(res => res.data),
  updateQueue: (id: string, data: UpdateClinicalQueueRequest) => api.put<ClinicalQueue>(`/api/clinical-queues/${id}`, data).then(res => res.data),
  deleteQueue: (id: string) => api.delete(`/api/clinical-queues/${id}`).then(res => res.data),

  // Medical Visits
  getVisits: () => api.get<MedicalVisit[]>('/api/medical-visits').then(res => res.data),
  getVisit: (id: string) => api.get<MedicalVisit>(`/api/medical-visits/${id}`).then(res => res.data),
  createVisit: (data: CreateMedicalVisitRequest) => api.post<MedicalVisit>('/api/medical-visits', data).then(res => res.data),
  updateVisit: (id: string, data: UpdateMedicalVisitRequest) => api.put<MedicalVisit>(`/api/medical-visits/${id}`, data).then(res => res.data),
  deleteVisit: (id: string) => api.delete(`/api/medical-visits/${id}`).then(res => res.data),
};
