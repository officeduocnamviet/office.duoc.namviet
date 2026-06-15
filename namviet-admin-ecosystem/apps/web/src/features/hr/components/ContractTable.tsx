import React, { useState } from 'react';
import { Table, Button, Tag, Modal, Form, Input, DatePicker, Select } from 'antd';
import { useEmploymentContracts, useCreateEmploymentContract } from '../hooks';
import dayjs from 'dayjs';
import { FileSignature, Plus } from 'lucide-react';

export const ContractTable = () => {
  const { data: contracts, isLoading } = useEmploymentContracts();
  const createMutation = useCreateEmploymentContract();
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const handleFinish = (values: any) => {
    createMutation.mutate({
      user_id: values.user_id,
      contract_code: values.contract_type === 'THU_VIEC' ? 'PROBATION' : 'OFFICIAL',
      valid_from: values.start_date.toISOString(),
      valid_to: values.end_date?.toISOString(),
      base_salary: 10000000
    }, {
      onSuccess: () => {
        setIsModalOpen(false);
        form.resetFields();
      }
    });
  };

  const columns = [
    {
      title: 'Nhân viên',
      dataIndex: 'user_id',
      key: 'user_id',
    },
    {
      title: 'Loại HĐ',
      dataIndex: 'contract_type',
      key: 'contract_type',
      render: (type: string) => <Tag color="blue">{type}</Tag>
    },
    {
      title: 'Từ ngày',
      dataIndex: 'start_date',
      key: 'start_date',
      render: (val: string) => val ? dayjs(val).format('DD/MM/YYYY') : '---'
    },
    {
      title: 'Đến ngày',
      dataIndex: 'end_date',
      key: 'end_date',
      render: (val: string) => val ? dayjs(val).format('DD/MM/YYYY') : 'Vô thời hạn'
    },
    {
      title: 'Lương cơ bản',
      dataIndex: 'base_salary',
      key: 'base_salary',
      render: (val: number) => val ? `${val.toLocaleString()} đ` : '---'
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'ACTIVE' ? 'success' : 'default'}>{status}</Tag>
      )
    }
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-lg font-bold text-slate-800">Danh sách Hợp đồng</h2>
        <Button type="primary" className="bg-blue-600" icon={<Plus size={16}/>} onClick={() => setIsModalOpen(true)}>
          Thêm Hợp đồng
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={contracts} 
        rowKey="id" 
        loading={isLoading}
        pagination={{ pageSize: 15 }}
      />

      <Modal
        title="Thêm Hợp đồng mới"
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleFinish} className="mt-4">
          <Form.Item name="user_id" label="Mã Nhân viên" rules={[{ required: true }]}>
            <Input placeholder="Nhập mã nhân viên" />
          </Form.Item>
          <Form.Item name="contract_type" label="Loại hợp đồng" rules={[{ required: true }]}>
            <Select>
              <Select.Option value="THU_VIEC">Thử việc (2 tháng)</Select.Option>
              <Select.Option value="CO_THOI_HAN">Có thời hạn (1 năm)</Select.Option>
              <Select.Option value="VO_THOI_HAN">Vô thời hạn</Select.Option>
            </Select>
          </Form.Item>
          <div className="flex gap-4">
            <Form.Item name="start_date" label="Từ ngày" rules={[{ required: true }]} className="flex-1">
              <DatePicker className="w-full" format="DD/MM/YYYY" />
            </Form.Item>
            <Form.Item name="end_date" label="Đến ngày" className="flex-1">
              <DatePicker className="w-full" format="DD/MM/YYYY" />
            </Form.Item>
          </div>
          <div className="flex justify-end gap-2 mt-6">
            <Button onClick={() => setIsModalOpen(false)}>Hủy</Button>
            <Button type="primary" htmlType="submit" className="bg-blue-600" loading={createMutation.isPending}>
              Tạo mới
            </Button>
          </div>
        </Form>
      </Modal>
    </div>
  );
};
