import React, { useState } from 'react';
import { Table, Button, Space, Modal, Popconfirm, Tag, Input, Drawer, Select, Divider } from 'antd';
import { Plus, Pencil, Trash2, Search, Box, Filter } from 'lucide-react';
import { BatchesBatch } from '@namviet/shared-types/src/backend.d';
import { useBatches, useDeleteBatch } from '../hooks';
import { useProducts } from '@/features/products/hooks/useProducts';
import { toast } from 'sonner';
import { BatchForm } from './BatchForm';
import dayjs from 'dayjs';

export const BatchTable = () => {
  const { data: batches, isLoading } = useBatches();
  const { data: products } = useProducts();
  const deleteMutation = useDeleteBatch();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingBatch, setEditingBatch] = useState<BatchesBatch | null>(null);
  const [searchText, setSearchText] = useState('');

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa Lô thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const getProductName = (productId?: number) => {
    const product = products?.find(p => p.id === productId);
    return product ? product.name : 'Không xác định';
  };

  const filteredBatches = batches?.filter(b => 
    b.batch_code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã Lô',
      dataIndex: 'batch_code',
      key: 'batch_code',
      render: (text: string) => <strong className="text-blue-600">{text}</strong>,
    },
    {
      title: 'Sản phẩm',
      dataIndex: 'product_id',
      key: 'product',
      render: (id: number) => (
        <div className="flex items-center gap-2">
          <Box className="w-4 h-4 text-gray-400" />
          {getProductName(id)}
        </div>
      )
    },
    {
      title: 'Hạn sử dụng',
      dataIndex: 'expiry_date',
      key: 'expiry_date',
      render: (date: string) => {
        if (!date) return '---';
        const isExpiring = dayjs(date).diff(dayjs(), 'day') < 90; // Sắp hết hạn trong 90 ngày
        return (
          <Tag color={isExpiring ? 'red' : 'green'}>
            {dayjs(date).format('DD/MM/YYYY')}
          </Tag>
        );
      }
    },

    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: BatchesBatch) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingBatch(record);
              setIsModalOpen(true);
            }}
          />
          <Popconfirm
            title="Xóa Lô?"
            description="Lưu ý: Xóa lô có thể làm hỏng dữ liệu nhập xuất nếu lô đã phát sinh giao dịch."
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
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Quản lý Lô Hàng</h2>
          <p className="text-sm text-gray-500">Quản lý thông tin Lô, Hạn sử dụng của các sản phẩm</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            prefix={<Search className="w-4 h-4 text-gray-400" />} 
            placeholder="Tìm theo Mã lô..." 
            className="w-full sm:w-64"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <Button 
            icon={<Filter className="w-4 h-4" />}
            onClick={() => setIsFilterDrawerOpen(true)}
          >
            Lọc
          </Button>
          <Button 
            type="primary" 
            icon={<Plus className="w-4 h-4" />}
            className="bg-blue-600"
            onClick={() => {
              setEditingBatch(null);
              setIsModalOpen(true);
            }}
          >
            Tạo Lô mới
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredBatches} 
            rowKey="id"
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredBatches?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có lô hàng nào.</div>
          ) : (
            filteredBatches?.map((batch) => (
              <div key={batch.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-center">
                  <div className="w-10 h-10 bg-amber-50 rounded-lg flex items-center justify-center text-amber-500 flex-shrink-0">
                    <Box className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-blue-600 text-sm truncate">{batch.batch_code}</h3>
                    <div className="text-xs text-gray-700 font-medium truncate mt-1">SP: {getProductName(batch.product_id)}</div>
                    <div className="text-xs text-gray-500 mt-1">HSD: {batch.expiry_date ? dayjs(batch.expiry_date).format('DD/MM/YYYY') : '---'}</div>
                  </div>
                </div>
                
                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                  <Button 
                    type="text" 
                    size="small"
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingBatch(batch);
                      setIsModalOpen(true);
                    }}
                  >
                    Sửa
                  </Button>
                  <Popconfirm
                    title="Xóa lô?"
                    onConfirm={() => batch.id && handleDelete(batch.id)}
                    okText="Xóa"
                    cancelText="Hủy"
                  >
                    <Button 
                      type="text" 
                      size="small"
                      danger
                      className="bg-red-50 hover:bg-red-100 px-3"
                      icon={<Trash2 className="w-3.5 h-3.5" />} 
                    >
                      Xóa
                    </Button>
                  </Popconfirm>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      <Modal
        title={editingBatch ? "Cập nhật Thông tin Lô" : "Tạo Lô mới"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
        width={600}
      >
        <BatchForm 
          initialData={editingBatch} 
          onSuccess={() => setIsModalOpen(false)}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>

      <Drawer
        title={<span className="font-semibold text-gray-800">Lọc & Tìm kiếm</span>}
        placement="right"
        onClose={() => setIsFilterDrawerOpen(false)}
        open={isFilterDrawerOpen}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 500}
        styles={{ body: { padding: '24px' } }}
      >
        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái Hạn sử dụng</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'EXPIRING_SOON', label: 'Sắp hết hạn (Dưới 90 ngày)' },
                { value: 'EXPIRED', label: 'Đã hết hạn' }
              ]} 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Sản phẩm</label>
            <Select 
              className="w-full" 
              placeholder="Chọn sản phẩm" 
              allowClear 
              options={products?.map(p => ({ value: p.id, label: p.name }))}
            />
          </div>
          <Divider />
          <div className="flex justify-end gap-2">
            <Button onClick={() => setIsFilterDrawerOpen(false)}>Hủy</Button>
            <Button type="primary" className="bg-blue-600" onClick={() => setIsFilterDrawerOpen(false)}>Áp dụng</Button>
          </div>
        </div>
      </Drawer>
    </div>
  );
};
