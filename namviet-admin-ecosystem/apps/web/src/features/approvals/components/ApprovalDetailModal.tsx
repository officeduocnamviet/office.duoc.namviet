import React, { useState } from 'react';
import { Modal, Steps, Button, Tag, Input, Form, Divider, Descriptions } from 'antd';
import { ApprovalRequest, ApprovalStep } from '../api/approvalApi';
import { useApprovalSteps, useUpdateApprovalRequest, useUpdateApprovalStep } from '../hooks/useApprovals';
import { CheckCircle, XCircle, Clock, AlertCircle } from 'lucide-react';
import dayjs from 'dayjs';
import { toast } from 'sonner';

interface Props {
  request: ApprovalRequest | null;
  isOpen: boolean;
  onClose: () => void;
}

export const ApprovalDetailModal: React.FC<Props> = ({ request, isOpen, onClose }) => {
  const { data: steps, isLoading } = useApprovalSteps(request?.id);
  const updateReqMutation = useUpdateApprovalRequest();
  const updateStepMutation = useUpdateApprovalStep();
  const [comments, setComments] = useState('');

  if (!request) return null;

  const handleAction = (status: 'APPROVED' | 'REJECTED') => {
    // Example: find the current pending step to update
    const currentStepObj = steps?.find(s => s.status === 'PENDING');
    
    if (currentStepObj && currentStepObj.id) {
      updateStepMutation.mutate(
        { stepId: currentStepObj.id, status, comments },
        {
          onSuccess: () => {
            toast.success(`Đã ${status === 'APPROVED' ? 'duyệt' : 'từ chối'} yêu cầu`);
            onClose();
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    } else {
      // Direct request update if no steps found or bypassing steps
      updateReqMutation.mutate(
        { id: request.id!, status },
        {
          onSuccess: () => {
            toast.success(`Đã ${status === 'APPROVED' ? 'duyệt' : 'từ chối'} yêu cầu`);
            onClose();
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    }
  };

  const getStatusColor = (status?: string) => {
    switch (status) {
      case 'APPROVED': return 'success';
      case 'REJECTED': return 'error';
      case 'PENDING': return 'processing';
      default: return 'default';
    }
  };

  const stepItems = steps?.sort((a, b) => (a.step_order || 0) - (b.step_order || 0)).map((step) => ({
    title: `Bước ${step.step_order}`,
    description: (
      <div className="mt-2 text-xs text-slate-500">
        <div><strong className="text-slate-700">Người duyệt:</strong> {step.approver_id || step.approver_role || 'Chưa phân công'}</div>
        <div><strong className="text-slate-700">Trạng thái:</strong> <Tag color={getStatusColor(step.status)} className="mt-1">{step.status}</Tag></div>
        {step.comments && <div><strong className="text-slate-700">Ghi chú:</strong> {step.comments}</div>}
        {step.action_at && <div><strong className="text-slate-700">Thời gian:</strong> {dayjs(step.action_at).format('DD/MM/YYYY HH:mm')}</div>}
      </div>
    ),
    status: step.status === 'APPROVED' ? 'finish' : step.status === 'REJECTED' ? 'error' : step.status === 'PENDING' ? 'process' : 'wait' as any,
  })) || [];

  return (
    <Modal
      title={
        <div className="flex items-center gap-2">
          Chi tiết Phê duyệt
          <Tag color={getStatusColor(request.status)}>{request.status}</Tag>
        </div>
      }
      open={isOpen}
      onCancel={onClose}
      footer={null}
      width={700}
      destroyOnClose
    >
      <div className="mt-4 space-y-6">
        <Descriptions bordered size="small" column={2}>
          <Descriptions.Item label="Mã Yêu cầu">{request.id}</Descriptions.Item>
          <Descriptions.Item label="Loại">{request.request_type}</Descriptions.Item>
          <Descriptions.Item label="Mã tham chiếu">{request.ref_id || '---'}</Descriptions.Item>
          <Descriptions.Item label="Người tạo">{request.requester_id || '---'}</Descriptions.Item>
          <Descriptions.Item label="Ngày tạo" span={2}>
            {request.created_at ? dayjs(request.created_at).format('DD/MM/YYYY HH:mm') : '---'}
          </Descriptions.Item>
          <Descriptions.Item label="Dữ liệu (Payload)" span={2}>
            <div className="max-h-32 overflow-y-auto bg-slate-50 p-2 rounded text-xs font-mono">
              {typeof request.payload === 'object' ? JSON.stringify(request.payload, null, 2) : String(request.payload || '---')}
            </div>
          </Descriptions.Item>
        </Descriptions>

        <div>
          <h3 className="text-sm font-bold text-slate-800 mb-4">Tiến trình Duyệt</h3>
          {isLoading ? (
            <div className="text-center py-4 text-slate-400">Đang tải tiến trình...</div>
          ) : stepItems.length > 0 ? (
            <Steps
              direction="vertical"
              current={steps?.findIndex(s => s.status === 'PENDING') ?? 0}
              items={stepItems}
              size="small"
            />
          ) : (
            <div className="text-slate-500 text-sm italic">Không có chi tiết các bước duyệt.</div>
          )}
        </div>

        {request.status === 'PENDING' && (
          <>
            <Divider />
            <div>
              <h3 className="text-sm font-bold text-slate-800 mb-2">Hành động của bạn</h3>
              <Input.TextArea 
                rows={3} 
                placeholder="Nhập ghi chú phản hồi (nếu có)..." 
                value={comments}
                onChange={e => setComments(e.target.value)}
                className="mb-4"
              />
              <div className="flex justify-end gap-3">
                <Button 
                  danger 
                  icon={<XCircle className="w-4 h-4" />}
                  onClick={() => handleAction('REJECTED')}
                  loading={updateReqMutation.isPending || updateStepMutation.isPending}
                >
                  Từ chối
                </Button>
                <Button 
                  type="primary" 
                  className="bg-emerald-600 hover:bg-emerald-700" 
                  icon={<CheckCircle className="w-4 h-4" />}
                  onClick={() => handleAction('APPROVED')}
                  loading={updateReqMutation.isPending || updateStepMutation.isPending}
                >
                  Phê duyệt
                </Button>
              </div>
            </div>
          </>
        )}
      </div>
    </Modal>
  );
};
