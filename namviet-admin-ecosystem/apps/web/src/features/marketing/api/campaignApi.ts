import { apiClient } from '@/lib/axios';

export type MarketingCampaign = import('@namviet/shared-types/src/backend').MarketingCampaignsMarketingCampaign;
export type CreateMarketingCampaignRequest = import('@namviet/shared-types/src/backend').MarketingCampaignsCreateMarketingCampaignRequest;
export type UpdateMarketingCampaignRequest = import('@namviet/shared-types/src/backend').MarketingCampaignsUpdateMarketingCampaignRequest;

export const campaignApi = {
  getCampaigns: async (): Promise<MarketingCampaign[]> => {
    const response = await apiClient.get<MarketingCampaign[]>('/marketing-campaigns');
    return response.data;
  },

  createCampaign: async (data: CreateMarketingCampaignRequest): Promise<MarketingCampaign> => {
    const response = await apiClient.post<MarketingCampaign>('/marketing-campaigns', data);
    return response.data;
  },

  updateCampaign: async (id: string, data: UpdateMarketingCampaignRequest): Promise<MarketingCampaign> => {
    const response = await apiClient.put<MarketingCampaign>(`/marketing-campaigns/${id}`, data);
    return response.data;
  },

  deleteCampaign: async (id: string): Promise<void> => {
    await apiClient.delete(`/marketing-campaigns/${id}`);
  }
};
