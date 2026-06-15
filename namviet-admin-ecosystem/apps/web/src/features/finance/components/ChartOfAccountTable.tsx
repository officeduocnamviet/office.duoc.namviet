import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input } from 'antd';
import { Plus, Pencil, Search, BookOpen } from 'lucide-react';
import { useChartOfAccounts } from '../hooks';
import { ChartOfAccount } from '../api';

export const ChartOfAccountTable = () => {
  const { data: accounts = [], isLoading } = useChartOfAccounts();
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [searchText, setSearchText] = useState('');

  const filteredAccounts = accounts.filter((a: any) => 
    a.account_name?.toLowerCase().includes(searchText.toLowerCase()) ||
    a.account_code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã TK',
      dataIndex: 'account_code',
      key: 'account_code',
      render: (code: string) => <span className="font-mono font-bold text-blue-700">{code}</span>
    },
    {
      title: 'Tên Tài khoản',
      dataIndex: 'account_name',
      key: 'account_name',
      render: (name: string) => <span className="font-semibold text-gray-800">{name}</span>
    },
    {
      title: 'Loại',
      dataIndex: 'account_type',
      key: 'account_type',
      render: (type: string) => (
        type === 'ASSET' ? <Tag color="blue">Tài sản (Asset)</Tag> :
        type === 'LIABILITY' ? <Tag color="orange">Nợ phải trả (Liability)</Tag> :
        type === 'EQUITY' ? <Tag color="purple">Vốn CSH (Equity)</Tag> :
        type === 'REVENUE' ? <Tag color="green">Doanh thu (Revenue)</Tag> :
        <Tag color="red">Chi phí (Expense)</Tag>
      )
    },
    {
      title: 'Hạch toán',
      dataIndex: 'allow_posting',
      key: 'allow_posting',
      render: (allow: boolean) => (
        allow ? <Tag color="success">Được phép</Tag> : <Tag color="default">Chỉ tổng hợp</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 80,
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button type="text" icon={<Pencil size={16} className="text-blue-600" />} />
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Hệ thống Tài khoản (COA)</h2>
          <p className="text-sm text-gray-500">Danh mục tài khoản kế toán chuẩn mực</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm theo mã, tên TK..." 
            prefix={<Search className="w-4 h-4 text-gray-400" />}
            value={searchText}
            onChange={e => setSearchText(e.target.value)}
            className="w-full sm:w-64"
          />
          <Button 
            type="primary" 
            icon={<Plus className="w-4 h-4" />}
            className="bg-blue-600 flex items-center"
            onClick={() => setDrawerVisible(true)}
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
            pagination={{ pageSize: 20 }}
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
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-500">
                    <BookOpen className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{acc.account_name}</h3>
                    <div className="text-xs font-mono font-bold text-blue-700 mt-0.5">{acc.account_code}</div>
                  </div>
                  {acc.allow_posting ? (
                    <Tag color="success" className="m-0">Hạch toán</Tag>
                  ) : (
                    <Tag color="default" className="m-0">Tổng hợp</Tag>
                  )}
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      <Drawer
        title="Thêm Tài khoản mới"
        width={400}
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        destroyOnClose
      >
        <div className="p-4 text-center text-gray-500">
          Chức năng thêm tài khoản kế toán đang được hoàn thiện.
        </div>
      </Drawer>
    </div>
  );
};
