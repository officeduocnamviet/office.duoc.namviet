'use client';

import { useEffect } from 'react';
import { ExcelImport } from '@/features/inventory/components/ExcelImport';

export default function InventoryImportPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'inventory_import');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-y-auto bg-gray-50/30">
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-800 tracking-tight">Nhập Kho Hàng Loạt</h2>
        <p className="text-sm text-gray-500 mt-1">Sử dụng file Excel để thêm nhiều sản phẩm vào kho cùng lúc.</p>
      </div>
      <ExcelImport />
    </div>
  );
}
