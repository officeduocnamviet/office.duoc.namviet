import React, { useState } from 'react';
import { Form, Input, Select, Button, Row, Col, Typography, Card, Space, Divider, message, InputNumber, Table } from 'antd';
import { MinusCircleOutlined, PlusOutlined, ShoppingCartOutlined, UserOutlined } from '@ant-design/icons';
import { useCreateOrder } from '../hooks';
import { OrdersCreateOrderRequest } from '@namviet/shared-types/src/backend.d';
import { useForm, useFieldArray } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Store, UserCircle2 } from 'lucide-react';

const { Text, Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const orderItemSchema = z.object({
  product_id: z.number({ message: 'Chọn sản phẩm' }),
  quantity: z.number().min(1, 'Số lượng tối thiểu là 1'),
  unit_price: z.number().min(0, 'Giá không hợp lệ'),
  uom: z.string().min(1, 'Đơn vị tính'),
  discount: z.number().optional(),
});

const orderSchema = z.object({
  customer_id: z.number().optional(), // Bán lẻ có thể không cần
  note: z.string().optional(),
  items: z.array(orderItemSchema).min(1, 'Đơn hàng phải có ít nhất 1 sản phẩm'),
});

type OrderFormValues = z.infer<typeof orderSchema>;

export const RetailOrderForm = () => {
  const createMutation = useCreateOrder();
  
  const { control, handleSubmit, watch, setValue, formState: { errors } } = useForm<OrderFormValues>({
    resolver: zodResolver(orderSchema),
    defaultValues: {
      items: [{ product_id: 1, quantity: 1, unit_price: 50000, uom: 'Hộp', discount: 0 }],
    }
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: "items"
  });

  const items = watch('items');

  // Tính tổng tiền realtime
  const totalAmount = items.reduce((sum, item) => sum + (item.quantity * item.unit_price) - (item.discount || 0), 0);

  const onSubmit = (data: OrderFormValues) => {
    const payload: OrdersCreateOrderRequest = {
      code: `RET-${Date.now()}`,
      order_type: 'RETAIL',
      customer_id: data.customer_id,
      note: data.note,
      items: data.items.map(i => ({
        product_id: i.product_id,
        quantity: i.quantity,
        unit_price: i.unit_price,
        uom: i.uom,
        discount: i.discount || 0,
      }))
    };

    createMutation.mutate(payload, {
      onSuccess: () => {
        message.success('Tạo đơn hàng Bán lẻ thành công!');
        // Reset form or navigate
      },
      onError: () => {
        message.error('Lỗi khi tạo đơn hàng');
      }
    });
  };

  return (
    <div className="flex flex-col lg:flex-row gap-6">
      {/* Cột trái: Thông tin Chung */}
      <div className="w-full lg:w-1/3 space-y-6">
        <Card title={<><UserCircle2 size={18} className="inline mr-2" />Thông tin Khách hàng</>} className="shadow-sm">
          <div className="mb-4">
            <Text className="text-gray-500 mb-1 block">Khách hàng</Text>
            <Select 
              className="w-full" 
              placeholder="Chọn khách lẻ hoặc tìm kiếm..." 
              allowClear
              onChange={(val) => setValue('customer_id', val)}
            >
              <Option value={1}>Khách vãng lai</Option>
              <Option value={2}>Nguyễn Văn A - 0987654321</Option>
            </Select>
            {errors.customer_id && <Text type="danger" className="text-xs">{errors.customer_id.message}</Text>}
          </div>
          <div>
            <Text className="text-gray-500 mb-1 block">Ghi chú đơn hàng</Text>
            <TextArea 
              rows={3} 
              placeholder="VD: Khách lấy thêm túi ni lông..."
              onChange={(e) => setValue('note', e.target.value)}
            />
          </div>
        </Card>

        <Card title={<><Store size={18} className="inline mr-2" />Thanh toán</>} className="shadow-sm bg-blue-50/50">
          <div className="flex justify-between items-center mb-4">
            <Text className="text-lg">Tổng tiền hàng:</Text>
            <Text strong className="text-xl text-blue-600">
              {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(totalAmount)}
            </Text>
          </div>
          <Button 
            type="primary" 
            size="large" 
            className="w-full bg-blue-600 hover:bg-blue-700"
            onClick={handleSubmit(onSubmit)}
            loading={createMutation.isPending}
          >
            Tạo Đơn Hàng (F9)
          </Button>
        </Card>
      </div>

      {/* Cột phải: Danh sách Sản phẩm */}
      <div className="w-full lg:w-2/3">
        <Card title={<><ShoppingCartOutlined className="mr-2" />Giỏ hàng (Order Items)</>} className="shadow-sm min-h-[500px]">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="border-b border-gray-200 text-sm text-gray-500">
                  <th className="pb-3 w-10">#</th>
                  <th className="pb-3 min-w-[200px]">Sản phẩm</th>
                  <th className="pb-3 w-32">Số lượng</th>
                  <th className="pb-3 w-24">ĐVT</th>
                  <th className="pb-3 w-32">Đơn giá</th>
                  <th className="pb-3 w-32 text-right">Thành tiền</th>
                  <th className="pb-3 w-10"></th>
                </tr>
              </thead>
              <tbody>
                {fields.map((field, index) => (
                  <tr key={field.id} className="border-b border-gray-100 last:border-0 hover:bg-gray-50">
                    <td className="py-3 text-gray-400">{index + 1}</td>
                    <td className="py-3 pr-2">
                      <Select
                        className="w-full"
                        value={watch(`items.${index}.product_id`)}
                        onChange={(val) => setValue(`items.${index}.product_id`, val)}
                        options={[
                          { label: 'Panadol Extra Đỏ', value: 1 },
                          { label: 'Oresol Vị Cam', value: 2 },
                          { label: 'Khẩu trang 4D', value: 3 },
                        ]}
                      />
                    </td>
                    <td className="py-3 pr-2">
                      <InputNumber 
                        min={1} 
                        className="w-full"
                        value={watch(`items.${index}.quantity`)}
                        onChange={(val) => setValue(`items.${index}.quantity`, val || 1)}
                      />
                    </td>
                    <td className="py-3 pr-2">
                      <Input 
                        value={watch(`items.${index}.uom`)}
                        onChange={(e) => setValue(`items.${index}.uom`, e.target.value)}
                      />
                    </td>
                    <td className="py-3 pr-2">
                      <InputNumber 
                        className="w-full"
                        formatter={value => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                        value={watch(`items.${index}.unit_price`)}
                        onChange={(val) => setValue(`items.${index}.unit_price`, val || 0)}
                      />
                    </td>
                    <td className="py-3 text-right font-medium text-blue-600">
                      {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(
                        watch(`items.${index}.quantity`) * watch(`items.${index}.unit_price`)
                      )}
                    </td>
                    <td className="py-3 pl-2">
                      <Button 
                        type="text" 
                        danger 
                        icon={<MinusCircleOutlined />} 
                        onClick={() => remove(index)} 
                      />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          
          <Button 
            type="dashed" 
            onClick={() => append({ product_id: 2, quantity: 1, unit_price: 0, uom: 'Hộp', discount: 0 })} 
            block 
            icon={<PlusOutlined />} 
            className="mt-4 border-gray-300 text-gray-500 hover:text-blue-500 hover:border-blue-500"
          >
            Thêm dòng sản phẩm
          </Button>
          {errors.items && <Text type="danger" className="block mt-2">{errors.items.message}</Text>}
        </Card>
      </div>
    </div>
  );
};
