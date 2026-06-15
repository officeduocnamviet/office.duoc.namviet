"use client";

import React from 'react';
import { ContractTable } from '@/features/hr/components/ContractTable';
import { FileSignature } from 'lucide-react';

export default function ContractsPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <FileSignature className="w-6 h-6 text-blue-600" />
          Hợp đồng Lao động (Contracts)
        </h1>
        <p className="text-slate-500">
          Quản lý hợp đồng lao động, thời hạn thử việc và lương cơ bản của nhân viên.
        </p>
      </div>

      <ContractTable />
    </div>
  );
}
