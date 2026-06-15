"use client";

import React from 'react';
import { ApprovalBoard } from '@/features/approvals/components/ApprovalBoard';
import { ClipboardCheck } from 'lucide-react';

export default function ApprovalsPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <ClipboardCheck className="w-6 h-6 text-emerald-600" />
          Trung tâm Phê duyệt (Approval Center)
        </h1>
        <p className="text-slate-500">
          Quản lý các yêu cầu cần phê duyệt (Nghỉ phép, Chi tiêu, Bàn giao ca, Nhập hàng...).
        </p>
      </div>

      <ApprovalBoard />
    </div>
  );
}
