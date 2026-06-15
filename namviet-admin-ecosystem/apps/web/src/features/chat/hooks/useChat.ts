import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { chatApi, CreateInternalChannelRequest, CreateInternalMessageRequest, InternalMessage } from '../api/chatApi';
import { toast } from 'sonner';

export const useChannels = () => {
  return useQuery({
    queryKey: ['chat_channels'],
    queryFn: chatApi.getChannels,
  });
};

export const useCreateChannel = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateInternalChannelRequest) => chatApi.createChannel(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chat_channels'] });
      toast.success('Tạo kênh thành công');
    },
    onError: (err: any) => {
      toast.error(err.response?.data?.error || 'Lỗi tạo kênh');
    }
  });
};

export const useMessages = (channelId: number | undefined) => {
  return useQuery({
    queryKey: ['chat_messages', channelId],
    queryFn: () => chatApi.getMessages(channelId!),
    enabled: !!channelId,
    // Disable polling as requested
    refetchInterval: false,
    refetchOnWindowFocus: false,
  });
};

export const useSendMessage = () => {
  return useMutation({
    mutationFn: (data: CreateInternalMessageRequest) => chatApi.sendMessage(data),
    onError: (err: any) => {
      toast.error(err.response?.data?.error || 'Lỗi gửi tin nhắn');
    }
  });
};
