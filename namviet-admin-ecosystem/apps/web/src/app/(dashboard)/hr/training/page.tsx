"use client";

import React from 'react';
import { TrainingTable } from '@/features/hr/components/TrainingTable';
import { BookOpen } from 'lucide-react';

export default function TrainingPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <BookOpen className="w-6 h-6 text-indigo-600" />
          Khóa đào tạo (Training Courses)
        </h1>
        <p className="text-slate-500">
          Quản lý các khóa đào tạo nội bộ, nâng cao kỹ năng cho nhân viên.
        </p>
      </div>

      <TrainingTable />
    </div>
  );
}
