import React from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, Select, Button } from 'antd';
import { ManufacturersManufacturer } from '@namviet/shared-types/src/backend.d';
import { useCreateManufacturer, useUpdateManufacturer } from '../hooks';
import { toast } from 'sonner';

const schema = z.object({
  name: z.string().min(2, 'Tên nhà sản xuất phải có ít nhất 2 ký tự'),
  country: z.string().optional(),
  status: z.enum(['ACTIVE', 'INACTIVE']).optional(),
});

type FormData = z.infer<typeof schema>;

interface ManufacturerFormProps {
  initialData?: ManufacturersManufacturer | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const ManufacturerForm: React.FC<ManufacturerFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const createMutation = useCreateManufacturer();
  const updateMutation = useUpdateManufacturer();

  const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: initialData?.name || '',
      country: initialData?.country || '',
      status: (initialData?.status as any) || 'ACTIVE',
    }
  });

  const onSubmit = (data: FormData) => {
    if (isEditing && initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data },
        {
          onSuccess: () => {
            toast.success('Cập nhật NSX thành công');
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
            toast.success('Thêm NSX thành công');
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
        name="name"
        control={control}
        render={({ field }) => (
          <Form.Item 
            label="Tên nhà sản xuất" 
            validateStatus={errors.name ? 'error' : ''}
            help={errors.name?.message}
            required
          >
            <Input {...field} placeholder="VD: Công ty Cổ phần Dược phẩm Nam Việt" />
          </Form.Item>
        )}
      />

      <Controller
        name="country"
        control={control}
        render={({ field }) => (
          <Form.Item label="Quốc gia">
            <Input {...field} placeholder="VD: Việt Nam, USA, Đức..." />
          </Form.Item>
        )}
      />

      <Controller
        name="status"
        control={control}
        render={({ field }) => (
          <Form.Item label="Trạng thái">
            <Select {...field} options={[
              { value: 'ACTIVE', label: 'Hoạt động' },
              { value: 'INACTIVE', label: 'Ngừng hợp tác' },
            ]} />
          </Form.Item>
        )}
      />

      <div className="flex justify-end gap-2 mt-6">
        <Button onClick={onCancel} disabled={isLoading}>Hủy</Button>
        <Button type="primary" htmlType="submit" loading={isLoading} className="bg-blue-600">
          {isEditing ? 'Lưu thay đổi' : 'Tạo mới'}
        </Button>
      </div>
    </Form>
  );
};
