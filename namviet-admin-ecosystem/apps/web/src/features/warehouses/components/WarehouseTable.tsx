import React, { useState } from 'react';
import { Table, Button, Space, Tag, Modal, Popconfirm } from 'antd';
import { Plus, Pencil, Trash2, MapPin, Building2 } from 'lucide-react';
import { WarehousesWarehouse } from '@namviet/shared-types/src/backend.d';
import { useWarehouses, useDeleteWarehouse } from '../hooks';
import { toast } from 'sonner';
import { WarehouseForm } from './WarehouseForm';

export const WarehouseTable = () => {
  const { data: warehouses, isLoading } = useWarehouses();
  const deleteMutation = useDeleteWarehouse();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingWarehouse, setEditingWarehouse] = useState<WarehousesWarehouse | null>(null);

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa chi nhánh thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const columns = [
    {
      title: 'Tên Chi Nhánh',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <div className="flex items-center gap-2 font-medium text-gray-900">
          <Building2 className="w-4 h-4 text-blue-500" />
          {text}
        </div>
      ),
    },
    {
      title: 'Phân loại',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        <Tag color={type === 'MAIN' ? 'blue' : 'default'}>
          {type === 'MAIN' ? 'Tổng Kho' : 'Chi nhánh'}
        </Tag>
      ),
    },
    {
      title: 'Địa chỉ',
      dataIndex: 'address',
      key: 'address',
      render: (address: string) => (
        <div className="flex items-center gap-1 text-gray-500">
          <MapPin className="w-3 h-3" />
          <span className="truncate max-w-[200px]">{address || 'Chưa cập nhật'}</span>
        </div>
      ),
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'ACTIVE' ? 'success' : 'error'}>
          {status === 'ACTIVE' ? 'Hoạt động' : 'Đóng cửa'}
        </Tag>
      ),
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: WarehousesWarehouse) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingWarehouse(record);
              setIsModalOpen(true);
            }}
          />
          <Popconfirm
            title="Xóa chi nhánh?"
            description="Bạn có chắc chắn muốn xóa chi nhánh này không?"
            onConfirm={() => record.id && handleDelete(record.id)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger icon={<Trash2 className="w-4 h-4" />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Quản lý Chi Nhánh (Kho)</h2>
          <p className="text-sm text-gray-500">Danh sách các cơ sở, chi nhánh thuộc Nam Việt</p>
        </div>
        <Button 
          type="primary" 
          icon={<Plus className="w-4 h-4" />}
          className="bg-blue-600"
          onClick={() => {
            setEditingWarehouse(null);
            setIsModalOpen(true);
          }}
        >
          Thêm Chi Nhánh
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={warehouses} 
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editingWarehouse ? "Cập nhật Chi nhánh" : "Thêm Chi nhánh mới"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
        width={600}
      >
        <WarehouseForm 
          initialData={editingWarehouse} 
          onSuccess={() => setIsModalOpen(false)}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>
    </div>
  );
};
