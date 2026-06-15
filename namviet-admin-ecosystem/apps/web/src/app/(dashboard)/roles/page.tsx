'use client';

import { useEffect } from 'react';
import { RoleTable } from '@/features/roles/components/RoleTable';

export default function RolesPage() {
  useEffect(() => {
    sessionStorage.setItem('activePage', 'roles');
  }, []);

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-hidden">
      <RoleTable />
    </div>
  );
}
