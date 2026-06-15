import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { orderApi } from '../api';
import { OrdersCreateOrderRequest, OrdersUpdateOrderRequest } from '@namviet/shared-types/src/backend.d';

export const orderKeys = {
  all: ['orders'] as const,
  lists: () => [...orderKeys.all, 'list'] as const,
  list: (type?: string) => [...orderKeys.lists(), { type }] as const,
  details: () => [...orderKeys.all, 'detail'] as const,
  detail: (id: string) => [...orderKeys.details(), id] as const,
};

export const useOrders = (type?: string) => {
  return useQuery({
    queryKey: orderKeys.list(type),
    queryFn: () => orderApi.getAll(type),
  });
};

export const useOrder = (id: string) => {
  return useQuery({
    queryKey: orderKeys.detail(id),
    queryFn: () => orderApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateOrder = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: OrdersCreateOrderRequest) => orderApi.create(data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: orderKeys.lists() });
    },
  });
};

export const useUpdateOrder = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; data: OrdersUpdateOrderRequest }) => orderApi.updateStatus(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: orderKeys.lists() });
      queryClient.invalidateQueries({ queryKey: orderKeys.detail(variables.id) });
    },
  });
};
