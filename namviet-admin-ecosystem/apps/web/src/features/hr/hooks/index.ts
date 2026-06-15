import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { hrApi, CreateEmployeeRequest, UpdateEmployeeRequest, CreateTimeAttendanceRequest, UpdateTimeAttendanceRequest, CreatePayrollRequest, UpdatePayrollRequest } from '../api';
import { toast } from 'sonner';

// --- EMPLOYEES ---
export const useEmployees = () => {
  return useQuery({
    queryKey: ['employees'],
    queryFn: hrApi.getEmployees,
  });
};

export const useCreateEmployee = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateEmployeeRequest) => hrApi.createEmployee(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees'] });
      toast.success('Thêm nhân viên thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateEmployee = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateEmployeeRequest }) => hrApi.updateEmployee(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees'] });
      toast.success('Cập nhật nhân viên thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteEmployee = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.deleteEmployee(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employees'] });
      toast.success('Xóa nhân viên thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- TIME ATTENDANCE ---
export const useAttendances = () => {
  return useQuery({
    queryKey: ['attendances'],
    queryFn: hrApi.getAttendances,
  });
};

export const useCreateAttendance = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateTimeAttendanceRequest) => hrApi.createAttendance(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attendances'] });
      toast.success('Chấm công thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateAttendance = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateTimeAttendanceRequest }) => hrApi.updateAttendance(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attendances'] });
      toast.success('Cập nhật chấm công thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteAttendance = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.deleteAttendance(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attendances'] });
      toast.success('Xóa bản ghi chấm công thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- PAYROLLS ---
export const usePayrolls = () => {
  return useQuery({
    queryKey: ['payrolls'],
    queryFn: hrApi.getPayrolls,
  });
};

export const useCreatePayroll = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreatePayrollRequest) => hrApi.createPayroll(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payrolls'] });
      toast.success('Tạo bảng lương thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdatePayroll = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdatePayrollRequest }) => hrApi.updatePayroll(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payrolls'] });
      toast.success('Cập nhật bảng lương thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeletePayroll = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => hrApi.deletePayroll(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payrolls'] });
      toast.success('Xóa bảng lương thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- ATTENDANCE LOGS (CHECK-IN/OUT) ---
export const useAttendanceLogs = () => {
  return useQuery({
    queryKey: ['attendance_logs'],
    queryFn: hrApi.getAttendanceLogs,
  });
};

export const useCreateAttendanceLog = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: import('../api').CreateAttendanceLogRequest) => hrApi.createAttendanceLog(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attendance_logs'] });
      toast.success('Check-in thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateAttendanceLog = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: import('../api').UpdateAttendanceLogRequest }) => hrApi.updateAttendanceLog(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['attendance_logs'] });
      toast.success('Check-out thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- WORK SHIFTS & ASSIGNMENTS ---
export const useWorkShifts = () => {
  return useQuery({
    queryKey: ['work_shifts'],
    queryFn: hrApi.getWorkShifts,
  });
};

export const useShiftAssignments = () => {
  return useQuery({
    queryKey: ['shift_assignments'],
    queryFn: hrApi.getShiftAssignments,
  });
};

export const useCreateShiftAssignment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: import('../api').CreateShiftAssignmentRequest) => hrApi.createShiftAssignment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['shift_assignments'] });
      toast.success('Phân ca thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useShiftHandovers = () => {
  return useQuery({
    queryKey: ['shift_handovers'],
    queryFn: hrApi.getShiftHandovers,
  });
};

export const useCreateShiftHandover = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: import('../api').CreateShiftHandoverRequest) => hrApi.createShiftHandover(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['shift_handovers'] });
      toast.success('Gửi yêu cầu bàn giao ca thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- EMPLOYMENT CONTRACTS ---
export const useEmploymentContracts = () => {
  return useQuery({
    queryKey: ['employment_contracts'],
    queryFn: hrApi.getEmploymentContracts,
  });
};

export const useCreateEmploymentContract = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: import('../api').CreateEmploymentContractRequest) => hrApi.createEmploymentContract(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['employment_contracts'] });
      toast.success('Tạo hợp đồng thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- TRAINING COURSES ---
export const useTrainingCourses = () => {
  return useQuery({
    queryKey: ['training_courses'],
    queryFn: hrApi.getTrainingCourses,
  });
};

export const useCreateTrainingCourse = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: import('../api').CreateTrainingCourseRequest) => hrApi.createTrainingCourse(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['training_courses'] });
      toast.success('Tạo khóa đào tạo thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};
