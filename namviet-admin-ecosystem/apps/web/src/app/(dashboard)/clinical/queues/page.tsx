"use client";

import React from 'react';
import { ClinicalQueueTable } from '@/features/clinical/components/ClinicalQueueTable';

export default function ClinicalQueuesPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <ClinicalQueueTable />
    </div>
  );
}
