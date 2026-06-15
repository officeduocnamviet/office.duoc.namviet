import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { inventoryApi } from '../api';
import { InventoryCreateTransactionRequest } from '@namviet/shared-types/src/backend.d';

export const inventoryKeys = {
  transactions: ['inventory', 'transactions'] as const,
  transactionDetail: (id: number) => [...inventoryKeys.transactions, id] as const,
  levels: ['inventory', 'levels'] as const,
};

export const useTransactions = (warehouseId?: number, type?: string) => {
  return useQuery({
    queryKey: [...inventoryKeys.transactions, { warehouseId, type }],
    queryFn: () => inventoryApi.getTransactions({ warehouse_id: warehouseId, type }),
  });
};

export const useTransaction = (id: number) => {
  return useQuery({
    queryKey: inventoryKeys.transactionDetail(id),
    queryFn: () => inventoryApi.getTransactionById(id),
    enabled: !!id,
  });
};

export const useCreateTransaction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: InventoryCreateTransactionRequest) => inventoryApi.createTransaction(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: inventoryKeys.transactions });
      queryClient.invalidateQueries({ queryKey: inventoryKeys.levels });
    },
  });
};

export const useInventoryLevels = (warehouseId?: number) => {
  return useQuery({
    queryKey: [...inventoryKeys.levels, { warehouseId }],
    queryFn: () => inventoryApi.getInventoryLevels(warehouseId),
  });
};
