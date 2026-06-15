"use client";

import React from 'react';
import { ShiftCalendar } from '@/features/hr/components/ShiftCalendar';
import { CalendarDays } from 'lucide-react';

export default function ShiftsPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <CalendarDays className="w-6 h-6 text-blue-600" />
          Lịch Phân Ca (Work Shifts)
        </h1>
        <p className="text-slate-500">
          Quản lý lịch làm việc, phân ca cho nhân viên theo từng ngày.
        </p>
      </div>

      <ShiftCalendar />
    </div>
  );
}
