import { apiClient } from '@/lib/axios';
import { 
  RolesRole, 
  RolesCreateRoleRequest, 
  RolesUpdateRoleRequest 
} from '@namviet/shared-types/src/backend.d';

export const roleApi = {
  getAll: async (): Promise<RolesRole[]> => {
    const response = await apiClient.get<RolesRole[]>('/roles');
    return response.data;
  },

  getById: async (id: string): Promise<RolesRole> => {
    const response = await apiClient.get<RolesRole>(`/roles/${id}`);
    return response.data;
  },

  create: async (data: RolesCreateRoleRequest): Promise<RolesRole> => {
    const response = await apiClient.post<RolesRole>('/roles', data);
    return response.data;
  },

  update: async (id: string, data: RolesUpdateRoleRequest): Promise<RolesRole> => {
    const response = await apiClient.put<RolesRole>(`/roles/${id}`, data);
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await apiClient.delete(`/roles/${id}`);
  }
};
