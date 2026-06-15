import { apiClient } from '@/lib/axios';
import { 
  InventoryInventoryTransaction, 
  InventoryCreateTransactionRequest,
  InventoryInventoryBatch
} from '@namviet/shared-types/src/backend.d';

export const inventoryApi = {
  getTransactions: async (params?: { warehouse_id?: number, type?: string }): Promise<InventoryInventoryTransaction[]> => {
    const response = await apiClient.get<InventoryInventoryTransaction[]>('/inventory/transactions', { params });
    return response.data;
  },

  getTransactionById: async (id: number): Promise<InventoryInventoryTransaction> => {
    const response = await apiClient.get<InventoryInventoryTransaction>(`/inventory/transactions/${id}`);
    return response.data;
  },

  createTransaction: async (data: InventoryCreateTransactionRequest): Promise<InventoryInventoryTransaction> => {
    const response = await apiClient.post<InventoryInventoryTransaction>('/inventory/transactions', data);
    return response.data;
  },

  getInventoryLevels: async (warehouseId?: number): Promise<InventoryInventoryBatch[]> => {
    const response = await apiClient.get<InventoryInventoryBatch[]>('/inventory/levels', { 
      params: { warehouse_id: warehouseId } 
    });
    return response.data;
  }
};
