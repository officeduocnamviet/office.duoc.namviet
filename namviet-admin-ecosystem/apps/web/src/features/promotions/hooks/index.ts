import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { promotionApi } from '../api';
import { PromotionsCreatePromotionRequest, PromotionsUpdatePromotionRequest } from '@namviet/shared-types/src/backend.d';

export const promotionKeys = {
  all: ['promotions'] as const,
  lists: () => [...promotionKeys.all, 'list'] as const,
  list: (filters: string) => [...promotionKeys.lists(), { filters }] as const,
  details: () => [...promotionKeys.all, 'detail'] as const,
  detail: (id: string) => [...promotionKeys.details(), id] as const,
};

export const usePromotions = () => {
  return useQuery({
    queryKey: promotionKeys.lists(),
    queryFn: () => promotionApi.getAll(),
  });
};

export const usePromotion = (id: string) => {
  return useQuery({
    queryKey: promotionKeys.detail(id),
    queryFn: () => promotionApi.getById(id),
    enabled: !!id,
  });
};

export const useCreatePromotion = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: PromotionsCreatePromotionRequest) => promotionApi.create(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: promotionKeys.lists() });
    },
  });
};

export const useUpdatePromotion = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { id: string; data: PromotionsUpdatePromotionRequest }) => promotionApi.update(payload),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: promotionKeys.lists() });
      queryClient.invalidateQueries({ queryKey: promotionKeys.detail(variables.id) });
    },
  });
};

export const useDeletePromotion = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => promotionApi.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: promotionKeys.lists() });
    },
  });
};
