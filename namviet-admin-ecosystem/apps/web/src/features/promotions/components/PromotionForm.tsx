import React, { useEffect } from 'react';
import { Form, Input, Button, DatePicker, message, Space, Typography } from 'antd';
import { useCreatePromotion, useUpdatePromotion } from '../hooks';
import { PromotionsPromotion, PromotionsCreatePromotionRequest } from '@namviet/shared-types/src/backend.d';
import dayjs from 'dayjs';
import { MinusCircleOutlined, PlusOutlined } from '@ant-design/icons';

interface PromotionFormProps {
  initialData?: PromotionsPromotion;
  onSuccess: () => void;
  onCancel: () => void;
}

const { RangePicker } = DatePicker;
const { Text } = Typography;

export const PromotionForm: React.FC<PromotionFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const [form] = Form.useForm();
  const createMutation = useCreatePromotion();
  const updateMutation = useUpdatePromotion();

  useEffect(() => {
    if (initialData) {
      form.setFieldsValue({
        ...initialData,
        dates: [
          initialData.start_date ? dayjs(initialData.start_date) : undefined,
          initialData.end_date ? dayjs(initialData.end_date) : undefined,
        ],
        rules: initialData.rules || [],
      });
    } else {
      form.resetFields();
    }
  }, [initialData, form]);

  const onFinish = (values: any) => {
    const payload: PromotionsCreatePromotionRequest = {
      name: values.name,
      code: values.code,
      start_date: values.dates?.[0]?.toISOString() || new Date().toISOString(),
      end_date: values.dates?.[1]?.toISOString() || new Date().toISOString(),
      rules: values.rules || [],
    };

    if (initialData?.id) {
      updateMutation.mutate({ id: initialData.id, data: payload }, {
        onSuccess: () => {
          message.success('Cập nhật Voucher thành công');
          onSuccess();
        },
        onError: (err: any) => message.error(err.response?.data?.error || 'Lỗi cập nhật')
      });
    } else {
      createMutation.mutate(payload, {
        onSuccess: () => {
          message.success('Tạo Voucher mới thành công');
          onSuccess();
        },
        onError: (err: any) => message.error(err.response?.data?.error || 'Lỗi tạo mới')
      });
    }
  };

  return (
    <Form layout="vertical" form={form} onFinish={onFinish}>
      <Form.Item name="name" label="Tên Chương trình/Voucher" rules={[{ required: true }]}>
        <Input placeholder="VD: Khuyến mãi mùa hè" />
      </Form.Item>
      
      <Form.Item name="code" label="Mã áp dụng (Code)" rules={[{ required: true }]}>
        <Input placeholder="VD: SUMMER2024" className="uppercase" />
      </Form.Item>

      <Form.Item name="dates" label="Thời gian áp dụng" rules={[{ required: true }]}>
        <RangePicker className="w-full" format="DD/MM/YYYY HH:mm" showTime />
      </Form.Item>

      <div className="bg-slate-50 p-4 rounded-lg mb-4 border border-slate-200">
        <Text strong className="block mb-3 text-slate-700">Điều kiện / Luật áp dụng (Rules)</Text>
        <Form.List name="rules">
          {(fields, { add, remove }) => (
            <>
              {fields.map(({ key, name, ...restField }) => (
                <Space key={key} style={{ display: 'flex', marginBottom: 8 }} align="baseline">
                  <Form.Item
                    {...restField}
                    name={[name]}
                    rules={[{ required: true, message: 'Missing rule' }]}
                    className="mb-0 w-80"
                  >
                    <Input placeholder="VD: Giảm 10% tối đa 50k" />
                  </Form.Item>
                  <MinusCircleOutlined onClick={() => remove(name)} className="text-red-500" />
                </Space>
              ))}
              <Form.Item className="mb-0 mt-2">
                <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                  Thêm Rule
                </Button>
              </Form.Item>
            </>
          )}
        </Form.List>
      </div>

      <div className="flex justify-end gap-2 mt-4">
        <Button onClick={onCancel}>Hủy</Button>
        <Button 
          type="primary" 
          htmlType="submit" 
          loading={createMutation.isPending || updateMutation.isPending}
        >
          {initialData ? 'Cập nhật' : 'Tạo mới'}
        </Button>
      </div>
    </Form>
  );
};
