import { apiClient } from '@/lib/axios';
import { 
  CategoriesCategory, 
  CategoriesCreateCategoryRequest, 
  CategoriesUpdateCategoryRequest 
} from '@namviet/shared-types/src/backend.d';

export const categoryApi = {
  getAll: async (): Promise<CategoriesCategory[]> => {
    const response = await apiClient.get<CategoriesCategory[]>('/categories');
    return response.data;
  },

  getById: async (id: number): Promise<CategoriesCategory> => {
    const response = await apiClient.get<CategoriesCategory>(`/categories/${id}`);
    return response.data;
  },

  create: async (data: CategoriesCreateCategoryRequest): Promise<CategoriesCategory> => {
    const response = await apiClient.post<CategoriesCategory>('/categories', data);
    return response.data;
  },

  update: async (id: number, data: CategoriesUpdateCategoryRequest): Promise<CategoriesCategory> => {
    const response = await apiClient.put<CategoriesCategory>(`/categories/${id}`, data);
    return response.data;
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/categories/${id}`);
  }
};
