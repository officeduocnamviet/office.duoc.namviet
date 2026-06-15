import React, { useState } from 'react';
import { Table, Button, Tag, Modal, Form, Input, DatePicker, Select } from 'antd';
import { useCampaigns, useCreateCampaign } from '../hooks/useCampaigns';
import dayjs from 'dayjs';
import { Megaphone, Plus } from 'lucide-react';
import { MarketingCampaign } from '../api/campaignApi';

export const CampaignTable = () => {
  const { data: campaigns, isLoading } = useCampaigns();
  const createMutation = useCreateCampaign();
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const handleFinish = (values: any) => {
    createMutation.mutate({
      name: values.name,
      description: values.description,
      start_date: values.start_date.toISOString(),
      end_date: values.end_date?.toISOString(),
      budget: values.budget || 0,
      target_segment: values.target_segment || ['ALL']
    }, {
      onSuccess: () => {
        setIsModalOpen(false);
        form.resetFields();
      }
    });
  };

  const columns = [
    {
      title: 'Tên chiến dịch',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => <strong className="text-slate-800">{text}</strong>
    },
    {
      title: 'Tập khách hàng (Segment)',
      dataIndex: 'target_segment',
      key: 'target_segment',
      render: (segments: string[]) => segments?.map(s => <Tag color="blue" key={s}>{s}</Tag>)
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
      render: (val: string) => val ? dayjs(val).format('DD/MM/YYYY') : '---'
    },
    {
      title: 'Ngân sách',
      dataIndex: 'budget',
      key: 'budget',
      render: (val: number) => val ? `${val.toLocaleString()} đ` : '---'
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'ACTIVE' ? 'success' : status === 'DRAFT' ? 'default' : 'warning'}>{status}</Tag>
      )
    }
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
          <Megaphone className="text-blue-600 w-5 h-5"/>
          Danh sách Chiến dịch
        </h2>
        <Button type="primary" className="bg-blue-600" icon={<Plus size={16}/>} onClick={() => setIsModalOpen(true)}>
          Thêm Chiến dịch
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={campaigns} 
        rowKey="id" 
        loading={isLoading}
        pagination={{ pageSize: 15 }}
      />

      <Modal
        title="Thêm Chiến dịch mới"
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleFinish} className="mt-4">
          <Form.Item name="name" label="Tên chiến dịch" rules={[{ required: true }]}>
            <Input placeholder="VD: Khuyến mãi mùa hè" />
          </Form.Item>
          <Form.Item name="target_segment" label="Tập khách hàng" rules={[{ required: true }]}>
            <Select mode="multiple">
              <Select.Option value="ALL">Tất cả (All)</Select.Option>
              <Select.Option value="VIP">Khách VIP</Select.Option>
              <Select.Option value="NEW">Khách Mới</Select.Option>
            </Select>
          </Form.Item>
          <Form.Item name="budget" label="Ngân sách (VNĐ)" rules={[{ required: true }]}>
            <Input type="number" placeholder="10000000" />
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
