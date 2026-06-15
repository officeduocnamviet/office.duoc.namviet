import React, { useState } from 'react';
import { Table, Button, Space, Drawer, Popconfirm, Tag, Input, Divider, Select } from 'antd';
import { Plus, Pencil, Trash2, Factory, Search, Filter } from 'lucide-react';
import { ManufacturersManufacturer } from '@namviet/shared-types/src/backend.d';
import { useManufacturers, useDeleteManufacturer } from '../hooks';
import { toast } from 'sonner';
import { ManufacturerForm } from './ManufacturerForm';

export const ManufacturerTable = () => {
  const { data: manufacturers, isLoading } = useManufacturers();
  const deleteMutation = useDeleteManufacturer();

  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingManufacturer, setEditingManufacturer] = useState<ManufacturersManufacturer | null>(null);
  const [searchText, setSearchText] = useState('');

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa NSX thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const columns = [
    {
      title: 'Nhà sản xuất / Đối tác',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <div className="flex items-center gap-2 font-medium text-gray-900">
          <Factory className="w-4 h-4 text-blue-500" />
          {text}
        </div>
      ),
    },
    {
      title: 'Quốc gia',
      dataIndex: 'country',
      key: 'country',
      render: (country: string) => <span className="text-gray-600">{country || 'Chưa cập nhật'}</span>
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'ACTIVE' ? 'success' : 'default'}>
          {status === 'ACTIVE' ? 'Đang hợp tác' : 'Ngừng hợp tác'}
        </Tag>
      ),
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: ManufacturersManufacturer) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingManufacturer(record);
              setIsDrawerOpen(true);
            }}
          />
          <Popconfirm
            title="Xóa nhà sản xuất?"
            description="Xóa đối tác này có thể ảnh hưởng đến sản phẩm liên quan."
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
          <h2 className="text-lg font-semibold text-gray-900">Danh sách Nhà sản xuất</h2>
          <p className="text-sm text-gray-500">Quản lý các đối tác cung cấp, sản xuất Dược phẩm</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            prefix={<Search className="w-4 h-4 text-gray-400" />} 
            placeholder="Tìm nhà sản xuất..." 
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
              setEditingManufacturer(null);
              setIsDrawerOpen(true);
            }}
          >
            Thêm
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={manufacturers?.filter(m => m.name?.toLowerCase().includes(searchText.toLowerCase()))} 
            rowKey="id"
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : manufacturers?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có nhà sản xuất nào.</div>
          ) : (
            manufacturers?.filter(m => m.name?.toLowerCase().includes(searchText.toLowerCase())).map((manu) => (
              <div key={manu.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-center">
                  <div className="w-10 h-10 bg-indigo-50 rounded-lg flex items-center justify-center text-indigo-500 flex-shrink-0">
                    <Factory className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-gray-900 text-sm truncate">{manu.name}</h3>
                    <div className="text-xs text-gray-500 mt-1 truncate">Quốc gia: {manu.country || '---'}</div>
                  </div>
                  <Tag color={manu.status === 'ACTIVE' ? 'success' : 'default'} className="m-0">
                    {manu.status === 'ACTIVE' ? 'HĐ' : 'Ẩn'}
                  </Tag>
                </div>
                
                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                  <Button 
                    type="text" 
                    size="small"
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingManufacturer(manu);
                      setIsDrawerOpen(true);
                    }}
                  >
                    Sửa
                  </Button>
                  <Popconfirm
                    title="Xóa đối tác?"
                    onConfirm={() => manu.id && handleDelete(manu.id)}
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

      <Drawer
        title={editingManufacturer ? "Cập nhật NSX" : "Thêm Nhà sản xuất mới"}
        width={450}
        onClose={() => setIsDrawerOpen(false)}
        open={isDrawerOpen}
        destroyOnClose
      >
        <ManufacturerForm 
          initialData={editingManufacturer} 
          onSuccess={() => setIsDrawerOpen(false)}
          onCancel={() => setIsDrawerOpen(false)}
        />
      </Drawer>

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
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả trạng thái" 
              allowClear
              options={[
                { value: 'ACTIVE', label: 'Đang hợp tác' },
                { value: 'INACTIVE', label: 'Ngừng hợp tác' }
              ]} 
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
