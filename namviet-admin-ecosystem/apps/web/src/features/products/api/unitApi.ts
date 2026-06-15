import { apiClient } from '@/lib/axios';
import { 
  ProductUnitsProductUnit, 
  ProductUnitsCreateProductUnitRequest, 
  ProductUnitsUpdateProductUnitRequest 
} from '@namviet/shared-types/src/backend.d';

export type ProductUnit = ProductUnitsProductUnit;
export type CreateProductUnitRequest = ProductUnitsCreateProductUnitRequest;
export type UpdateProductUnitRequest = ProductUnitsUpdateProductUnitRequest;

export const productUnitApi = {
  getByProductId: async (productId: number): Promise<ProductUnit[]> => {
    const response = await apiClient.get<ProductUnit[]>(`/products/${productId}/units`);
    return response.data;
  },

  create: async (productId: number, data: CreateProductUnitRequest): Promise<ProductUnit> => {
    const response = await apiClient.post<ProductUnit>(`/products/${productId}/units`, data);
    return response.data;
  },

  update: async (productId: number, unitId: number, data: UpdateProductUnitRequest): Promise<ProductUnit> => {
    const response = await apiClient.put<ProductUnit>(`/products/${productId}/units/${unitId}`, data);
    return response.data;
  },

  delete: async (productId: number, unitId: number): Promise<void> => {
    await apiClient.delete(`/products/${productId}/units/${unitId}`);
  }
};
