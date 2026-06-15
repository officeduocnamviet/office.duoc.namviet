import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { financeApi, CreateFundAccountRequest, UpdateFundAccountRequest } from '../api';
import { toast } from 'sonner';

export const useFundAccounts = () => {
  return useQuery({
    queryKey: ['fundAccounts'],
    queryFn: financeApi.getFundAccounts,
  });
};

export const useCreateFundAccount = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: CreateFundAccountRequest) => financeApi.createFundAccount(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fundAccounts'] });
      toast.success('Thêm tài khoản thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateFundAccount = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: UpdateFundAccountRequest }) => 
      financeApi.updateFundAccount(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fundAccounts'] });
      toast.success('Cập nhật tài khoản thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteFundAccount = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => financeApi.deleteFundAccount(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['fundAccounts'] });
      toast.success('Xóa tài khoản thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useFinanceTransactions = () => {
  return useQuery({
    queryKey: ['financeTransactions'],
    queryFn: financeApi.getFinanceTransactions,
  });
};

export const useCreateFinanceTransaction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: any) => financeApi.createFinanceTransaction(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['financeTransactions'] });
      toast.success('Ghi nhận giao dịch thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useUpdateFinanceTransaction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }: { id: string; data: any }) => 
      financeApi.updateFinanceTransaction(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['financeTransactions'] });
      toast.success('Cập nhật giao dịch thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useDeleteFinanceTransaction = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => financeApi.deleteFinanceTransaction(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['financeTransactions'] });
      toast.success('Xóa giao dịch thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useChartOfAccounts = () => {
  return useQuery({
    queryKey: ['chartOfAccounts'],
    queryFn: financeApi.getChartOfAccounts,
  });
};

export const useCreateChartOfAccount = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: any) => financeApi.createChartOfAccount(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['chartOfAccounts'] });
      toast.success('Tạo tài khoản kế toán thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};

export const useAccountingJournals = () => {
  return useQuery({
    queryKey: ['accountingJournals'],
    queryFn: financeApi.getAccountingJournals,
  });
};

export const useCreateAccountingJournal = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (data: any) => financeApi.createAccountingJournal(data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['accountingJournals'] });
      toast.success('Ghi sổ thành công');
    },
    onError: (error: any) => {
      toast.error(error?.response?.data?.error || 'Có lỗi xảy ra');
    },
  });
};
