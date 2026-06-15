'use client';

import { useEffect } from 'react';
import { RetailOrderForm } from '@/features/orders/components/RetailOrderForm';

export default function RetailCreatePage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'retail_create');
  }, []);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-gray-800">Tạo Đơn Hàng Bán Lẻ</h1>
        <p className="text-gray-500">POS Bán hàng nhanh tại quầy</p>
      </div>
      <RetailOrderForm />
    </div>
  );
}
