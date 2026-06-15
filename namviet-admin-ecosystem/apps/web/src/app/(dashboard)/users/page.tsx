'use client';

import { useEffect } from 'react';
import { UserTable } from '@/features/users/components/UserTable';

export default function UsersPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'users');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-hidden">
      <UserTable />
    </div>
  );
}
