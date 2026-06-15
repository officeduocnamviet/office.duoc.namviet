import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { campaignApi, CreateMarketingCampaignRequest, UpdateMarketingCampaignRequest } from '../api/campaignApi';
import { toast } from 'sonner';

export const useCampaigns = () => {
  return useQuery({
    queryKey: ['campaigns'],
    queryFn: campaignApi.getCampaigns,
  });
};

export const useCreateCampaign = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateMarketingCampaignRequest) => campaignApi.createCampaign(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] });
      toast.success('Tạo chiến dịch thành công');
    },
    onError: (error: any) => {
      toast.error(`Lỗi: ${error.message}`);
    },
  });
};

export const useUpdateCampaign = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateMarketingCampaignRequest }) => campaignApi.updateCampaign(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['campaigns'] });
      toast.success('Cập nhật chiến dịch thành công');
    },
    onError: (error: any) => {
      toast.error(`Lỗi: ${error.message}`);
    },
  });
};
