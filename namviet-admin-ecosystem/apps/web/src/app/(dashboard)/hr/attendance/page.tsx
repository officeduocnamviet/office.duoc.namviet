"use client";

import React from 'react';
import { TimeAttendanceTable } from '@/features/hr/components/TimeAttendanceTable';

import { CheckInOutPanel } from '@/features/hr/components/CheckInOutPanel';

export default function TimeAttendancesPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <CheckInOutPanel />
      <TimeAttendanceTable />
    </div>
  );
}
