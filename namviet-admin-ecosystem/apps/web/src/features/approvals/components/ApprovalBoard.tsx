import React, { useState } from 'react';
import { Tabs, Tag, Input, Button } from 'antd';
import { useApprovalRequests } from '../hooks/useApprovals';
import { ApprovalRequest } from '../api/approvalApi';
import { ApprovalDetailModal } from './ApprovalDetailModal';
import { Clock, CheckCircle, XCircle, FileText, Search } from 'lucide-react';
import dayjs from 'dayjs';

export const ApprovalBoard = () => {
  const { data: requests, isLoading } = useApprovalRequests();
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedReq, setSelectedReq] = useState<ApprovalRequest | null>(null);

  const filteredRequests = requests?.filter(req => 
    req.id?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    req.request_type?.toLowerCase().includes(searchTerm.toLowerCase())
  ) || [];

  const pendingRequests = filteredRequests.filter(req => req.status === 'PENDING');
  const approvedRequests = filteredRequests.filter(req => req.status === 'APPROVED');
  const rejectedRequests = filteredRequests.filter(req => req.status === 'REJECTED');

  const renderCardList = (list: ApprovalRequest[]) => {
    if (isLoading) return <div className="p-4 text-center text-slate-500">Đang tải dữ liệu...</div>;
    if (list.length === 0) return <div className="p-8 text-center text-slate-400 border border-dashed rounded-xl border-slate-200">Không có yêu cầu nào</div>;

    return (
      <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        {list.map(req => (
          <div 
            key={req.id} 
            className="bg-white p-4 rounded-xl border border-slate-200 shadow-sm hover:shadow-md transition-shadow cursor-pointer flex flex-col gap-3"
            onClick={() => setSelectedReq(req)}
          >
            <div className="flex justify-between items-start">
              <div>
                <Tag color="blue" className="mb-2">{req.request_type}</Tag>
                <h3 className="font-bold text-slate-800 text-sm line-clamp-1">{req.id}</h3>
              </div>
              {req.status === 'PENDING' && <Clock className="w-5 h-5 text-orange-500" />}
              {req.status === 'APPROVED' && <CheckCircle className="w-5 h-5 text-emerald-500" />}
              {req.status === 'REJECTED' && <XCircle className="w-5 h-5 text-red-500" />}
            </div>

            <div className="text-xs text-slate-500 space-y-1">
              <div className="flex justify-between">
                <span>Người tạo:</span>
                <span className="font-medium text-slate-700">{req.requester_id || '---'}</span>
              </div>
              <div className="flex justify-between">
                <span>Mã tham chiếu:</span>
                <span className="font-medium text-slate-700">{req.ref_id || '---'}</span>
              </div>
              <div className="flex justify-between">
                <span>Ngày tạo:</span>
                <span>{req.created_at ? dayjs(req.created_at).format('DD/MM/YYYY HH:mm') : '---'}</span>
              </div>
            </div>
            
            <div className="mt-2 pt-3 border-t border-slate-100 flex justify-between items-center">
              <span className="text-xs text-slate-400 flex items-center gap-1">
                <FileText className="w-3 h-3" />
                Bước hiện tại: {req.current_step || 1}
              </span>
              <Button type="link" size="small" className="p-0 h-auto text-[13px]">Xem chi tiết</Button>
            </div>
          </div>
        ))}
      </div>
    );
  };

  const tabItems = [
    {
      key: 'pending',
      label: <span className="flex items-center gap-2">Chờ duyệt <span className="bg-orange-100 text-orange-600 py-0.5 px-2 rounded-full text-xs font-bold">{pendingRequests.length}</span></span>,
      children: renderCardList(pendingRequests)
    },
    {
      key: 'approved',
      label: <span className="flex items-center gap-2">Đã duyệt <span className="bg-emerald-100 text-emerald-600 py-0.5 px-2 rounded-full text-xs font-bold">{approvedRequests.length}</span></span>,
      children: renderCardList(approvedRequests)
    },
    {
      key: 'rejected',
      label: <span className="flex items-center gap-2">Đã từ chối <span className="bg-red-100 text-red-600 py-0.5 px-2 rounded-full text-xs font-bold">{rejectedRequests.length}</span></span>,
      children: renderCardList(rejectedRequests)
    }
  ];

  return (
    <div className="space-y-4">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 bg-white p-4 rounded-xl shadow-sm border border-slate-100">
        <Input
          placeholder="Tìm kiếm mã yêu cầu hoặc loại..."
          prefix={<Search className="w-4 h-4 text-slate-400" />}
          value={searchTerm}
          onChange={e => setSearchTerm(e.target.value)}
          className="w-full md:w-80"
        />
        {/* Further filters can be added here */}
      </div>

      <Tabs 
        defaultActiveKey="pending" 
        items={tabItems} 
        className="bg-white p-4 rounded-xl shadow-sm border border-slate-100"
      />

      <ApprovalDetailModal 
        isOpen={!!selectedReq} 
        request={selectedReq} 
        onClose={() => setSelectedReq(null)} 
      />
    </div>
  );
};
