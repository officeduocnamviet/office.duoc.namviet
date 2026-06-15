import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { batchApi } from '../api';
import { 
  BatchesCreateBatchRequest, 
  BatchesUpdateBatchRequest 
} from '@namviet/shared-types/src/backend.d';

export const batchKeys = {
  all: ['batches'] as const,
  details: () => [...batchKeys.all, 'detail'] as const,
  detail: (id: number) => [...batchKeys.details(), id] as const,
};

export const useBatches = () => {
  return useQuery({
    queryKey: batchKeys.all,
    queryFn: batchApi.getAll,
  });
};

export const useBatch = (id: number) => {
  return useQuery({
    queryKey: batchKeys.detail(id),
    queryFn: () => batchApi.getById(id),
    enabled: !!id,
  });
};

export const useCreateBatch = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: BatchesCreateBatchRequest) => batchApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: batchKeys.all });
    },
  });
};

export const useUpdateBatch = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: number; data: BatchesUpdateBatchRequest }) => 
      batchApi.update(id, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: batchKeys.all });
      queryClient.invalidateQueries({ queryKey: batchKeys.detail(variables.id) });
    },
  });
};

export const useDeleteBatch = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => batchApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: batchKeys.all });
    },
  });
};
