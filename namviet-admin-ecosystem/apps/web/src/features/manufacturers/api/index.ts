import { apiClient } from '@/lib/axios';
import { 
  ManufacturersManufacturer, 
  ManufacturersCreateManufacturerRequest, 
  ManufacturersUpdateManufacturerRequest 
} from '@namviet/shared-types/src/backend.d';

export const manufacturerApi = {
  getAll: async (): Promise<ManufacturersManufacturer[]> => {
    const response = await apiClient.get<ManufacturersManufacturer[]>('/manufacturers');
    return response.data;
  },

  getById: async (id: number): Promise<ManufacturersManufacturer> => {
    const response = await apiClient.get<ManufacturersManufacturer>(`/manufacturers/${id}`);
    return response.data;
  },

  create: async (data: ManufacturersCreateManufacturerRequest): Promise<ManufacturersManufacturer> => {
    const response = await apiClient.post<ManufacturersManufacturer>('/manufacturers', data);
    return response.data;
  },

  update: async (id: number, data: ManufacturersUpdateManufacturerRequest): Promise<ManufacturersManufacturer> => {
    const response = await apiClient.put<ManufacturersManufacturer>(`/manufacturers/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/manufacturers/${id}`);
  }
};
