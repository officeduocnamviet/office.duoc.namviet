import { apiClient } from '@/lib/axios';

export type InternalChannel = import('@namviet/shared-types/src/backend.d').InternalCommunicationsInternalChannel;
export type CreateInternalChannelRequest = import('@namviet/shared-types/src/backend.d').InternalCommunicationsCreateInternalChannelRequest;
export type UpdateInternalChannelRequest = import('@namviet/shared-types/src/backend.d').InternalCommunicationsUpdateInternalChannelRequest;

export type InternalMessage = import('@namviet/shared-types/src/backend.d').InternalCommunicationsInternalMessage;
export type CreateInternalMessageRequest = import('@namviet/shared-types/src/backend.d').InternalCommunicationsCreateInternalMessageRequest;

export const chatApi = {
  getChannels: async (): Promise<InternalChannel[]> => {
    const { data } = await apiClient.get<InternalChannel[]>('/api/internal-channels');
    return data;
  },
  
  createChannel: async (req: CreateInternalChannelRequest): Promise<InternalChannel> => {
    const { data } = await apiClient.post<InternalChannel>('/api/internal-channels', req);
    return data;
  },

  getMessages: async (channelId: number): Promise<InternalMessage[]> => {
    const { data } = await apiClient.get<InternalMessage[]>(`/api/internal-channels/${channelId}/messages`);
    return data;
  },

  sendMessage: async (req: CreateInternalMessageRequest): Promise<InternalMessage> => {
    const { data } = await apiClient.post<InternalMessage>('/api/internal-messages', req);
    return data;
  }
};
