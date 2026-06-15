'use client';

import { useEffect } from 'react';
import { OrderTable } from '@/features/orders/components/OrderTable';

export default function RetailListPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'retail_list');
  }, []);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <OrderTable orderType="RETAIL" />
    </div>
  );
}
