import React from 'react';
import { Drawer, Tabs, Table, Tag, Empty } from 'antd';
import { CustomersCustomer } from '@namviet/shared-types/src/backend.d';
import { CustomerForm } from './CustomerForm';
import { useCustomerVaccinations, useCustomerVouchers } from '../hooks';
import dayjs from 'dayjs';
import { Syringe, Ticket, User } from 'lucide-react';

interface CustomerProfileDrawerProps {
  customer?: CustomersCustomer;
  visible: boolean;
  onClose: () => void;
}

export const CustomerProfileDrawer: React.FC<CustomerProfileDrawerProps> = ({ customer, visible, onClose }) => {
  const isNew = !customer;

  return (
    <Drawer
      title={isNew ? "Thêm mới Khách hàng" : <span className="font-semibold">Hồ sơ: {customer.name}</span>}
      width={700}
      onClose={onClose}
      open={visible}
      destroyOnClose
    >
      {isNew ? (
        <CustomerForm 
          initialData={undefined} 
          onSuccess={onClose} 
          onCancel={onClose} 
        />
      ) : (
        <Tabs defaultActiveKey="1" items={[
          {
            key: '1',
            label: <span className="flex items-center gap-2"><User size={16}/> Thông tin chung</span>,
            children: <CustomerForm initialData={customer} onSuccess={onClose} onCancel={onClose} />
          },
          {
            key: '2',
            label: <span className="flex items-center gap-2"><Syringe size={16}/> Sổ Tiêm chủng</span>,
            children: <VaccinationTab customerId={customer.id!} />
          },
          {
            key: '3',
            label: <span className="flex items-center gap-2"><Ticket size={16}/> Vouchers</span>,
            children: <VoucherTab customerId={customer.id!} />
          }
        ]} />
      )}
    </Drawer>
  );
};

const VaccinationTab = ({ customerId }: { customerId: number }) => {
  const { data: vaccinations, isLoading } = useCustomerVaccinations(customerId);

  const columns = [
    { title: 'Vắc xin', dataIndex: 'vaccine_name', key: 'vaccine_name', render: (t: string) => <strong className="text-blue-700">{t}</strong> },
    { title: 'Ngày tiêm', dataIndex: 'vaccination_date', key: 'vaccination_date', render: (d: string) => dayjs(d).format('DD/MM/YYYY') },
    { title: 'Mũi tiêm', dataIndex: 'dose_number', key: 'dose_number', render: (d: number) => `Mũi ${d}` },
    { title: 'Lịch nhắc', dataIndex: 'next_due_date', key: 'next_due_date', render: (d: string) => d ? <Tag color="warning">{dayjs(d).format('DD/MM/YYYY')}</Tag> : '---' },
  ];

  return (
    <div className="mt-4">
      <Table 
        columns={columns} 
        dataSource={vaccinations} 
        rowKey="id" 
        loading={isLoading} 
        pagination={false} 
        locale={{ emptyText: <Empty description="Chưa có lịch sử tiêm chủng" /> }}
      />
    </div>
  );
};

const VoucherTab = ({ customerId }: { customerId: number }) => {
  const { data: vouchers, isLoading } = useCustomerVouchers(customerId);

  const columns = [
    { title: 'Mã Voucher', dataIndex: 'voucher_code', key: 'voucher_code', render: (t: string) => <Tag color="blue" className="font-mono text-base">{t}</Tag> },
    { title: 'Loại giảm giá', dataIndex: 'discount_type', key: 'discount_type', render: (t: string) => t === 'PERCENT' ? '%' : 'VND' },
    { title: 'Mức giảm', dataIndex: 'discount_value', key: 'discount_value', render: (v: number) => <strong className="text-red-500">{v.toLocaleString()}</strong> },
    { title: 'Hạn sử dụng', dataIndex: 'valid_until', key: 'valid_until', render: (d: string) => dayjs(d).format('DD/MM/YYYY') },
    { title: 'Trạng thái', dataIndex: 'status', key: 'status', render: (s: string) => <Tag color={s === 'ACTIVE' ? 'success' : 'default'}>{s}</Tag> },
  ];

  return (
    <div className="mt-4">
      <Table 
        columns={columns} 
        dataSource={vouchers} 
        rowKey="id" 
        loading={isLoading} 
        pagination={false} 
        locale={{ emptyText: <Empty description="Không có voucher nào" /> }}
      />
    </div>
  );
};
