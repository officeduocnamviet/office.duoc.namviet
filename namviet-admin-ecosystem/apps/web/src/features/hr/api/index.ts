import { apiClient as api } from '@/lib/axios';
import {
  EmployeesEmployee as Employee,
  EmployeesCreateEmployeeRequest as CreateEmployeeRequest,
  EmployeesUpdateEmployeeRequest as UpdateEmployeeRequest,
  TimeAttendanceTimeAttendance as TimeAttendance,
  TimeAttendanceCreateTimeAttendanceRequest as CreateTimeAttendanceRequest,
  TimeAttendanceUpdateTimeAttendanceRequest as UpdateTimeAttendanceRequest,
  PayrollsPayroll as Payroll,
  PayrollsCreatePayrollRequest as CreatePayrollRequest,
  PayrollsUpdatePayrollRequest as UpdatePayrollRequest
} from '@namviet/shared-types/src/backend.d';

export type {
  Employee,
  CreateEmployeeRequest,
  UpdateEmployeeRequest,
  TimeAttendance,
  CreateTimeAttendanceRequest,
  UpdateTimeAttendanceRequest,
  Payroll,
  CreatePayrollRequest,
  UpdatePayrollRequest
};

export type AttendanceLog = import('@namviet/shared-types/src/backend').AttendanceLogsAttendanceLog;
export type CreateAttendanceLogRequest = import('@namviet/shared-types/src/backend').AttendanceLogsCreateAttendanceLogRequest;
export type UpdateAttendanceLogRequest = import('@namviet/shared-types/src/backend').AttendanceLogsUpdateAttendanceLogRequest;

export type WorkShift = import('@namviet/shared-types/src/backend').WorkShiftsWorkShift;
export type CreateWorkShiftRequest = import('@namviet/shared-types/src/backend').WorkShiftsCreateWorkShiftRequest;
export type UpdateWorkShiftRequest = import('@namviet/shared-types/src/backend').WorkShiftsUpdateWorkShiftRequest;

export type ShiftAssignment = import('@namviet/shared-types/src/backend').WorkShiftsShiftAssignment;
export type CreateShiftAssignmentRequest = import('@namviet/shared-types/src/backend').WorkShiftsCreateShiftAssignmentRequest;

export type ShiftHandover = import('@namviet/shared-types/src/backend').WorkShiftsShiftHandover;
export type CreateShiftHandoverRequest = import('@namviet/shared-types/src/backend').WorkShiftsCreateShiftHandoverRequest;

export type EmploymentContract = import('@namviet/shared-types/src/backend').EmploymentContractsEmploymentContract;
export type CreateEmploymentContractRequest = import('@namviet/shared-types/src/backend').EmploymentContractsCreateEmploymentContractRequest;

export type TrainingCourse = import('@namviet/shared-types/src/backend').TrainingCoursesTrainingCourse;
export type CreateTrainingCourseRequest = import('@namviet/shared-types/src/backend').TrainingCoursesCreateTrainingCourseRequest;

export const hrApi = {
  // Employees
  getEmployees: () => api.get<Employee[]>('/employees').then(res => res.data),
  getEmployee: (id: string) => api.get<Employee>(`/employees/${id}`).then(res => res.data),
  createEmployee: (data: CreateEmployeeRequest) => api.post<Employee>('/employees', data).then(res => res.data),
  updateEmployee: (id: string, data: UpdateEmployeeRequest) => api.put<Employee>(`/employees/${id}`, data).then(res => res.data),
  deleteEmployee: (id: string) => api.delete(`/employees/${id}`).then(res => res.data),

  // Time Attendance
  getAttendances: () => api.get<TimeAttendance[]>('/time-attendances').then(res => res.data),
  getAttendance: (id: string) => api.get<TimeAttendance>(`/time-attendances/${id}`).then(res => res.data),
  createAttendance: (data: CreateTimeAttendanceRequest) => api.post<TimeAttendance>('/time-attendances', data).then(res => res.data),
  updateAttendance: (id: string, data: UpdateTimeAttendanceRequest) => api.put<TimeAttendance>(`/time-attendances/${id}`, data).then(res => res.data),
  deleteAttendance: (id: string) => api.delete(`/time-attendances/${id}`).then(res => res.data),

  // Payrolls
  getPayrolls: () => api.get<Payroll[]>('/payrolls').then(res => res.data),
  getPayroll: (id: string) => api.get<Payroll>(`/payrolls/${id}`).then(res => res.data),
  createPayroll: (data: CreatePayrollRequest) => api.post<Payroll>('/payrolls', data).then(res => res.data),
  updatePayroll: (id: string, data: UpdatePayrollRequest) => api.put<Payroll>(`/payrolls/${id}`, data).then(res => res.data),
  deletePayroll: (id: string) => api.delete(`/payrolls/${id}`).then(res => res.data),

  // Attendance Logs (Check-in/Check-out)
  getAttendanceLogs: () => api.get<AttendanceLog[]>('/attendance-logs').then(res => res.data),
  createAttendanceLog: (data: CreateAttendanceLogRequest) => api.post<AttendanceLog>('/attendance-logs', data).then(res => res.data),
  updateAttendanceLog: (id: string, data: UpdateAttendanceLogRequest) => api.put<AttendanceLog>(`/attendance-logs/${id}`, data).then(res => res.data),

  // Work Shifts
  getWorkShifts: () => api.get<WorkShift[]>('/work-shifts').then(res => res.data),
  createWorkShift: (data: CreateWorkShiftRequest) => api.post<WorkShift>('/work-shifts', data).then(res => res.data),
  
  // Shift Assignments
  getShiftAssignments: () => api.get<ShiftAssignment[]>('/shift-assignments').then(res => res.data),
  createShiftAssignment: (data: CreateShiftAssignmentRequest) => api.post<ShiftAssignment>('/shift-assignments', data).then(res => res.data),

  // Shift Handovers
  getShiftHandovers: () => api.get<ShiftHandover[]>('/shift-handovers').then(res => res.data),
  createShiftHandover: (data: CreateShiftHandoverRequest) => api.post<ShiftHandover>('/shift-handovers', data).then(res => res.data),

  // Employment Contracts
  getEmploymentContracts: () => api.get<EmploymentContract[]>('/employment-contracts').then(res => res.data),
  createEmploymentContract: (data: CreateEmploymentContractRequest) => api.post<EmploymentContract>('/employment-contracts', data).then(res => res.data),

  // Training Courses
  getTrainingCourses: () => api.get<TrainingCourse[]>('/training-courses').then(res => res.data),
  createTrainingCourse: (data: CreateTrainingCourseRequest) => api.post<TrainingCourse>('/training-courses', data).then(res => res.data),
};
