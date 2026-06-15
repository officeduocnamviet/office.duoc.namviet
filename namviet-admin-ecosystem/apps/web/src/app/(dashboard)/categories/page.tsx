'use client';

import { useEffect } from 'react';
import { CategoryTable } from '@/features/categories/components/CategoryTable';

export default function CategoriesPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'categories');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-hidden">
      <CategoryTable />
    </div>
  );
}
