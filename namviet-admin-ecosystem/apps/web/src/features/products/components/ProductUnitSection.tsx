import React, { useState } from 'react';
import { Table, Button, Space, Modal, Form, Input, InputNumber, Popconfirm } from 'antd';
import { Plus, Pencil, Trash2, ArrowRight } from 'lucide-react';
import { ProductUnitsProductUnit } from '@namviet/shared-types/src/backend.d';
import { useProductUnits, useCreateProductUnit, useUpdateProductUnit, useDeleteProductUnit } from '../hooks/useProductUnits';
import { toast } from 'sonner';

interface ProductUnitSectionProps {
  productId: number;
  baseUnit: string;
}

export const ProductUnitSection: React.FC<ProductUnitSectionProps> = ({ productId, baseUnit }) => {
  const { data: units, isLoading } = useProductUnits(productId);
  const createMutation = useCreateProductUnit(productId);
  const updateMutation = useUpdateProductUnit(productId);
  const deleteMutation = useDeleteProductUnit(productId);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUnit, setEditingUnit] = useState<ProductUnitsProductUnit | null>(null);
  const [form] = Form.useForm();

  const handleDelete = (unitId: number) => {
    deleteMutation.mutate(unitId, {
      onSuccess: () => toast.success('Xóa đơn vị quy đổi thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const openModal = (unit?: ProductUnitsProductUnit) => {
    setEditingUnit(unit || null);
    if (unit) {
      form.setFieldsValue(unit);
    } else {
      form.resetFields();
      form.setFieldValue('conversion_factor', 1);
    }
    setIsModalOpen(true);
  };

  const handleSubmit = (values: any) => {
    if (editingUnit?.id) {
      updateMutation.mutate(
        { unitId: editingUnit.id, data: values },
        {
          onSuccess: () => {
            toast.success('Cập nhật thành công');
            setIsModalOpen(false);
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    } else {
      createMutation.mutate(
        values,
        {
          onSuccess: () => {
            toast.success('Thêm đơn vị thành công');
            setIsModalOpen(false);
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    }
  };

  const columns = [
    {
      title: 'Tên Đơn vị',
      dataIndex: 'unit_name',
      key: 'unit_name',
      render: (text: string) => <strong className="text-blue-600">{text}</strong>
    },
    {
      title: 'Hệ số quy đổi',
      key: 'conversion',
      render: (_: any, record: ProductUnitsProductUnit) => (
        <div className="flex items-center gap-2 text-sm text-gray-600">
          <span>1 {record.unit_name}</span>
          <ArrowRight className="w-3 h-3 text-gray-400" />
          <strong className="text-gray-900">{record.conversion_factor}</strong> {baseUnit}
        </div>
      ),
    },
    {
      title: 'Mã vạch (Barcode)',
      dataIndex: 'barcode',
      key: 'barcode',
      render: (text: string) => text || '---'
    },
    {
      title: 'Giá bán riêng (VNĐ)',
      dataIndex: 'price_sell',
      key: 'price_sell',
      render: (price: number) => price ? price.toLocaleString('vi-VN') + 'đ' : <span className="text-gray-400 italic">Theo tỷ lệ hệ số</span>
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 100,
      render: (_: any, record: ProductUnitsProductUnit) => (
        <Space size="small">
          <Button type="text" size="small" icon={<Pencil className="w-4 h-4 text-gray-500" />} onClick={() => openModal(record)} />
          <Popconfirm
            title="Xóa đơn vị quy đổi?"
            onConfirm={() => record.id && handleDelete(record.id)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger size="small" icon={<Trash2 className="w-4 h-4" />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="mt-4">
      <div className="flex justify-between items-center mb-4">
        <div>
          <h3 className="font-medium text-gray-800">Danh sách Đơn vị quy đổi</h3>
          <p className="text-xs text-gray-500">Đơn vị cơ bản của sản phẩm này là: <strong className="text-blue-600">{baseUnit}</strong></p>
        </div>
        <Button type="dashed" icon={<Plus className="w-4 h-4" />} onClick={() => openModal()}>
          Thêm Quy đổi
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={units} 
        rowKey="id"
        loading={isLoading}
        pagination={false}
        size="small"
        bordered
      />

      <Modal
        title={editingUnit ? "Cập nhật Đơn vị tính" : "Thêm Đơn vị tính"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleSubmit} className="mt-4">
          <Form.Item 
            name="unit_name" 
            label="Tên đơn vị" 
            rules={[{ required: true, message: 'Vui lòng nhập tên đơn vị (VD: Hộp, Vỉ...)' }]}
          >
            <Input placeholder="VD: Hộp, Vỉ, Thùng..." />
          </Form.Item>

          <Form.Item 
            name="conversion_factor" 
            label={`Hệ số quy đổi (1 Đơn vị này = bao nhiêu ${baseUnit}?)`}
            rules={[{ required: true, message: 'Vui lòng nhập hệ số quy đổi' }]}
          >
            <InputNumber className="w-full" min={1} />
          </Form.Item>

          <Form.Item name="barcode" label="Mã vạch riêng cho đơn vị này">
            <Input placeholder="Quét mã vạch..." />
          </Form.Item>

          <Form.Item name="price_sell" label="Giá bán tùy chỉnh (VNĐ)" help="Để trống nếu muốn hệ thống tự tính giá = Giá cơ bản × Hệ số quy đổi">
            <InputNumber className="w-full" formatter={value => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')} parser={value => value?.replace(/\$\s?|(,*)/g, '') as any} />
          </Form.Item>

          <div className="flex justify-end gap-2 mt-6">
            <Button onClick={() => setIsModalOpen(false)}>Hủy</Button>
            <Button type="primary" htmlType="submit" className="bg-blue-600">
              Lưu đơn vị
            </Button>
          </div>
        </Form>
      </Modal>
    </div>
  );
};
