import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { clinicalApi, CreateAppointmentRequest, UpdateAppointmentRequest, CreateClinicalQueueRequest, UpdateClinicalQueueRequest, CreateMedicalVisitRequest, UpdateMedicalVisitRequest } from '../api';
import { toast } from 'sonner';

// --- APPOINTMENTS ---
export const useAppointments = () => {
  return useQuery({
    queryKey: ['appointments'],
    queryFn: clinicalApi.getAppointments,
  });
};

export const useCreateAppointment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateAppointmentRequest) => clinicalApi.createAppointment(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['appointments'] });
      toast.success('Đặt lịch hẹn thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateAppointment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateAppointmentRequest }) => clinicalApi.updateAppointment(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['appointments'] });
      toast.success('Cập nhật lịch hẹn thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteAppointment = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => clinicalApi.deleteAppointment(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['appointments'] });
      toast.success('Hủy lịch hẹn thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- CLINICAL QUEUES ---
export const useQueues = () => {
  return useQuery({
    queryKey: ['clinicalQueues'],
    queryFn: clinicalApi.getQueues,
  });
};

export const useCreateQueue = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateClinicalQueueRequest) => clinicalApi.createQueue(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clinicalQueues'] });
      toast.success('Đưa vào hàng đợi thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateQueue = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateClinicalQueueRequest }) => clinicalApi.updateQueue(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clinicalQueues'] });
      toast.success('Cập nhật hàng đợi thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteQueue = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => clinicalApi.deleteQueue(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['clinicalQueues'] });
      toast.success('Xóa khỏi hàng đợi thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

// --- MEDICAL VISITS ---
export const useVisits = () => {
  return useQuery({
    queryKey: ['medicalVisits'],
    queryFn: clinicalApi.getVisits,
  });
};

export const useCreateVisit = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateMedicalVisitRequest) => clinicalApi.createVisit(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medicalVisits'] });
      toast.success('Tạo hồ sơ khám thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateVisit = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMedicalVisitRequest }) => clinicalApi.updateVisit(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medicalVisits'] });
      toast.success('Cập nhật hồ sơ khám thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteVisit = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => clinicalApi.deleteVisit(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['medicalVisits'] });
      toast.success('Xóa hồ sơ khám thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};
