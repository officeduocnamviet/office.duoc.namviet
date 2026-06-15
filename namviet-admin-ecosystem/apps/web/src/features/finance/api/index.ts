import { apiClient as api } from '@/lib/axios';
import {
  FundAccountsFundAccount as FundAccount,
  FundAccountsCreateFundAccountRequest as CreateFundAccountRequest,
  FundAccountsUpdateFundAccountRequest as UpdateFundAccountRequest,
  FinanceTransactionsFinanceTransaction as FinanceTransaction,
  FinanceTransactionsCreateFinanceTransactionRequest as CreateFinanceTransactionRequest,
  FinanceTransactionsUpdateFinanceTransactionRequest as UpdateFinanceTransactionRequest,
  ChartOfAccountsChartOfAccount as ChartOfAccount,
  ChartOfAccountsCreateChartOfAccountRequest as CreateChartOfAccountRequest,
  AccountingJournalsAccountingJournal as AccountingJournal,
  AccountingJournalsCreateAccountingJournalRequest as CreateAccountingJournalRequest
} from '@namviet/shared-types/src/backend.d';

export type {
  FundAccount,
  CreateFundAccountRequest,
  UpdateFundAccountRequest,
  FinanceTransaction,
  CreateFinanceTransactionRequest,
  UpdateFinanceTransactionRequest,
  ChartOfAccount,
  CreateChartOfAccountRequest,
  AccountingJournal,
  CreateAccountingJournalRequest
};

export const financeApi = {
  // Fund Accounts
  getFundAccounts: () => api.get<FundAccount[]>('/api/fund-accounts').then(res => res.data),
  getFundAccount: (id: string) => api.get<FundAccount>(`/api/fund-accounts/${id}`).then(res => res.data),
  createFundAccount: (data: CreateFundAccountRequest) => api.post<FundAccount>('/api/fund-accounts', data).then(res => res.data),
  updateFundAccount: (id: string, data: UpdateFundAccountRequest) => api.put<FundAccount>(`/api/fund-accounts/${id}`, data).then(res => res.data),
  deleteFundAccount: (id: string) => api.delete(`/api/fund-accounts/${id}`).then(res => res.data),

  // Finance Transactions
  getFinanceTransactions: () => api.get<FinanceTransaction[]>('/api/finance-transactions').then(res => res.data),
  getFinanceTransaction: (id: string) => api.get<FinanceTransaction>(`/api/finance-transactions/${id}`).then(res => res.data),
  createFinanceTransaction: (data: CreateFinanceTransactionRequest) => api.post<FinanceTransaction>('/api/finance-transactions', data).then(res => res.data),
  updateFinanceTransaction: (id: string, data: UpdateFinanceTransactionRequest) => api.put<FinanceTransaction>(`/api/finance-transactions/${id}`, data).then(res => res.data),
  deleteFinanceTransaction: (id: string) => api.delete(`/api/finance-transactions/${id}`).then(res => res.data),

  // COA
  getChartOfAccounts: () => api.get<ChartOfAccount[]>('/api/chart-of-accounts').then(res => res.data),
  createChartOfAccount: (data: CreateChartOfAccountRequest) => api.post<ChartOfAccount>('/api/chart-of-accounts', data).then(res => res.data),

  // Journals
  getAccountingJournals: () => api.get<AccountingJournal[]>('/api/accounting-journals').then(res => res.data),
  createAccountingJournal: (data: CreateAccountingJournalRequest) => api.post<AccountingJournal>('/api/accounting-journals', data).then(res => res.data),
};
