'use client';

import { useEffect } from 'react';
import { ProductTable } from '@/features/products/components/ProductTable';

export default function ProductsPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'products');
  }, []);

  return (
    <div className="p-0 sm:p-6 h-[calc(100vh-64px)] sm:h-[calc(100vh-80px)] overflow-hidden">
      <ProductTable />
    </div>
  );
}
