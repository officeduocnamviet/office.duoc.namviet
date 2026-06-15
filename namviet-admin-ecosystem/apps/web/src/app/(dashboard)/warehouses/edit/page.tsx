"use client";

import React, { Suspense } from 'react';
import { WarehouseForm } from '@/components/features/warehouses/warehouse-form';
import { useWarehouse, useUpdateWarehouse } from '@/hooks/queries/use-warehouses';
import { WarehouseFormData } from '@namviet/shared-types/src/warehouse.types';
import { toast } from 'sonner';
import { useRouter, useSearchParams } from 'next/navigation';

function EditWarehouseContent() {
  const searchParams = useSearchParams();
  const idParam = searchParams.get('id');
  const id = idParam ? parseInt(idParam, 10) : 0;
  const router = useRouter();

  const { data: warehouse, isLoading: isFetching, error } = useWarehouse(id);
  const updateMutation = useUpdateWarehouse();

  const handleSubmit = (data: WarehouseFormData) => {
    updateMutation.mutate({ id, data }, {
      onSuccess: () => {
        toast.success('Cập nhật chi nhánh thành công!');
      },
      onError: (err) => {
        toast.error(`Lỗi: ${err.message}`);
      }
    });
  };

  if (!id) {
    return <div className="p-6">Không tìm thấy mã chi nhánh</div>;
  }

  if (error) {
    return (
      <div className="p-6 h-[calc(100vh-80px)] flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 font-bold mb-4">Lỗi tải dữ liệu chi nhánh</p>
          <button onClick={() => router.push('/warehouses')} className="text-primary hover:underline">Quay lại danh sách</button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-6 h-[calc(100vh-80px)] overflow-y-auto bg-slate-50/50 dark:bg-slate-900/50">
      {isFetching ? (
        <div className="flex justify-center p-10"><span className="animate-spin text-primary">⏳</span></div>
      ) : (
        <WarehouseForm 
          initialData={warehouse}
          onSubmit={handleSubmit}
          isLoading={updateMutation.isPending}
        />
      )}
    </div>
  );
}

export default function EditWarehousePage() {
  return (
    <Suspense fallback={<div>Đang tải...</div>}>
      <EditWarehouseContent />
    </Suspense>
  );
}
