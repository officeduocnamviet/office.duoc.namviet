'use client';

import { useEffect } from 'react';
import { B2bOrderForm } from '@/features/orders/components/B2bOrderForm';

export default function B2bCreatePage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'b2b_create');
  }, []);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Tạo Đơn Hàng B2B (Bán Sỉ)</h1>
        <p className="text-gray-500">Kênh bán buôn / Khách hàng doanh nghiệp</p>
      </div>
      <B2bOrderForm />
    </div>
  );
}
