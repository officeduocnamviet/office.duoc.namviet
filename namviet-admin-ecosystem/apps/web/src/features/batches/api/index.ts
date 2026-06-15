import { apiClient } from '@/lib/axios';
import { 
  BatchesBatch, 
  BatchesCreateBatchRequest, 
  BatchesUpdateBatchRequest 
} from '@namviet/shared-types/src/backend.d';

export const batchApi = {
  getAll: async (): Promise<BatchesBatch[]> => {
    const response = await apiClient.get<BatchesBatch[]>('/batches');
    return response.data;
  },

  getById: async (id: number): Promise<BatchesBatch> => {
    const response = await apiClient.get<BatchesBatch>(`/batches/${id}`);
    return response.data;
  },

  create: async (data: BatchesCreateBatchRequest): Promise<BatchesBatch> => {
    const response = await apiClient.post<BatchesBatch>('/batches', data);
    return response.data;
  },

  update: async (id: number, data: BatchesUpdateBatchRequest): Promise<BatchesBatch> => {
    const response = await apiClient.put<BatchesBatch>(`/batches/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/batches/${id}`);
  }
};
