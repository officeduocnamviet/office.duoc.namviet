import React, { useState } from 'react';
import { Table, Button, Tag, Modal, Form, Input, DatePicker, Select } from 'antd';
import { useTrainingCourses, useCreateTrainingCourse } from '../hooks';
import dayjs from 'dayjs';
import { BookOpen, Plus } from 'lucide-react';

export const TrainingTable = () => {
  const { data: courses, isLoading } = useTrainingCourses();
  const createMutation = useCreateTrainingCourse();
  
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const handleFinish = (values: any) => {
    createMutation.mutate({
      title: values.name,
      content_type: 'VIDEO',
      passing_score: Number(values.passing_score) || 0,
    }, {
      onSuccess: () => {
        setIsModalOpen(false);
        form.resetFields();
      }
    });
  };

  const columns = [
    {
      title: 'Tên Khóa học',
      dataIndex: 'title',
      key: 'title',
      render: (text: string) => <strong className="text-slate-800">{text}</strong>
    },
    {
      title: 'Loại nội dung',
      dataIndex: 'content_type',
      key: 'content_type',
    },
    {
      title: 'Điểm qua môn',
      dataIndex: 'passing_score',
      key: 'passing_score',
      render: (val: number) => val ? val.toString() : '---'
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'UPCOMING' ? 'blue' : status === 'IN_PROGRESS' ? 'processing' : 'success'}>{status}</Tag>
      )
    }
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-lg font-bold text-slate-800">Khóa Đào tạo Nội bộ</h2>
        <Button type="primary" className="bg-blue-600" icon={<Plus size={16}/>} onClick={() => setIsModalOpen(true)}>
          Thêm Khóa học
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={courses} 
        rowKey="id" 
        loading={isLoading}
        pagination={{ pageSize: 15 }}
      />

      <Modal
        title="Thêm Khóa Đào tạo"
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleFinish} className="mt-4">
          <Form.Item name="name" label="Tên khóa học" rules={[{ required: true }]}>
            <Input placeholder="VD: Đào tạo Hội nhập Nhân viên mới" />
          </Form.Item>
          <Form.Item name="passing_score" label="Điểm qua môn" rules={[{ required: true }]}>
            <Input type="number" placeholder="VD: 80" />
          </Form.Item>
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
