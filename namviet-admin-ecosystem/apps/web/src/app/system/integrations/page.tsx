"use client";

import React from 'react';
import { IntegrationBoard } from '@/features/system/components/IntegrationBoard';
import { Link2 } from 'lucide-react';

export default function IntegrationsPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <Link2 className="w-6 h-6 text-purple-600" />
          Đối tác & Tích hợp (Integrations)
        </h1>
        <p className="text-slate-500">
          Quản lý đối tác vận chuyển (GHTK, GHN, Viettel Post...), các kết nối API nội bộ và theo dõi lịch sử Webhook.
        </p>
      </div>

      <IntegrationBoard />
    </div>
  );
}
