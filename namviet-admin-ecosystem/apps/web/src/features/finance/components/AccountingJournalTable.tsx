import React, { useState } from 'react';
import { Table, Button, Tag, Space, Input, DatePicker } from 'antd';
import { Search, BookText } from 'lucide-react';
import { useAccountingJournals } from '../hooks';
import dayjs from 'dayjs';

export const AccountingJournalTable = () => {
  const { data: journals = [], isLoading } = useAccountingJournals();
  const [searchText, setSearchText] = useState('');

  const filteredJournals = journals.filter((j: any) => 
    j.description?.toLowerCase().includes(searchText.toLowerCase()) ||
    j.transaction_code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã GD',
      dataIndex: 'transaction_code',
      key: 'transaction_code',
      render: (code: string) => <span className="font-mono text-gray-600">{code}</span>
    },
    {
      title: 'Ngày hạch toán',
      key: 'date',
      render: (_: any, record: any) => dayjs(record.journal_date).format('DD/MM/YYYY')
    },
    {
      title: 'Diễn giải',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'TK Nợ',
      dataIndex: 'account_debit',
      key: 'account_debit',
      render: (acc: string) => <span className="font-semibold text-blue-600">{acc}</span>
    },
    {
      title: 'TK Có',
      dataIndex: 'account_credit',
      key: 'account_credit',
      render: (acc: string) => <span className="font-semibold text-orange-600">{acc}</span>
    },
    {
      title: 'Số tiền',
      dataIndex: 'amount',
      key: 'amount',
      align: 'right' as const,
      render: (amount: number) => (
        <span className="font-bold">
          {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount || 0)}
        </span>
      )
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Sổ Nhật ký chung</h2>
          <p className="text-sm text-gray-500">Lịch sử hạch toán kế toán toàn hệ thống</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm diễn giải, mã GD..." 
            prefix={<Search className="w-4 h-4 text-gray-400" />}
            value={searchText}
            onChange={e => setSearchText(e.target.value)}
            className="w-full sm:w-64"
          />
          <DatePicker.RangePicker className="w-full sm:w-auto" format="DD/MM/YYYY" />
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredJournals} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 20 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredJournals?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có dữ liệu.</div>
          ) : (
            filteredJournals?.map((journal: any) => (
              <div key={journal.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-slate-50 text-slate-500">
                    <BookText className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{journal.description}</h3>
                    <div className="text-xs font-mono text-gray-500 mt-0.5">{journal.transaction_code} • {dayjs(journal.journal_date).format('DD/MM/YYYY')}</div>
                  </div>
                </div>
                
                <div className="bg-gray-50 p-3 rounded-lg flex flex-col gap-2 mt-2 text-sm">
                  <div className="flex justify-between items-center border-b border-gray-100 pb-2">
                    <div className="flex flex-col">
                      <span className="text-xs text-gray-500">Nợ</span>
                      <span className="font-semibold text-blue-600">{journal.account_debit}</span>
                    </div>
                    <div className="flex flex-col text-right">
                      <span className="text-xs text-gray-500">Có</span>
                      <span className="font-semibold text-orange-600">{journal.account_credit}</span>
                    </div>
                  </div>
                  <div className="flex justify-between items-center text-gray-600 pt-1">
                    <span>Số tiền:</span>
                    <span className="font-bold text-gray-800">{new Intl.NumberFormat('vi-VN').format(journal.amount || 0)}đ</span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
};
