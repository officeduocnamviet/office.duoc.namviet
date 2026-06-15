"use client";

import React from 'react';
import { ShiftHandoverTable } from '@/features/hr/components/ShiftHandoverTable';
import { HandCoins } from 'lucide-react';

export default function ShiftHandoversPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <HandCoins className="w-6 h-6 text-green-600" />
          Bàn giao ca (Shift Handovers)
        </h1>
        <p className="text-slate-500">
          Chốt ca, đối soát tiền mặt, tiền COD và gửi yêu cầu phê duyệt doanh thu cuối ca.
        </p>
      </div>

      <ShiftHandoverTable />
    </div>
  );
}
