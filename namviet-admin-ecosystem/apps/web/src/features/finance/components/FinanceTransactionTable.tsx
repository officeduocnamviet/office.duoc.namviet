import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input, Select, Divider } from 'antd';
import { Search, Filter, Plus, Pencil, ArrowRightLeft } from 'lucide-react';
import { useFinanceTransactions, useDeleteFinanceTransaction } from '../hooks';
import { FinanceTransactionForm } from './FinanceTransactionForm';
import { FinanceTransaction } from '../api';
import dayjs from 'dayjs';

export const FinanceTransactionTable = () => {
  const { data: transactions = [], isLoading } = useFinanceTransactions();
  const deleteMutation = useDeleteFinanceTransaction();
  
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingTransaction, setEditingTransaction] = useState<FinanceTransaction | undefined>();
  const [searchText, setSearchText] = useState('');

  const filteredTransactions = transactions.filter((t: any) => 
    t.transaction_code?.toLowerCase().includes(searchText.toLowerCase()) ||
    t.description?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã GD',
      dataIndex: 'transaction_code',
      key: 'transaction_code',
      render: (code: string) => <span className="font-mono text-gray-600">{code}</span>
    },
    {
      title: 'Thời gian',
      key: 'time',
      render: (_: any, record: any) => dayjs(record.transaction_date).format('DD/MM/YYYY HH:mm')
    },
    {
      title: 'Loại GD',
      dataIndex: 'flow',
      key: 'flow',
      render: (flow: string) => (
        flow === 'IN' ? <Tag color="success">Thu tiền</Tag> : <Tag color="error">Chi tiền</Tag>
      )
    },
    {
      title: 'Số tiền',
      dataIndex: 'amount',
      key: 'amount',
      align: 'right' as const,
      render: (amount: number, record: any) => (
        <span className={`font-bold ${record.flow === 'IN' ? 'text-green-600' : 'text-red-600'}`}>
          {record.flow === 'IN' ? '+' : '-'} {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount || 0)}
        </span>
      )
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'COMPLETED' ? <Tag color="success">Thành công</Tag> : 
        status === 'PENDING' ? <Tag color="warning">Chờ xử lý</Tag> : 
        status === 'CANCELLED' ? <Tag color="default">Đã hủy</Tag> : 
        <Tag color="default">{status}</Tag>
      )
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Giao dịch Tài chính (Thu/Chi)</h2>
          <p className="text-sm text-gray-500">Quản lý dòng tiền vào và ra của hệ thống</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm mã GD, lý do..." 
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
              setEditingTransaction(undefined);
              setDrawerVisible(true);
            }}
            className="bg-blue-600"
          >
            Thêm Giao dịch
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredTransactions} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 15 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredTransactions?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có dữ liệu.</div>
          ) : (
            filteredTransactions?.map((t: any) => (
              <div key={t.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className={`w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 ${t.flow === 'IN' ? 'bg-green-50 text-green-500' : 'bg-red-50 text-red-500'}`}>
                    <ArrowRightLeft className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{t.description || 'Giao dịch mới'}</h3>
                    <div className="text-xs text-gray-500 font-mono mt-0.5">{t.transaction_code} • {dayjs(t.transaction_date).format('DD/MM/YYYY')}</div>
                  </div>
                  <div className={`font-bold ${t.flow === 'IN' ? 'text-green-600' : 'text-red-600'}`}>
                    {t.flow === 'IN' ? '+' : '-'}{new Intl.NumberFormat('vi-VN').format(t.amount || 0)}đ
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      <Drawer
        title={editingTransaction ? 'Sửa Giao dịch' : 'Thêm Giao dịch mới'}
        placement="right"
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 500}
        destroyOnClose
      >
        <FinanceTransactionForm 
          initialData={editingTransaction}
          onSuccess={() => setDrawerVisible(false)}
          onCancel={() => setDrawerVisible(false)}
        />
      </Drawer>
    </div>
  );
};
