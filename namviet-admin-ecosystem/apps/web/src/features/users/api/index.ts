import { apiClient } from '@/lib/axios';
import { 
  UsersUser, 
  UsersCreateUserRequest, 
  UsersUpdateUserRequest 
} from '@namviet/shared-types/src/backend.d';

export const userApi = {
  getAll: async (): Promise<UsersUser[]> => {
    const response = await apiClient.get<UsersUser[]>('/users');
    return response.data;
  },

  getById: async (id: string): Promise<UsersUser> => {
    const response = await apiClient.get<UsersUser>(`/users/${id}`);
    return response.data;
  },

  create: async (data: UsersCreateUserRequest): Promise<UsersUser> => {
    const response = await apiClient.post<UsersUser>('/users', data);
    return response.data;
  },

  update: async (id: string, data: UsersUpdateUserRequest): Promise<UsersUser> => {
    const response = await apiClient.put<UsersUser>(`/users/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/users/${id}`);
  }
};
