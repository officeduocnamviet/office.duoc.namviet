import { apiClient } from '@/lib/axios';

export interface ThirdPartyConnection {
  id?: string;
  partner_code?: string;
  connection_type?: string;
  credentials?: any;
  status?: string;
  created_at?: string;
}

export interface WebhookLog {
  id?: string;
  connection_id?: string;
  event_type?: string;
  payload?: any;
  response_status?: number;
  response_body?: string;
  created_at?: string;
}

export interface ShippingPartner {
  id?: string;
  code?: string;
  name?: string;
  partner_type?: string;
  status?: string;
  tracking_url_template?: string;
}

export const integrationApi = {
  getConnections: async (): Promise<ThirdPartyConnection[]> => {
    const response = await apiClient.get<ThirdPartyConnection[]>('/connections');
    return response.data;
  },
  getWebhookLogs: async (): Promise<WebhookLog[]> => {
    const response = await apiClient.get<WebhookLog[]>('/webhook-logs');
    return response.data;
  },
  getShippingPartners: async (): Promise<ShippingPartner[]> => {
    const response = await apiClient.get<ShippingPartner[]>('/shipping-partners');
    return response.data;
  }
};
