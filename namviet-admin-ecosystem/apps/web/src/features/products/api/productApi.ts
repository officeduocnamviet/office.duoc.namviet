import { apiClient } from '@/lib/axios';
import { 
  ProductsProduct, 
  ProductsCreateProductRequest, 
  ProductsUpdateProductRequest 
} from '@namviet/shared-types/src/backend.d';

export const productApi = {
  getAll: async (): Promise<ProductsProduct[]> => {
    const response = await apiClient.get<ProductsProduct[]>('/products');
    return response.data;
  },

  getById: async (id: number): Promise<ProductsProduct> => {
    const response = await apiClient.get<ProductsProduct>(`/products/${id}`);
    return response.data;
  },

  create: async (data: ProductsCreateProductRequest): Promise<ProductsProduct> => {
    const response = await apiClient.post<ProductsProduct>('/products', data);
    return response.data;
  },

  update: async (id: number, data: ProductsUpdateProductRequest): Promise<ProductsProduct> => {
    const response = await apiClient.put<ProductsProduct>(`/products/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/products/${id}`);
  }
};
