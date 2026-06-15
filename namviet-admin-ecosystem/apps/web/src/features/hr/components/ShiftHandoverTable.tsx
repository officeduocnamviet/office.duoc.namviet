import React, { useState } from 'react';
import { Table, Button, Space, Tag, Modal, Form, InputNumber } from 'antd';
import { useShiftHandovers, useCreateShiftHandover } from '../hooks';
import dayjs from 'dayjs';
import { CheckCircle, Clock } from 'lucide-react';

export const ShiftHandoverTable = () => {
  const { data: handovers, isLoading } = useShiftHandovers();
  const createHandoverMutation = useCreateShiftHandover();
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const handleFinish = (values: any) => {
    createHandoverMutation.mutate({
      actual_cash_submitted: values.actual_cash_submitted,
      assignment_id: 1, // Mock assignment ID for demo
      branch_id: 1,
      system_cash_amount: 1500000,
      system_cod_amount: 500000,
      user_id: 'EMPLOYEE_1',
    }, {
      onSuccess: () => {
        setIsModalOpen(false);
        form.resetFields();
      }
    });
  };

  const columns = [
    {
      title: 'Mã Bàn Giao',
      dataIndex: 'id',
      key: 'id',
      render: (id: string) => <span className="font-mono text-gray-600">{id}</span>
    },
    {
      title: 'Nhân viên',
      dataIndex: 'user_id',
      key: 'user_id',
    },
    {
      title: 'Tiền mặt (Hệ thống)',
      dataIndex: 'system_cash_amount',
      key: 'system_cash_amount',
      render: (val: number) => val ? `${val.toLocaleString()} đ` : '0 đ'
    },
    {
      title: 'Tiền mặt (Thực nộp)',
      dataIndex: 'actual_cash_submitted',
      key: 'actual_cash_submitted',
      render: (val: number) => val ? <strong className="text-blue-600">{val.toLocaleString()} đ</strong> : '0 đ'
    },
    {
      title: 'Chênh lệch',
      key: 'diff',
      render: (_: any, record: any) => {
        const diff = (record.actual_cash_submitted || 0) - (record.system_cash_amount || 0);
        if (diff === 0) return <span className="text-green-500 font-medium">Khớp</span>;
        if (diff > 0) return <span className="text-blue-500 font-medium">+{diff.toLocaleString()} đ</span>;
        return <span className="text-red-500 font-medium">{diff.toLocaleString()} đ</span>;
      }
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        if (status === 'APPROVED') return <Tag color="success" icon={<CheckCircle size={14} className="mr-1" />}>Đã duyệt</Tag>;
        if (status === 'PENDING') return <Tag color="warning" icon={<Clock size={14} className="mr-1" />}>Chờ duyệt</Tag>;
        return <Tag>{status}</Tag>;
      }
    },
    {
      title: 'Thời gian',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (val: string) => val ? dayjs(val).format('DD/MM/YYYY HH:mm') : '---'
    }
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-lg font-bold text-slate-800">Lịch sử Bàn giao ca</h2>
        <Button type="primary" className="bg-blue-600" onClick={() => setIsModalOpen(true)}>
          Tạo Bàn Giao Ca
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={handovers} 
        rowKey="id" 
        loading={isLoading}
        pagination={{ pageSize: 15 }}
      />

      <Modal
        title="Tạo Bàn Giao Ca"
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleFinish} className="mt-4">
          <div className="bg-slate-50 p-4 rounded-lg mb-4 space-y-2 text-sm text-slate-600">
            <div className="flex justify-between">
              <span>Tổng tiền mặt trên hệ thống:</span>
              <strong className="text-slate-800">1,500,000 đ</strong>
            </div>
            <div className="flex justify-between">
              <span>Tổng tiền COD trên hệ thống:</span>
              <strong className="text-slate-800">500,000 đ</strong>
            </div>
          </div>

          <Form.Item 
            name="actual_cash_submitted" 
            label="Tiền mặt thực tế nộp"
            rules={[{ required: true, message: 'Vui lòng nhập số tiền thực nộp' }]}
          >
            <InputNumber 
              className="w-full" 
              formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
              parser={(value) => value!.replace(/\$\s?|(,*)/g, '')}
              prefix="VND"
            />
          </Form.Item>

          <div className="flex justify-end gap-2 mt-6">
            <Button onClick={() => setIsModalOpen(false)}>Hủy</Button>
            <Button type="primary" htmlType="submit" className="bg-blue-600" loading={createHandoverMutation.isPending}>
              Xác nhận Bàn giao
            </Button>
          </div>
        </Form>
      </Modal>
    </div>
  );
};
