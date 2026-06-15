"use client";

import React from 'react';
import { EmployeeTable } from '@/features/hr/components/EmployeeTable';

export default function EmployeesPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <EmployeeTable />
    </div>
  );
}
