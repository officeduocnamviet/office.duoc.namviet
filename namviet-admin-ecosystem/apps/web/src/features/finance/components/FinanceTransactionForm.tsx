import React, { useEffect } from 'react';
import { Form, Input, Select, Button, InputNumber } from 'antd';
import { useCreateFinanceTransaction, useUpdateFinanceTransaction } from '../hooks';
import { FinanceTransaction } from '../api';

const { TextArea } = Input;

interface Props {
  initialData?: any;
  onSuccess: () => void;
  onCancel: () => void;
}

export const FinanceTransactionForm: React.FC<Props> = ({ initialData, onSuccess, onCancel }) => {
  const [form] = Form.useForm();
  const createMutation = useCreateFinanceTransaction();
  const updateMutation = useUpdateFinanceTransaction();

  useEffect(() => {
    if (initialData) {
      form.setFieldsValue({
        code: (initialData as any).transaction_code || (initialData as any).code,
        flow: initialData.flow,
        amount: initialData.amount,
        description: initialData.description,
        business_type: initialData.business_type,
        fund_account_id: initialData.fund_account_id,
      });
    } else {
      form.resetFields();
      form.setFieldsValue({
        flow: 'IN',
        amount: 0,
      });
    }
  }, [initialData, form]);

  const onFinish = (values: any) => {
    if (initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id.toString(), data: values },
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
      <div className="grid grid-cols-2 gap-4">
        <Form.Item
          name="flow"
          label="Loại giao dịch"
          rules={[{ required: true }]}
        >
          <Select
            options={[
              { value: 'IN', label: 'Phiếu Thu' },
              { value: 'OUT', label: 'Phiếu Chi' },
            ]}
          />
        </Form.Item>

        <Form.Item
          name="amount"
          label="Số tiền"
          rules={[{ required: true }]}
        >
          <InputNumber 
            className="w-full" 
            formatter={(value) => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
            parser={(value) => value!.replace(/\$\s?|(,*)/g, '') as any}
          />
        </Form.Item>
      </div>

      <Form.Item
        name="fund_account_id"
        label="Tài khoản / Quỹ (ID)"
        rules={[{ required: true, message: 'Vui lòng chọn quỹ' }]}
      >
        <InputNumber className="w-full" placeholder="ID của Quỹ" />
      </Form.Item>

      <Form.Item
        name="description"
        label="Lý do thu / chi"
        rules={[{ required: true, message: 'Vui lòng nhập lý do' }]}
      >
        <TextArea rows={3} placeholder="VD: Thanh toán tiền điện tháng 6..." />
      </Form.Item>

      <Form.Item
        name="business_type"
        label="Nghiệp vụ"
      >
        <Input placeholder="VD: PAYMENT, REFUND, INVENTORY_PURCHASE..." />
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
