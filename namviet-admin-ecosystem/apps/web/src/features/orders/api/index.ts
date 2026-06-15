import { apiClient } from '@/lib/axios';
import { 
  OrdersOrder, 
  OrdersCreateOrderRequest, 
  OrdersUpdateOrderRequest 
} from '@namviet/shared-types/src/backend.d';

export const orderApi = {
  getAll: async (type?: string): Promise<OrdersOrder[]> => {
    const params = type ? { order_type: type } : {};
    const { data } = await apiClient.get<OrdersOrder[]>('/api/orders', { params });
    return data;
  },

  getById: async (id: string): Promise<OrdersOrder> => {
    const { data } = await apiClient.get<OrdersOrder>(`/api/orders/${id}`);
    return data;
  },

  create: async (request: OrdersCreateOrderRequest): Promise<OrdersOrder> => {
    const { data } = await apiClient.post<OrdersOrder>('/api/orders', request);
    return data;
  },

  updateStatus: async ({ id, data }: { id: string; data: OrdersUpdateOrderRequest }): Promise<OrdersOrder> => {
    const { data: resData } = await apiClient.put<OrdersOrder>(`/api/orders/${id}`, data);
    return resData;
  },
};
