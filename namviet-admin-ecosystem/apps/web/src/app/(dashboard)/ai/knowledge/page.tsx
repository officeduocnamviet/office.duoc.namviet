"use client";

import React from 'react';
import { KnowledgeVectorTable } from '@/features/ai/components/KnowledgeVectorTable';
import { BrainCircuit } from 'lucide-react';

export default function KnowledgeBasePage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <BrainCircuit className="w-6 h-6 text-indigo-600" />
          Quản trị Tri thức AI (Knowledge Vectors)
        </h1>
        <p className="text-slate-500">
          Quản lý cơ sở dữ liệu Vector phục vụ cho Chatbot AI trả lời tự động.
        </p>
      </div>

      <KnowledgeVectorTable />
    </div>
  );
}
