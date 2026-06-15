import { apiClient } from '@/lib/axios';

export interface SystemConfig {
  id?: string;
  config_key?: string;
  config_value?: any; // The API might return JSON
  description?: string;
  updated_at?: string;
  updated_by?: string;
}

export const systemConfigApi = {
  getAll: async (): Promise<SystemConfig[]> => {
    const response = await apiClient.get<SystemConfig[]>('/system-configs');
    return response.data;
  },

  getByKey: async (key: string): Promise<SystemConfig> => {
    const response = await apiClient.get<SystemConfig>(`/system-configs/${key}`);
    return response.data;
  },

  create: async (data: any): Promise<SystemConfig> => {
    const response = await apiClient.post<SystemConfig>('/system-configs', data);
    return response.data;
  },

  update: async (key: string, data: any): Promise<SystemConfig> => {
    const response = await apiClient.put<SystemConfig>(`/system-configs/${key}`, data);
    return response.data;
  },

  delete: async (key: string): Promise<void> => {
    await apiClient.delete(`/system-configs/${key}`);
  }
};
