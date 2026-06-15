import { useQuery } from '@tanstack/react-query';
import { integrationApi } from '../api/integrationApi';

export const integrationKeys = {
  connections: ['connections'] as const,
  webhooks: ['webhooks'] as const,
  shipping: ['shippingPartners'] as const,
};

export const useConnections = () => {
  return useQuery({
    queryKey: integrationKeys.connections,
    queryFn: integrationApi.getConnections,
  });
};

export const useWebhookLogs = () => {
  return useQuery({
    queryKey: integrationKeys.webhooks,
    queryFn: integrationApi.getWebhookLogs,
  });
};

export const useShippingPartners = () => {
  return useQuery({
    queryKey: integrationKeys.shipping,
    queryFn: integrationApi.getShippingPartners,
  });
};
