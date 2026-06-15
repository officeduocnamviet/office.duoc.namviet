import { apiClient } from '@/lib/axios';
import { 
  CustomersCustomer, 
  CustomersCreateCustomerRequest, 
  CustomersUpdateCustomerRequest,
  CustomerRecordsCustomerVaccinationRecord,
  CustomerRecordsCustomerVoucher
} from '@namviet/shared-types/src/backend.d';

export const customerApi = {
  getAll: async (): Promise<CustomersCustomer[]> => {
    const { data } = await apiClient.get<CustomersCustomer[]>('/api/customers');
    return data;
  },

  getById: async (id: number): Promise<CustomersCustomer> => {
    const { data } = await apiClient.get<CustomersCustomer>(`/api/customers/${id}`);
    return data;
  },

  create: async (request: CustomersCreateCustomerRequest): Promise<CustomersCustomer> => {
    const { data } = await apiClient.post<CustomersCustomer>('/api/customers', request);
    return data;
  },

  update: async ({ id, data }: { id: number; data: CustomersUpdateCustomerRequest }): Promise<CustomersCustomer> => {
    const { data: resData } = await apiClient.put<CustomersCustomer>(`/api/customers/${id}`, data);
    return resData;
  },

  delete: async (id: number): Promise<void> => {
    await apiClient.delete(`/api/customers/${id}`);
  },

  getVaccinations: async (customerId: number): Promise<CustomerRecordsCustomerVaccinationRecord[]> => {
    const { data } = await apiClient.get<CustomerRecordsCustomerVaccinationRecord[]>(`/api/customer-vaccinations?customer_id=${customerId}`);
    return data;
  },

  getVouchers: async (customerId: number): Promise<CustomerRecordsCustomerVoucher[]> => {
    const { data } = await apiClient.get<CustomerRecordsCustomerVoucher[]>(`/api/customer-vouchers?customer_id=${customerId}`);
    return data;
  }
};
