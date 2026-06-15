"use client";

import React from 'react';
import { AppointmentTable } from '@/features/clinical/components/AppointmentTable';

export default function AppointmentsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <AppointmentTable />
    </div>
  );
}
