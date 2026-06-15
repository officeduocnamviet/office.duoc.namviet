"use client";

import React, { useEffect } from 'react';
import { WarehouseTable } from '@/features/warehouses/components/WarehouseTable';
import { toast } from 'sonner';

export default function WarehousesPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'warehouses');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-hidden">
      <WarehouseTable />
    </div>
  );
}
