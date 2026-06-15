"use client";

import React from 'react';
import { AccountingJournalTable } from '@/features/finance/components/AccountingJournalTable';

export default function AccountingJournalsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <AccountingJournalTable />
    </div>
  );
}
