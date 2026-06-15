import React from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, DatePicker, Select, Button } from 'antd';
import { BatchesBatch } from '@namviet/shared-types/src/backend.d';
import { useCreateBatch, useUpdateBatch } from '../hooks';
import { useProducts } from '@/features/products/hooks/useProducts';
import { toast } from 'sonner';
import dayjs from 'dayjs';

const schema = z.object({
  product_id: z.number().min(1, 'Vui lòng chọn sản phẩm'),
  batch_code: z.string().min(2, 'Mã lô bắt buộc'),
  manufacturing_date: z.string().optional(),
  expiry_date: z.string().min(1, 'Hạn sử dụng bắt buộc'),
});

type FormData = z.infer<typeof schema>;

interface BatchFormProps {
  initialData?: BatchesBatch | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const BatchForm: React.FC<BatchFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const createMutation = useCreateBatch();
  const updateMutation = useUpdateBatch();
  const { data: products = [], isLoading: loadingProducts } = useProducts();

  const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      product_id: initialData?.product_id || 0,
      batch_code: initialData?.batch_code || '',
      manufacturing_date: initialData?.manufacturing_date ? dayjs(initialData.manufacturing_date).format('YYYY-MM-DD') : undefined,
      expiry_date: initialData?.expiry_date ? dayjs(initialData.expiry_date).format('YYYY-MM-DD') : undefined,
    }
  });

  const onSubmit = (data: FormData) => {
    // Chuyển đổi định dạng ngày nếu cần thiết (để backend hiểu)
    if (isEditing && initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data },
        {
          onSuccess: () => {
            toast.success('Cập nhật lô thành công');
            onSuccess();
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    } else {
      createMutation.mutate(
        data,
        {
          onSuccess: () => {
            toast.success('Thêm lô thành công');
            onSuccess();
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  return (
    <Form layout="vertical" onFinish={handleSubmit(onSubmit)} className="mt-4">
      <Controller
        name="product_id"
        control={control}
        render={({ field }) => (
          <Form.Item 
            label="Sản phẩm"
            validateStatus={errors.product_id ? 'error' : ''}
            help={errors.product_id?.message}
            required
          >
            <Select 
              {...field} 
              placeholder="Chọn sản phẩm"
              loading={loadingProducts}
              showSearch
              filterOption={(input, option) =>
                (option?.label ?? '').toString().toLowerCase().includes(input.toLowerCase())
              }
              options={products.map(p => ({ value: p.id, label: p.name }))}
            />
          </Form.Item>
        )}
      />

      <div className="grid grid-cols-2 gap-4">
        <Controller
          name="batch_code"
          control={control}
          render={({ field }) => (
            <Form.Item 
              label="Mã Lô (Hệ thống)" 
              validateStatus={errors.batch_code ? 'error' : ''}
              help={errors.batch_code?.message}
              required
            >
              <Input {...field} placeholder="VD: LO-2023-001" />
            </Form.Item>
          )}
        />
      </div>
      <div className="grid grid-cols-2 gap-4 mt-4">
        <Controller
          name="manufacturing_date"
          control={control}
          render={({ field }) => (
            <Form.Item label="Ngày sản xuất (NSX)">
              <DatePicker 
                className="w-full" 
                format="DD/MM/YYYY"
                value={field.value ? dayjs(field.value) : undefined}
                onChange={(date, dateString) => field.onChange(date ? date.format('YYYY-MM-DD') : undefined)}
              />
            </Form.Item>
          )}
        />
        <Controller
          name="expiry_date"
          control={control}
          render={({ field }) => (
            <Form.Item label="Hạn sử dụng (HSD)" validateStatus={errors.expiry_date ? 'error' : ''} help={errors.expiry_date?.message} required>
              <DatePicker 
                className="w-full" 
                format="DD/MM/YYYY"
                value={field.value ? dayjs(field.value) : undefined}
                onChange={(date, dateString) => field.onChange(date ? date.format('YYYY-MM-DD') : undefined)}
              />
            </Form.Item>
          )}
        />
      </div>

      <div className="flex justify-end gap-2 mt-6">
        <Button onClick={onCancel} disabled={isLoading}>Hủy</Button>
        <Button type="primary" htmlType="submit" loading={isLoading} className="bg-blue-600">
          {isEditing ? 'Lưu thay đổi' : 'Tạo Lô mới'}
        </Button>
      </div>
    </Form>
  );
};
