import React, { useEffect } from 'react';
import { Form, Input, Button } from 'antd';
import { useCreateSystemConfig, useUpdateSystemConfig } from '../hooks/useConfigs';
import { SystemConfig } from '../api/configApi';

const { TextArea } = Input;

interface Props {
  initialData?: SystemConfig;
  onSuccess: () => void;
  onCancel: () => void;
}

export const ConfigForm: React.FC<Props> = ({ initialData, onSuccess, onCancel }) => {
  const [form] = Form.useForm();
  const createMutation = useCreateSystemConfig();
  const updateMutation = useUpdateSystemConfig();

  useEffect(() => {
    if (initialData) {
      form.setFieldsValue({
        config_key: initialData.config_key,
        config_value: JSON.stringify(initialData.config_value, null, 2),
        description: initialData.description,
      });
    } else {
      form.resetFields();
    }
  }, [initialData, form]);

  const onFinish = (values: any) => {
    let parsedValue = values.config_value;
    try {
      parsedValue = JSON.parse(values.config_value);
    } catch (e) {
      // If not valid JSON, leave it as string
    }

    const payload = {
      config_key: values.config_key,
      config_value: parsedValue,
      description: values.description,
    };

    if (initialData?.config_key) {
      updateMutation.mutate(
        { key: initialData.config_key, data: payload },
        { onSuccess }
      );
    } else {
      createMutation.mutate(payload, { onSuccess });
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={onFinish}
    >
      <Form.Item
        name="config_key"
        label="Khóa cấu hình (Key)"
        rules={[{ required: true, message: 'Vui lòng nhập khóa cấu hình' }]}
      >
        <Input placeholder="VD: PAYMENT_GATEWAY_CONFIG" disabled={!!initialData} />
      </Form.Item>

      <Form.Item
        name="description"
        label="Mô tả"
      >
        <Input placeholder="Nhập mô tả chức năng của cấu hình" />
      </Form.Item>

      <Form.Item
        name="config_value"
        label="Giá trị (JSON/Text)"
        rules={[{ required: true, message: 'Vui lòng nhập giá trị cấu hình' }]}
      >
        <TextArea rows={12} placeholder='VD: { "apiKey": "123", "secret": "abc" }' className="font-mono text-sm" />
      </Form.Item>

      <div className="flex justify-end gap-2 mt-8">
        <Button onClick={onCancel}>Hủy</Button>
        <Button 
          type="primary" 
          htmlType="submit"
          className="bg-blue-600"
          loading={createMutation.isPending || updateMutation.isPending}
        >
          {initialData ? 'Cập nhật' : 'Tạo mới'}
        </Button>
      </div>
    </Form>
  );
};
