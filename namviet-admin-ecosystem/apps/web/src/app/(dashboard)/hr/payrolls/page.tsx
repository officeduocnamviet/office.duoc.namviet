"use client";

import React from 'react';
import { PayrollTable } from '@/features/hr/components/PayrollTable';

export default function PayrollsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <PayrollTable />
    </div>
  );
}
