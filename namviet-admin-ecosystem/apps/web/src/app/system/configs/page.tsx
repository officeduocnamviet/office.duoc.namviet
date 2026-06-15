"use client";

import React from 'react';
import { ConfigTable } from '@/features/system/components/ConfigTable';
import { Settings } from 'lucide-react';

export default function SystemConfigsPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <Settings className="w-6 h-6 text-blue-600" />
          Cấu hình Hệ thống (System Configs)
        </h1>
        <p className="text-slate-500">
          Quản lý các biến môi trường, API keys bên thứ ba, cấu hình tính năng chung cho toàn bộ hệ thống.
        </p>
      </div>

      <ConfigTable />
    </div>
  );
}
