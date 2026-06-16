import React, { useEffect, useState } from 'react';
import { Form, Input, Select, Button, Row, Col, Switch, DatePicker, message } from 'antd';
import { useCreateCustomer, useUpdateCustomer } from '../hooks';
import { CustomersCustomer, CustomersCreateCustomerRequest } from '@namviet/shared-types/src/backend.d';
import dayjs from 'dayjs';

interface CustomerFormProps {
  initialData?: CustomersCustomer;
  onSuccess: () => void;
  onCancel: () => void;
}

const { Option } = Select;
const { TextArea } = Input;

export const CustomerForm: React.FC<CustomerFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const [form] = Form.useForm();
  const createMutation = useCreateCustomer();
  const updateMutation = useUpdateCustomer();
  const [isB2B, setIsB2B] = useState(false);

  useEffect(() => {
    if (initialData) {
      setIsB2B(initialData.customer_type === 'B2B');
      form.setFieldsValue({
        ...initialData,
        dob: initialData.dob ? dayjs(initialData.dob) : undefined,
        b2b_tax_code: (initialData.b2b_metadata as any)?.tax_code || '',
        b2b_company: (initialData.b2b_metadata as any)?.company || '',
      });
    } else {
      setIsB2B(false);
      form.resetFields();
      form.setFieldsValue({ customer_type: 'B2C' });
    }
  }, [initialData, form]);

  const onFinish = (values: any) => {
    const payload: CustomersCreateCustomerRequest = {
      name: values.name,
      phone: values.phone,
      email: values.email,
      address: values.address,
      customer_type: isB2B ? 'B2B' : 'B2C',
      gender: values.gender,
      cccd: values.cccd,
      dob: values.dob ? values.dob.format('YYYY-MM-DD') : undefined,
      b2b_metadata: isB2B ? ({ tax_code: values.b2b_tax_code || '', company: values.b2b_company || '' } as any) : undefined,
    };

    if (initialData?.id) {
      updateMutation.mutate({ id: initialData.id, data: payload }, {
        onSuccess: () => {
          message.success('Cập nhật khách hàng thành công');
          onSuccess();
        },
        onError: (err: any) => message.error(err.response?.data?.error || 'Lỗi cập nhật')
      });
    } else {
      createMutation.mutate(payload, {
        onSuccess: () => {
          message.success('Tạo khách hàng mới thành công');
          onSuccess();
        },
        onError: (err: any) => message.error(err.response?.data?.error || 'Lỗi tạo mới')
      });
    }
  };

  return (
    <Form layout="vertical" form={form} onFinish={onFinish}>
      <Row gutter={16} className="mb-4">
        <Col span={24}>
          <div className="flex items-center gap-3 bg-slate-50 p-3 rounded-lg border border-slate-200">
            <span className="font-medium">Loại khách hàng:</span>
            <Switch 
              checkedChildren="Doanh nghiệp (B2B)" 
              unCheckedChildren="Cá nhân (B2C)" 
              checked={isB2B}
              onChange={(checked) => setIsB2B(checked)}
            />
          </div>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Form.Item name="name" label="Tên khách hàng" rules={[{ required: true, message: 'Vui lòng nhập tên' }]}>
            <Input placeholder="Nguyễn Văn A" />
          </Form.Item>
        </Col>
        <Col span={12}>
          <Form.Item name="phone" label="Số điện thoại" rules={[{ required: true, message: 'Vui lòng nhập SĐT' }]}>
            <Input placeholder="0987654321" />
          </Form.Item>
        </Col>
      </Row>

      <Row gutter={16}>
        <Col span={12}>
          <Form.Item name="email" label="Email" rules={[{ type: 'email', message: 'Email không hợp lệ' }]}>
            <Input placeholder="email@example.com" />
          </Form.Item>
        </Col>
        <Col span={6}>
          <Form.Item name="gender" label="Giới tính">
            <Select placeholder="Chọn">
              <Option value="Nam">Nam</Option>
              <Option value="Nữ">Nữ</Option>
              <Option value="Khác">Khác</Option>
            </Select>
          </Form.Item>
        </Col>
        <Col span={6}>
          <Form.Item name="dob" label="Ngày sinh">
            <DatePicker className="w-full" format="DD/MM/YYYY" placeholder="Chọn ngày" />
          </Form.Item>
        </Col>
      </Row>

      {!isB2B && (
        <Row gutter={16}>
          <Col span={12}>
            <Form.Item name="cccd" label="CCCD/CMND">
              <Input placeholder="Nhập số CCCD" />
            </Form.Item>
          </Col>
        </Row>
      )}

      {isB2B && (
        <div className="bg-blue-50/50 p-4 rounded-lg mb-4 border border-blue-100">
          <h4 className="font-semibold text-blue-800 mb-3">Thông tin Doanh nghiệp</h4>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="b2b_tax_code" label="Mã số thuế" rules={[{ required: true, message: 'Vui lòng nhập MST' }]}>
                <Input placeholder="Nhập mã số thuế" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="b2b_company" label="Tên Công ty" rules={[{ required: true, message: 'Vui lòng nhập Tên Cty' }]}>
                <Input placeholder="Công ty TNHH ABC" />
              </Form.Item>
            </Col>
          </Row>
        </div>
      )}

      <Row gutter={16}>
        <Col span={24}>
          <Form.Item name="address" label="Địa chỉ">
            <TextArea rows={2} placeholder="Nhập địa chỉ chi tiết" />
          </Form.Item>
        </Col>
      </Row>

      <div className="flex justify-end gap-2 mt-4">
        <Button onClick={onCancel}>Hủy</Button>
        <Button 
          type="primary" 
          htmlType="submit" 
          loading={createMutation.isPending || updateMutation.isPending}
        >
          {initialData ? 'Cập nhật' : 'Thêm mới'}
        </Button>
      </div>
    </Form>
  );
};
