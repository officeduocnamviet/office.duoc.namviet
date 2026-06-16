import { apiClient } from './index';

export interface User {
  id: string;
  email?: string;
  full_name?: string;
  phone?: string;
  status: string;
  role_id?: string;
  company_id?: string;
  warehouse_id?: number;
  role_name?: string;
  company_name?: string;
  warehouse_name?: string;
}

export const userService = {
  getAll: (): Promise<User[]> => {
    return apiClient.get('/users');
  },
  create: (data: {
    email: string;
    password?: string;
    full_name: string;
    phone?: string;
    role_id: string;
    company_id: string;
    warehouse_id: number;
  }): Promise<User> => {
    return apiClient.post('/users', data);
  },
};
