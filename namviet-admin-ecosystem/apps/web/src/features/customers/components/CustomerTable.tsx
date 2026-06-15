import React, { useState } from 'react';
import { Table, Button, Tag, Space, Popconfirm, Drawer, Input, Modal, Typography, Select, Divider } from 'antd';
import { Plus, Pencil, Trash2, Search, Building2, User2, Filter } from 'lucide-react';
import { useCustomers, useDeleteCustomer } from '../hooks';
import { CustomerForm } from './CustomerForm';
import { CustomersCustomer } from '@namviet/shared-types/src/backend.d';
import { CustomerProfileDrawer } from './CustomerProfileDrawer';

const { Text } = Typography;

export const CustomerTable = () => {
  const { data: customers = [], isLoading } = useCustomers();
  const deleteMutation = useDeleteCustomer();
  
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingCustomer, setEditingCustomer] = useState<CustomersCustomer | undefined>();
  const [searchText, setSearchText] = useState('');

  const filteredCustomers = customers.filter(c => 
    c.name?.toLowerCase().includes(searchText.toLowerCase()) ||
    c.phone?.includes(searchText) ||
    c.email?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã KH',
      dataIndex: 'customer_code',
      key: 'customer_code',
      width: 100,
      render: (code: string) => <Text strong className="text-gray-600">{code}</Text>
    },
    {
      title: 'Khách hàng',
      key: 'name',
      render: (_: any, record: CustomersCustomer) => (
        <div>
          <div className="font-semibold text-gray-800">{record.name}</div>
          <div className="text-xs text-gray-400">{record.phone}</div>
        </div>
      )
    },
    {
      title: 'Loại',
      dataIndex: 'customer_type',
      key: 'customer_type',
      render: (type: string) => (
        type === 'B2B' 
          ? <Tag color="blue" icon={<Building2 size={12} className="mr-1 inline" />}>Doanh nghiệp</Tag>
          : <Tag color="green" icon={<User2 size={12} className="mr-1 inline" />}>Cá nhân</Tag>
      )
    },
    {
      title: 'Công nợ',
      dataIndex: 'current_debt',
      key: 'current_debt',
      render: (val: number) => (
        <span className={val > 0 ? 'text-red-500 font-semibold' : 'text-gray-500'}>
          {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(val || 0)}
        </span>
      )
    },
    {
      title: 'Điểm HT',
      dataIndex: 'loyalty_points',
      key: 'loyalty_points',
      render: (val: number) => <span className="font-medium text-amber-500">{val || 0}</span>
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 120,
      render: (_: any, record: CustomersCustomer) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil size={16} className="text-blue-600" />} 
            onClick={() => {
              setEditingCustomer(record);
              setDrawerVisible(true);
            }}
          />
          <Popconfirm
            title="Xóa khách hàng này?"
            onConfirm={() => deleteMutation.mutate(record.id!)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger icon={<Trash2 size={16} />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Khách hàng & Đối tác</h2>
          <p className="text-sm text-gray-500">Quản lý danh sách khách hàng B2B và B2C</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm theo tên, SĐT..." 
            prefix={<Search className="w-4 h-4 text-gray-400" />}
            value={searchText}
            onChange={e => setSearchText(e.target.value)}
            className="w-full sm:w-64"
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
            onClick={() => {
              setEditingCustomer(undefined);
              setDrawerVisible(true);
            }}
            className="bg-blue-600 flex items-center"
          >
            Thêm mới
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredCustomers} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredCustomers?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có khách hàng nào.</div>
          ) : (
            filteredCustomers?.map((customer) => (
              <div key={customer.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-center">
                  <div className={`w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 ${customer.customer_type === 'B2B' ? 'bg-blue-50 text-blue-500' : 'bg-green-50 text-green-500'}`}>
                    {customer.customer_type === 'B2B' ? <Building2 className="w-5 h-5" /> : <User2 className="w-5 h-5" />}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-gray-900 text-sm truncate">{customer.name}</h3>
                    <div className="text-xs text-gray-500 mt-1 truncate">SĐT: {customer.phone || '---'}</div>
                  </div>
                  <Tag color={customer.customer_type === 'B2B' ? 'blue' : 'green'} className="m-0">
                    {customer.customer_type === 'B2B' ? 'B2B' : 'Cá nhân'}
                  </Tag>
                </div>
                
                <div className="flex justify-between items-center bg-gray-50 p-2 rounded">
                  <div className="text-xs text-gray-500">
                    Công nợ: <span className={customer.current_debt! > 0 ? 'text-red-500 font-semibold' : 'text-gray-900 font-semibold'}>{new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(customer.current_debt || 0)}</span>
                  </div>
                  <div className="text-xs text-gray-500">
                    Điểm HT: <span className="text-amber-500 font-semibold">{customer.loyalty_points || 0}</span>
                  </div>
                </div>

                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                  <Button 
                    type="text" 
                    size="small"
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingCustomer(customer);
                      setDrawerVisible(true);
                    }}
                  >
                    Sửa
                  </Button>
                  <Popconfirm
                    title="Xóa khách hàng?"
                    onConfirm={() => customer.id && deleteMutation.mutate(customer.id)}
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

      <CustomerProfileDrawer
        customer={editingCustomer}
        visible={drawerVisible}
        onClose={() => setDrawerVisible(false)}
      />

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
            <label className="block text-sm font-medium text-gray-700 mb-1">Loại khách hàng</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'B2B', label: 'Doanh nghiệp (B2B)' },
                { value: 'RETAIL', label: 'Cá nhân (B2C)' }
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
