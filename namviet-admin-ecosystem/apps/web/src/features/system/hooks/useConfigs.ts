import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { systemConfigApi, SystemConfig } from '../api/configApi';

export const configKeys = {
  all: ['systemConfigs'] as const,
  detail: (key: string) => ['systemConfigs', key] as const,
};

export const useSystemConfigs = () => {
  return useQuery({
    queryKey: configKeys.all,
    queryFn: systemConfigApi.getAll,
  });
};

export const useCreateSystemConfig = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: systemConfigApi.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.all });
    },
  });
};

export const useUpdateSystemConfig = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ key, data }: { key: string; data: any }) => systemConfigApi.update(key, data),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: configKeys.all });
      queryClient.invalidateQueries({ queryKey: configKeys.detail(variables.key) });
    },
  });
};

export const useDeleteSystemConfig = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: systemConfigApi.delete,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: configKeys.all });
    },
  });
};
