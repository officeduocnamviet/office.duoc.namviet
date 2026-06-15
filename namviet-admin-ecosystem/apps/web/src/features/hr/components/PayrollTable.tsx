import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, DatePicker, Select, Divider } from 'antd';
import { Filter, Banknote, PlayCircle, Eye } from 'lucide-react';
import { usePayrolls } from '../hooks';
import dayjs from 'dayjs';

export const PayrollTable = () => {
  const { data: payrolls = [], isLoading } = usePayrolls();
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [selectedMonth, setSelectedMonth] = useState<string>(dayjs().format('MM-YYYY'));

  const filteredPayrolls = payrolls.filter((p: any) => 
    p.period_month === parseInt(selectedMonth.split('-')[0]) &&
    p.period_year === parseInt(selectedMonth.split('-')[1])
  );

  const columns = [
    {
      title: 'Mã Bảng Lương',
      dataIndex: 'payroll_code',
      key: 'payroll_code',
      render: (code: string) => <span className="font-mono text-gray-600 font-medium">{code}</span>
    },
    {
      title: 'Kỳ Lương',
      key: 'period',
      render: (_: any, record: any) => <span className="font-semibold text-blue-700">Tháng {record.period_month}/{record.period_year}</span>
    },
    {
      title: 'Tổng Lương Cơ Bản',
      dataIndex: 'total_base_salary',
      key: 'total_base_salary',
      align: 'right' as const,
      render: (amount: number) => new Intl.NumberFormat('vi-VN').format(amount || 0)
    },
    {
      title: 'Phụ cấp / Thưởng',
      dataIndex: 'total_allowance',
      key: 'total_allowance',
      align: 'right' as const,
      render: (amount: number) => <span className="text-green-600">+{new Intl.NumberFormat('vi-VN').format(amount || 0)}</span>
    },
    {
      title: 'Khấu trừ / Phạt',
      dataIndex: 'total_deduction',
      key: 'total_deduction',
      align: 'right' as const,
      render: (amount: number) => <span className="text-red-600">-{new Intl.NumberFormat('vi-VN').format(amount || 0)}</span>
    },
    {
      title: 'Thực Lĩnh',
      dataIndex: 'net_salary',
      key: 'net_salary',
      align: 'right' as const,
      render: (amount: number) => <span className="font-bold text-gray-800 text-base">{new Intl.NumberFormat('vi-VN').format(amount || 0)}</span>
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'PAID' ? <Tag color="success">Đã thanh toán</Tag> : 
        status === 'APPROVED' ? <Tag color="blue">Đã duyệt</Tag> : 
        <Tag color="warning">Nháp (Draft)</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 120,
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button type="text" className="text-blue-600 hover:text-blue-700 hover:bg-blue-50 p-1" title="Chi tiết">
            <Eye className="w-5 h-5" />
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Quản lý Bảng Lương (Payrolls)</h2>
          <p className="text-sm text-gray-500">Tính toán và quản lý lương nhân viên theo kỳ</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <DatePicker 
            picker="month"
            className="w-full sm:w-40" 
            defaultValue={dayjs()}
            onChange={(date) => setSelectedMonth(date ? date.format('MM-YYYY') : dayjs().format('MM-YYYY'))} 
            format="MM/YYYY"
          />
          <Button 
            icon={<Filter className="w-4 h-4" />}
            onClick={() => setIsFilterDrawerOpen(true)}
          >
            Lọc
          </Button>
          <Button 
            type="primary" 
            icon={<PlayCircle className="w-4 h-4" />}
            className="bg-green-600 hover:bg-green-700 border-none"
          >
            Chạy tính lương
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden lg:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredPayrolls} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 15 }}
          />
        </div>

        <div className="block lg:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredPayrolls?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Chưa có bảng lương cho tháng này.</div>
          ) : (
            filteredPayrolls?.map((payroll: any) => (
              <div key={payroll.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-green-50 text-green-500">
                    <Banknote className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm">Tháng {payroll.period_month}/{payroll.period_year}</h3>
                    <div className="text-xs text-gray-500 font-mono mt-0.5">{payroll.payroll_code}</div>
                  </div>
                  {payroll.status === 'PAID' ? <Tag color="success" className="m-0">Đã TT</Tag> : 
                   payroll.status === 'APPROVED' ? <Tag color="blue" className="m-0">Đã duyệt</Tag> : 
                   <Tag color="warning" className="m-0">Nháp</Tag>}
                </div>
                
                <div className="bg-gray-50 p-3 rounded-lg flex flex-col gap-2 mt-2 text-sm">
                  <div className="flex justify-between items-center text-gray-600">
                    <span>Lương cơ bản:</span>
                    <span className="font-medium">{new Intl.NumberFormat('vi-VN').format(payroll.total_base_salary || 0)}đ</span>
                  </div>
                  <div className="flex justify-between items-center text-gray-600">
                    <span>Phụ cấp/Thưởng:</span>
                    <span className="text-green-600">+{new Intl.NumberFormat('vi-VN').format(payroll.total_allowance || 0)}đ</span>
                  </div>
                  <div className="flex justify-between items-center text-gray-600">
                    <span>Khấu trừ:</span>
                    <span className="text-red-600">-{new Intl.NumberFormat('vi-VN').format(payroll.total_deduction || 0)}đ</span>
                  </div>
                  <div className="flex justify-between items-center pt-2 border-t border-gray-200 mt-1">
                    <span className="font-medium text-gray-800">Thực lĩnh:</span>
                    <span className="font-bold text-gray-800 text-base">{new Intl.NumberFormat('vi-VN').format(payroll.net_salary || 0)}đ</span>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      
      <Drawer
        title="Lọc Bảng lương"
        placement="right"
        onClose={() => setIsFilterDrawerOpen(false)}
        open={isFilterDrawerOpen}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 400}
        styles={{ body: { padding: '24px' } }}
      >
        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'DRAFT', label: 'Nháp' },
                { value: 'APPROVED', label: 'Đã duyệt' },
                { value: 'PAID', label: 'Đã thanh toán' }
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
