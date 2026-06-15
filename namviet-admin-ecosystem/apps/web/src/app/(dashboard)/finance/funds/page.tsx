"use client";

import React from 'react';
import { FundAccountTable } from '@/features/finance/components/FundAccountTable';

export default function FundAccountsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <FundAccountTable />
    </div>
  );
}
