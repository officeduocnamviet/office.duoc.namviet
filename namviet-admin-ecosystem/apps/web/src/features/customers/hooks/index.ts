import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { customerApi } from '../api';
import { CustomersCreateCustomerRequest, CustomersUpdateCustomerRequest } from '@namviet/shared-types/src/backend.d';

export const customerKeys = {
  all: ['customers'] as const,
  lists: () => [...customerKeys.all, 'list'] as const,
  list: (filters: string) => [...customerKeys.lists(), { filters }] as const,
  details: () => [...customerKeys.all, 'detail'] as const,
  detail: (id: number) => [...customerKeys.details(), id] as const,
};

export const useCustomers = () => {
  return useQuery({
    queryKey: customerKeys.lists(),
    queryFn: () => customerApi.getAll(),
  });
};

export const useCustomer = (id: number) => {
  return useQuery({
    queryKey: customerKeys.detail(id),
    queryFn: () => customerApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateCustomer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CustomersCreateCustomerRequest) => customerApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
    },
  });
};

export const useUpdateCustomer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: number; data: CustomersUpdateCustomerRequest }) => customerApi.update(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
      queryClient.invalidateQueries({ queryKey: customerKeys.detail(variables.id) });
    },
  });
};

export const useDeleteCustomer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => customerApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: customerKeys.lists() });
    },
  });
};

export const useCustomerVaccinations = (customerId: number | undefined) => {
  return useQuery({
    queryKey: ['customer_vaccinations', customerId],
    queryFn: () => customerApi.getVaccinations(customerId!),
    enabled: !!customerId,
  });
};

export const useCustomerVouchers = (customerId: number | undefined) => {
  return useQuery({
    queryKey: ['customer_vouchers', customerId],
    queryFn: () => customerApi.getVouchers(customerId!),
    enabled: !!customerId,
  });
};
