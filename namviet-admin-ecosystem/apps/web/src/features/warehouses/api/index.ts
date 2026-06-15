import { apiClient } from '@/lib/axios';
import { 
  WarehousesWarehouse, 
  WarehousesCreateWarehouseRequest, 
  WarehousesUpdateWarehouseRequest 
} from '@namviet/shared-types/src/backend.d';

export const warehouseApi = {
  getAll: async (): Promise<WarehousesWarehouse[]> => {
    const response = await apiClient.get<WarehousesWarehouse[]>('/warehouses');
    return response.data;
  },

  getById: async (id: number): Promise<WarehousesWarehouse> => {
    const response = await apiClient.get<WarehousesWarehouse>(`/warehouses/${id}`);
    return response.data;
  },

  create: async (data: WarehousesCreateWarehouseRequest): Promise<WarehousesWarehouse> => {
    const response = await apiClient.post<WarehousesWarehouse>('/warehouses', data);
    return response.data;
  },

  update: async (id: number, data: WarehousesUpdateWarehouseRequest): Promise<WarehousesWarehouse> => {
    const response = await apiClient.put<WarehousesWarehouse>(`/warehouses/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/warehouses/${id}`);
  }
};
