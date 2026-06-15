import React, { useEffect } from 'react';
import { Form, Input, Select, Button, InputNumber } from 'antd';
import { useCreateFundAccount, useUpdateFundAccount } from '../hooks';
import { FundAccount } from '../api';

interface Props {
  initialData?: any;
  onSuccess: () => void;
  onCancel: () => void;
}

export const FundAccountForm: React.FC<Props> = ({ initialData, onSuccess, onCancel }) => {
  const [form] = Form.useForm();
  const createMutation = useCreateFundAccount();
  const updateMutation = useUpdateFundAccount();

  useEffect(() => {
    if (initialData) {
      form.setFieldsValue({
        code: (initialData as any).account_code || (initialData as any).code,
        name: initialData.name,
        type: initialData.type,
        currency: initialData.currency,
        current_balance: initialData.current_balance,
        status: initialData.status,
      });
    } else {
      form.resetFields();
      form.setFieldsValue({
        type: 'CASH',
        currency: 'VND',
        status: 'ACTIVE',
        current_balance: 0,
      });
    }
  }, [initialData, form]);

  const onFinish = (values: any) => {
    if (initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data: values },
        { onSuccess }
      );
    } else {
      createMutation.mutate(values, { onSuccess });
    }
  };

  return (
    <Form
      form={form}
      layout="vertical"
      onFinish={onFinish}
    >
      <Form.Item
        name="code"
        label="Mã tài khoản"
        rules={[{ required: true, message: 'Vui lòng nhập mã' }]}
      >
        <Input placeholder="VD: TM01, NH01..." />
      </Form.Item>

      <Form.Item
        name="name"
        label="Tên tài khoản/quỹ"
        rules={[{ required: true, message: 'Vui lòng nhập tên' }]}
      >
        <Input placeholder="VD: Tiền mặt tại quỹ, Vietcombank..." />
      </Form.Item>

      <div className="grid grid-cols-2 gap-4">
        <Form.Item
          name="type"
          label="Loại"
          rules={[{ required: true }]}
        >
          <Select
            options={[
              { value: 'CASH', label: 'Tiền mặt' },
              { value: 'BANK', label: 'Ngân hàng' },
              { value: 'E_WALLET', label: 'Ví điện tử' },
            ]}
          />
        </Form.Item>

        <Form.Item
          name="currency"
          label="Tiền tệ"
          rules={[{ required: true }]}
        >
          <Select
            options={[
              { value: 'VND', label: 'VND' },
              { value: 'USD', label: 'USD' },
            ]}
          />
        </Form.Item>
      </div>

      <Form.Item
        name="current_balance"
        label="Số dư ban đầu"
      >
        <InputNumber 
          className="w-full" 
          formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
          parser={(value) => value!.replace(/\$\s?|(,*)/g, '') as any}
        />
      </Form.Item>

      <Form.Item
        name="status"
        label="Trạng thái"
        rules={[{ required: true }]}
      >
        <Select
          options={[
            { value: 'ACTIVE', label: 'Hoạt động' },
            { value: 'INACTIVE', label: 'Ngừng hoạt động' },
          ]}
        />
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
