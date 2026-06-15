import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input, Select, Divider, Popconfirm } from 'antd';
import { Search, Wallet, Filter, Plus, Pencil, Trash2 } from 'lucide-react';
import { useFundAccounts, useDeleteFundAccount } from '../hooks';
import { FundAccountForm } from './FundAccountForm';
import { FundAccount } from '../api';

export const FundAccountTable = () => {
  const { data: accounts = [], isLoading } = useFundAccounts();
  const deleteMutation = useDeleteFundAccount();
  
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingAccount, setEditingAccount] = useState<FundAccount | undefined>();
  const [searchText, setSearchText] = useState('');

  const filteredAccounts = accounts.filter((a: any) => 
    a.name?.toLowerCase().includes(searchText.toLowerCase()) ||
    a.code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Tài khoản',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: any) => (
        <div>
          <div className="font-semibold text-gray-800">{name}</div>
          <div className="text-xs text-gray-500 font-mono">{record.code}</div>
        </div>
      )
    },
    {
      title: 'Loại',
      dataIndex: 'type',
      key: 'type',
      render: (type: string) => (
        type === 'CASH' ? <Tag color="orange">Tiền mặt</Tag> : 
        type === 'BANK' ? <Tag color="blue">Ngân hàng</Tag> : 
        <Tag color="purple">Ví điện tử</Tag>
      )
    },
    {
      title: 'Số dư hiện tại',
      dataIndex: 'current_balance',
      key: 'current_balance',
      align: 'right' as const,
      render: (balance: number, record: any) => (
        <span className="font-bold text-green-600">
          {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: record.currency || 'VND' }).format(balance || 0)}
        </span>
      )
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'ACTIVE' ? <Tag color="success">Hoạt động</Tag> : <Tag color="default">Đóng</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 120,
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button 
            type="text" 
            className="text-blue-600 hover:text-blue-700 hover:bg-blue-50"
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingAccount(record as FundAccount);
              setDrawerVisible(true);
            }}
          />
          <Popconfirm
            title="Bạn có chắc chắn muốn xóa tài khoản này?"
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button 
              type="text" 
              danger
              className="hover:bg-red-50"
              icon={<Trash2 className="w-4 h-4" />} 
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Quỹ & Tài khoản Ngân hàng</h2>
          <p className="text-sm text-gray-500">Quản lý số dư, tiền mặt, tiền gửi ngân hàng</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm tên, mã tài khoản..." 
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
              setEditingAccount(undefined);
              setDrawerVisible(true);
            }}
            className="bg-blue-600"
          >
            Thêm Tài khoản
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredAccounts} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 15 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredAccounts?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có dữ liệu.</div>
          ) : (
            filteredAccounts?.map((acc: any) => (
              <div key={acc.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className={`w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 ${
                    acc.type === 'CASH' ? 'bg-orange-50 text-orange-500' :
                    acc.type === 'BANK' ? 'bg-blue-50 text-blue-500' :
                    'bg-purple-50 text-purple-500'
                  }`}>
                    <Wallet className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{acc.name}</h3>
                    <div className="text-xs text-gray-500 font-mono mt-0.5">{acc.code}</div>
                  </div>
                  {acc.status === 'ACTIVE' ? (
                    <Tag color="success" className="m-0">HĐ</Tag>
                  ) : (
                    <Tag color="default" className="m-0">Đóng</Tag>
                  )}
                </div>
                
                <div className="bg-gray-50 p-3 rounded-lg flex justify-between items-center">
                  <span className="text-sm text-gray-500">Số dư:</span>
                  <span className="font-bold text-green-600 text-base">
                    {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: acc.currency || 'VND' }).format(acc.current_balance || 0)}
                  </span>
                </div>

                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2 mt-1">
                  <Button 
                    type="text" 
                    size="small" 
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingAccount(acc as FundAccount);
                      setDrawerVisible(true);
                    }}
                  >
                    Sửa
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      <Drawer
        title={editingAccount ? 'Sửa Tài khoản' : 'Thêm Tài khoản mới'}
        placement="right"
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 500}
        destroyOnClose
      >
        <FundAccountForm 
          initialData={editingAccount}
          onSuccess={() => setDrawerVisible(false)}
          onCancel={() => setDrawerVisible(false)}
        />
      </Drawer>
    </div>
  );
};
