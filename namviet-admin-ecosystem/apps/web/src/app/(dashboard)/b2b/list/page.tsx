'use client';

import { useEffect } from 'react';
import { OrderTable } from '@/features/orders/components/OrderTable';

export default function B2bListPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'b2b_list');
  }, []);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <OrderTable orderType="B2B" />
    </div>
  );
}
