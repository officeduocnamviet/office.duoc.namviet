import React, { useState } from 'react';
import { ColumnDef } from '@tanstack/react-table';
import { DataTable } from '@/components/ui/DataTable';
import { Warehouse } from '@namviet/shared-types/src/warehouse.types';
import { Edit2, Trash2, MapPin, Building, Plus, Search, XCircle, CheckCircle2 } from 'lucide-react';
import { useRouter } from 'next/navigation';

interface WarehouseTableProps {
  data: Warehouse[];
  isLoading: boolean;
  onDelete: (id: number) => void;
  onToggleStatus: (id: number, currentStatus: string) => void;
}

export function WarehouseTable({ data, isLoading, onDelete, onToggleStatus }: WarehouseTableProps) {
  const router = useRouter();
  const [globalFilter, setGlobalFilter] = useState('');

  const columns: ColumnDef<Warehouse>[] = [
    {
      accessorKey: 'key',
      header: 'Mã (Key)',
      cell: ({ row }) => (
        <span className="font-bold text-slate-700 dark:text-slate-300">
          {row.original.key}
        </span>
      ),
    },
    {
      accessorKey: 'name',
      header: 'Tên Chi nhánh',
      cell: ({ row }) => (
        <div>
          <div 
            className="font-bold text-sm text-primary-600 dark:text-blue-400 hover:underline cursor-pointer"
            onClick={() => router.push(`/warehouses/edit?id=${row.original.id}`)}
          >
            {row.original.name}
          </div>
          {row.original.code && (
            <div className="text-xs text-slate-500 font-medium mt-1">
              Code: {row.original.code}
            </div>
          )}
        </div>
      ),
    },
    {
      accessorKey: 'address',
      header: 'Địa chỉ & Liên hệ',
      cell: ({ row }) => (
        <div className="text-sm">
          {row.original.address && (
            <div className="flex items-start gap-1.5 text-slate-600 dark:text-slate-400 mb-1 max-w-[200px] truncate" title={row.original.address}>
              <MapPin size={14} className="shrink-0 mt-0.5" /> {row.original.address}
            </div>
          )}
          {row.original.phone && (
            <div className="text-slate-500 text-xs">
              SĐT: {row.original.phone}
            </div>
          )}
        </div>
      ),
    },
    {
      accessorKey: 'type',
      header: 'Phân loại',
      cell: ({ row }) => (
        <span className="inline-flex items-center gap-1 px-2 py-1 rounded-md text-xs font-bold bg-slate-100 text-slate-600 border border-slate-200">
          <Building size={12} />
          {row.original.type === 'retail' ? 'Cửa hàng lẻ' : row.original.type === 'b2b' ? 'Kho sỉ (B2B)' : row.original.type}
        </span>
      ),
    },
    {
      accessorKey: 'status',
      header: 'Trạng thái',
      cell: ({ row }) => (
        <span className={`inline-flex px-2 py-1 rounded-md text-xs font-bold border ${row.original.status === 'active' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : 'bg-slate-50 text-slate-600 border-slate-200'}`}>
          {row.original.status === 'active' ? 'Hoạt động' : 'Tạm ngưng'}
        </span>
      ),
    },
    {
      id: 'actions',
      header: 'Thao tác',
      cell: ({ row }) => (
        <div className="flex items-center gap-1.5">
          <button 
            onClick={() => router.push(`/warehouses/edit?id=${row.original.id}`)}
            className="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
            title="Chỉnh sửa"
          >
            <Edit2 size={16} />
          </button>
          
          <button 
            onClick={() => {
              if (confirm(`Bạn muốn chuyển trạng thái thành ${row.original.status === 'active' ? 'Tạm ngưng' : 'Hoạt động'}?`)) {
                onToggleStatus(row.original.id!, row.original.status);
              }
            }}
            className="p-1.5 text-slate-400 hover:text-orange-600 hover:bg-orange-50 rounded-lg transition-colors"
            title={row.original.status === 'active' ? 'Xóa mềm (Ngưng hoạt động)' : 'Khôi phục (Hoạt động)'}
          >
            {row.original.status === 'active' ? <XCircle size={16} /> : <CheckCircle2 size={16} />}
          </button>

          <button 
            onClick={() => {
              if (confirm('CẢNH BÁO: Bạn đang thực hiện XÓA CỨNG.\nHành động này không thể hoàn tác và chỉ dành cho Admin. Bạn có chắc chắn?')) {
                onDelete(row.original.id!);
              }
            }}
            className="p-1.5 text-slate-300 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
            title="Xóa vĩnh viễn (Yêu cầu quyền Admin)"
          >
            <Trash2 size={16} />
          </button>
        </div>
      )
    }
  ];

  return (
    <div className="space-y-6 h-full flex flex-col">
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 shrink-0">
        <div>
          <h2 className="text-2xl font-black text-slate-800 dark:text-white tracking-tight flex items-center gap-2">
            <Building className="text-primary" /> Quản lý Chi nhánh
          </h2>
          <p className="text-sm text-slate-500 mt-1">Danh sách tất cả các kho và cửa hàng.</p>
        </div>
        <button 
          onClick={() => router.push('/warehouses/create')}
          className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white px-4 py-2 rounded-xl font-bold transition-all shadow-md active:scale-95"
        >
          <Plus size={18} /> Thêm Chi nhánh
        </button>
      </div>

      <div className="bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 shadow-sm flex gap-4 shrink-0">
        <div className="relative flex-1 max-w-md">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search size={16} className="text-slate-400" />
          </div>
          <input 
            type="text" 
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Tìm kiếm chi nhánh..." 
            className="w-full pl-10 pr-4 py-2 bg-slate-50 border border-slate-200 rounded-xl text-sm outline-none focus:border-primary"
          />
        </div>
      </div>

      <div className="flex-1 overflow-hidden">
        {isLoading ? (
          <div className="flex justify-center p-10"><span className="animate-spin text-primary">⏳</span></div>
        ) : (
          <DataTable 
            columns={columns} 
            data={data} 
            globalFilter={globalFilter}
            setGlobalFilter={setGlobalFilter}
          />
        )}
      </div>
    </div>
  );
}
