import { apiClient } from '@/lib/axios';

export interface ApprovalRequest {
  id?: string;
  ref_id?: string;
  request_type?: string;
  requester_id?: string;
  status?: string; // PENDING, APPROVED, REJECTED
  current_step?: number;
  payload?: any;
  created_at?: string;
  updated_at?: string;
}

export interface ApprovalStep {
  id?: string;
  request_id?: string;
  step_order?: number;
  approver_id?: string;
  approver_role?: string;
  status?: string; // PENDING, APPROVED, REJECTED
  comments?: string;
  action_at?: string;
  created_at?: string;
}

export const approvalApi = {
  getAllRequests: async (): Promise<ApprovalRequest[]> => {
    const response = await apiClient.get<ApprovalRequest[]>('/approval-requests');
    return response.data;
  },

  getStepsByRequestId: async (id: string): Promise<ApprovalStep[]> => {
    const response = await apiClient.get<ApprovalStep[]>(`/approval-requests/${id}/steps`);
    return response.data;
  },

  updateRequestStatus: async (id: string, status: string, payload?: any): Promise<ApprovalRequest> => {
    const response = await apiClient.put<ApprovalRequest>(`/approval-requests/${id}`, {
      status,
      payload
    });
    return response.data;
  },

  updateStepStatus: async (stepId: string, status: string, comments?: string): Promise<ApprovalStep> => {
    const response = await apiClient.put<ApprovalStep>(`/approval-steps/${stepId}`, {
      status,
      comments
    });
    return response.data;
  }
};
