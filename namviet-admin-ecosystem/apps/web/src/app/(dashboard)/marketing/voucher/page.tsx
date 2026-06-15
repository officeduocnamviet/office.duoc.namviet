'use client';

import { useEffect } from 'react';
import { PromotionTable } from '@/features/promotions/components/PromotionTable';

export default function VoucherPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'marketing_voucher');
  }, []);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <PromotionTable />
    </div>
  );
}
