import { CommonBoard } from '@/features/dashboard/components/CommonBoard';
import { RoleBasedBoard } from '@/features/dashboard/components/RoleBasedBoard';

export default function DashboardPage() {
  return (
    <div className="flex flex-col gap-6">
      <CommonBoard />
      <RoleBasedBoard />
    </div>
  );
}
