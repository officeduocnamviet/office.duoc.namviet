"use client";

import React from 'react';
import { ChartOfAccountTable } from '@/features/finance/components/ChartOfAccountTable';

export default function ChartOfAccountsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <ChartOfAccountTable />
    </div>
  );
}
