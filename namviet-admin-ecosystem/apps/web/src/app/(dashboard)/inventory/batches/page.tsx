'use client';

import { useEffect } from 'react';
import { BatchTable } from '@/features/batches/components/BatchTable';

export default function BatchesPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'inventory_batches');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-hidden">
      <BatchTable />
    </div>
  );
}
