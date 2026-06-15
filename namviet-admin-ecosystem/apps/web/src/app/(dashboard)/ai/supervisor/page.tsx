"use client";

import React from 'react';
import { SupervisorTable } from '@/features/ai/components/SupervisorTable';
import { ShieldAlert } from 'lucide-react';

export default function SupervisorPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <ShieldAlert className="w-6 h-6 text-indigo-600" />
          Giám sát Chatbot (Bot Supervisor)
        </h1>
        <p className="text-slate-500">
          Theo dõi các phiên tư vấn của AI, cảnh báo cảm xúc tiêu cực và hỗ trợ Human-Handoff.
        </p>
      </div>

      <SupervisorTable />
    </div>
  );
}
