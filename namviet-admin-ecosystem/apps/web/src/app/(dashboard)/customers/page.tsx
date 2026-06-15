'use client';

import { useEffect } from 'react';
import { CustomerTable } from '@/features/customers/components/CustomerTable';

export default function CustomersPage() {
  useEffect(() => {
    // Để active menu (mặc dù menu có id='customers')
    sessionStorage.setItem('activePage', 'customers');
  }, []);

  return (
    <div className="animate-in fade-in slide-in-from-bottom-4 duration-500">
      <CustomerTable />
    </div>
  );
}
