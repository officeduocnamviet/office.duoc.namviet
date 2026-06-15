import React, { useState } from 'react';
import { Table, Button, Tag, Tabs, Modal, Form, Input } from 'antd';
import { useMedicalVectors, useCreateMedicalVector, useProductVectors, useCreateProductVector } from '../hooks/useKnowledge';
import dayjs from 'dayjs';
import { Database, Plus } from 'lucide-react';

export const KnowledgeVectorTable = () => {
  const [activeTab, setActiveTab] = useState('medical');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const { data: medicalVectors, isLoading: loadMed } = useMedicalVectors();
  const createMed = useCreateMedicalVector();

  const { data: productVectors, isLoading: loadProd } = useProductVectors();
  const createProd = useCreateProductVector();

  const handleFinish = (values: any) => {
    if (activeTab === 'medical') {
      createMed.mutate({
        title: values.title,
        content: values.content,
        metadata: [values.keywords],
      }, {
        onSuccess: () => {
          setIsModalOpen(false);
          form.resetFields();
        }
      });
    } else {
      createProd.mutate({
        product_id: values.product_id,
        content: values.content,
        metadata: [],
      }, {
        onSuccess: () => {
          setIsModalOpen(false);
          form.resetFields();
        }
      });
    }
  };

  const medColumns = [
    { title: 'Chủ đề / Bệnh lý', dataIndex: 'title', key: 'title', render: (t: string) => <strong className="text-slate-800">{t}</strong> },
    { title: 'Metadata', dataIndex: 'metadata', key: 'metadata', render: (meta: string[]) => meta?.map((m: string) => <Tag key={m} color="blue">{m}</Tag>) },
    { title: 'Trạng thái Nhúng', dataIndex: 'embedding_status', key: 'embedding_status', render: () => <Tag color="success">COMPLETED</Tag> },
    { title: 'Cập nhật', dataIndex: 'created_at', key: 'created_at', render: (d: string) => d ? dayjs(d).format('DD/MM/YYYY HH:mm') : '---' },
  ];

  const prodColumns = [
    { title: 'Mã Sản phẩm', dataIndex: 'product_id', key: 'product_id', render: (t: number) => <span className="font-mono">{t}</span> },
    { title: 'Nội dung Vector', dataIndex: 'content', key: 'content', render: (t: string) => <div className="truncate max-w-md">{t}</div> },
    { title: 'Trạng thái Nhúng', dataIndex: 'embedding_status', key: 'embedding_status', render: () => <Tag color="success">COMPLETED</Tag> },
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
          <Database className="text-indigo-600 w-5 h-5"/>
          Quản trị Vector Database
        </h2>
        <Button type="primary" className="bg-indigo-600" icon={<Plus size={16}/>} onClick={() => setIsModalOpen(true)}>
          Thêm Tri thức
        </Button>
      </div>

      <Tabs 
        activeKey={activeTab} 
        onChange={setActiveTab}
        items={[
          {
            key: 'medical',
            label: 'Tri thức Y tế (Medical)',
            children: <Table columns={medColumns} dataSource={medicalVectors} rowKey="id" loading={loadMed} />
          },
          {
            key: 'product',
            label: 'Vector Sản phẩm (Products)',
            children: <Table columns={prodColumns} dataSource={productVectors} rowKey="id" loading={loadProd} />
          }
        ]}
      />

      <Modal
        title={activeTab === 'medical' ? "Thêm Tri thức Y tế" : "Thêm Vector Sản phẩm"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleFinish} className="mt-4">
          {activeTab === 'medical' && (
            <>
              <Form.Item name="title" label="Chủ đề / Tên bệnh lý" rules={[{ required: true }]}>
                <Input placeholder="VD: Cảm cúm" />
              </Form.Item>
              <Form.Item name="keywords" label="Từ khóa / Triệu chứng" rules={[{ required: true }]}>
                <Input placeholder="VD: Đau đầu, sốt, ho" />
              </Form.Item>
            </>
          )}
          {activeTab === 'product' && (
            <Form.Item name="product_id" label="Mã Sản phẩm" rules={[{ required: true }]}>
              <Input placeholder="Nhập mã sản phẩm" />
            </Form.Item>
          )}
          <Form.Item name="content" label="Nội dung chuyên sâu" rules={[{ required: true }]}>
            <Input.TextArea rows={4} placeholder="Nội dung dùng để tạo Vector embeddings..." />
          </Form.Item>
          <div className="flex justify-end gap-2 mt-6">
            <Button onClick={() => setIsModalOpen(false)}>Hủy</Button>
            <Button type="primary" htmlType="submit" className="bg-indigo-600" loading={createMed.isPending || createProd.isPending}>
              Tạo Vector
            </Button>
          </div>
        </Form>
      </Modal>
    </div>
  );
};
