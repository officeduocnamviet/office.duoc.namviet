import { apiClient } from '@/lib/axios';
import { 
  PromotionsPromotion, 
  PromotionsCreatePromotionRequest, 
  PromotionsUpdatePromotionRequest 
} from '@namviet/shared-types/src/backend.d';

export const promotionApi = {
  getAll: async (): Promise<PromotionsPromotion[]> => {
    const { data } = await apiClient.get<PromotionsPromotion[]>('/promotions');
    return data;
  },

  getById: async (id: string): Promise<PromotionsPromotion> => {
    const { data } = await apiClient.get<PromotionsPromotion>(`/promotions/${id}`);
    return data;
  },

  create: async (request: PromotionsCreatePromotionRequest): Promise<PromotionsPromotion> => {
    const { data } = await apiClient.post<PromotionsPromotion>('/promotions', request);
    return data;
  },

  update: async ({ id, data }: { id: string; data: PromotionsUpdatePromotionRequest }): Promise<PromotionsPromotion> => {
    const { data: resData } = await apiClient.put<PromotionsPromotion>(`/promotions/${id}`, data);
    return resData;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/promotions/${id}`);
  }
};
