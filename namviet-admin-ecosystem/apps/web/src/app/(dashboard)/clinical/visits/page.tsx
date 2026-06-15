"use client";

import React from 'react';
import { MedicalVisitTable } from '@/features/clinical/components/MedicalVisitTable';

export default function MedicalVisitsPage() {
  return (
    <div className="h-full max-h-full overflow-hidden">
      <MedicalVisitTable />
    </div>
  );
}
