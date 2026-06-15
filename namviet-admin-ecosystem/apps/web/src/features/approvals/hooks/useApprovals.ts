import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { approvalApi } from '../api/approvalApi';

export const approvalKeys = {
  allRequests: ['approvalRequests'] as const,
  steps: (requestId: string) => ['approvalSteps', requestId] as const,
};

export const useApprovalRequests = () => {
  return useQuery({
    queryKey: approvalKeys.allRequests,
    queryFn: approvalApi.getAllRequests,
  });
};

export const useApprovalSteps = (requestId: string | undefined) => {
  return useQuery({
    queryKey: approvalKeys.steps(requestId || ''),
    queryFn: () => approvalApi.getStepsByRequestId(requestId!),
    enabled: !!requestId,
  });
};

export const useUpdateApprovalRequest = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status, payload }: { id: string; status: string; payload?: any }) => 
      approvalApi.updateRequestStatus(id, status, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: approvalKeys.allRequests });
    },
  });
};

export const useUpdateApprovalStep = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ stepId, status, comments }: { stepId: string; status: string; comments?: string }) => 
      approvalApi.updateStepStatus(stepId, status, comments),
    onSuccess: (_, variables) => {
      // Invalidate relevant queries (would need requestId to be exact, but we can invalidate all requests to be safe)
      queryClient.invalidateQueries({ queryKey: approvalKeys.allRequests });
    },
  });
};
