'use client';

import { useEffect } from 'react';
import { ManufacturerTable } from '@/features/manufacturers/components/ManufacturerTable';

export default function ManufacturersPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'manufacturers');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-hidden">
      <ManufacturerTable />
    </div>
  );
}
