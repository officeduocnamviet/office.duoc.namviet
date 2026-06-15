"use client";

import React from 'react';
import { FinanceTransactionTable } from '@/features/finance/components/FinanceTransactionTable';

export default function FinanceTransactionsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <FinanceTransactionTable />
    </div>
  );
}
