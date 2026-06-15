import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { manufacturerApi } from '../api';
import { 
  ManufacturersCreateManufacturerRequest, 
  ManufacturersUpdateManufacturerRequest 
} from '@namviet/shared-types/src/backend.d';

export const manufacturerKeys = {
  all: ['manufacturers'] as const,
  details: () => [...manufacturerKeys.all, 'detail'] as const,
  detail: (id: number) => [...manufacturerKeys.details(), id] as const,
};

export const useManufacturers = () => {
  return useQuery({
    queryKey: manufacturerKeys.all,
    queryFn: manufacturerApi.getAll,
  });
};

export const useManufacturer = (id: number) => {
  return useQuery({
    queryKey: manufacturerKeys.detail(id),
    queryFn: () => manufacturerApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateManufacturer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: ManufacturersCreateManufacturerRequest) => manufacturerApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: manufacturerKeys.all });
    },
  });
};

export const useUpdateManufacturer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: ManufacturersUpdateManufacturerRequest }) => 
      manufacturerApi.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: manufacturerKeys.all });
      queryClient.invalidateQueries({ queryKey: manufacturerKeys.detail(variables.id) });
    },
  });
};

export const useDeleteManufacturer = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => manufacturerApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: manufacturerKeys.all });
    },
  });
};
