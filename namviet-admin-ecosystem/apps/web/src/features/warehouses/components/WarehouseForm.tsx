import React from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, Select, Button, Space } from 'antd';
import { WarehousesWarehouse } from '@namviet/shared-types/src/backend.d';
import { useCreateWarehouse, useUpdateWarehouse } from '../hooks';
import { toast } from 'sonner';

const schema = z.object({
  name: z.string().min(2, 'Tên chi nhánh phải có ít nhất 2 ký tự'),
  type: z.enum(['MAIN', 'BRANCH', 'STORE']),
  key: z.string().min(1, 'Key là bắt buộc'),
  unit: z.string().optional(),
  latitude: z.any().optional(),
  longitude: z.any().optional(),
  code: z.string().optional(),
  phone: z.string().optional(),
  address: z.string().optional(),
  manager: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

interface WarehouseFormProps {
  initialData?: WarehousesWarehouse | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const WarehouseForm: React.FC<WarehouseFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const createMutation = useCreateWarehouse();
  const updateMutation = useUpdateWarehouse();

  const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: initialData?.name || '',
      type: (initialData?.type as any) || 'BRANCH',
      key: (initialData as any)?.key || '',
      unit: (initialData as any)?.unit || '',
      latitude: (initialData as any)?.latitude || undefined,
      longitude: (initialData as any)?.longitude || undefined,
      code: (initialData as any)?.code || '',
      phone: (initialData as any)?.phone || '',
      address: initialData?.address || '',
      manager: initialData?.manager || '',
    }
  });

  const onSubmit = (data: FormData) => {
    if (isEditing && initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data },
        {
          onSuccess: () => {
            toast.success('Cập nhật chi nhánh thành công');
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
            toast.success('Thêm chi nhánh thành công');
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
            label="Tên chi nhánh" 
            validateStatus={errors.name ? 'error' : ''}
            help={errors.name?.message}
            required
          >
            <Input {...field} placeholder="VD: Kho Tổng Hà Nội" />
          </Form.Item>
        )}
      />

      <div className="grid grid-cols-2 gap-4">
        <Controller
          name="type"
          control={control}
          render={({ field }) => (
            <Form.Item label="Phân loại">
              <Select {...field} options={[
                { value: 'MAIN', label: 'Tổng Kho' },
                { value: 'BRANCH', label: 'Chi nhánh' },
                { value: 'STORE', label: 'Cửa hàng' },
              ]} />
            </Form.Item>
          )}
        />

        <Controller
          name="key"
          control={control}
          render={({ field }) => (
            <Form.Item 
              label="Mã định danh (Key)" 
              validateStatus={errors.key ? 'error' : ''}
              help={errors.key?.message}
              required
            >
              <Input {...field} placeholder="VD: KHO_HN" />
            </Form.Item>
          )}
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Controller
          name="code"
          control={control}
          render={({ field }) => (
            <Form.Item label="Mã số (Code)">
              <Input {...field} placeholder="VD: WH-001" />
            </Form.Item>
          )}
        />

        <Controller
          name="unit"
          control={control}
          render={({ field }) => (
            <Form.Item label="Đơn vị quản lý">
              <Input {...field} placeholder="VD: Chi nhánh MB" />
            </Form.Item>
          )}
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Controller
          name="phone"
          control={control}
          render={({ field }) => (
            <Form.Item label="Số điện thoại">
              <Input {...field} placeholder="0987654321" />
            </Form.Item>
          )}
        />
        <Controller
          name="manager"
          control={control}
          render={({ field }) => (
            <Form.Item label="Người quản lý">
              <Input {...field} placeholder="Nguyễn Văn A" />
            </Form.Item>
          )}
        />
      </div>

      <div className="grid grid-cols-2 gap-4">
        <Controller
          name="latitude"
          control={control}
          render={({ field }) => (
            <Form.Item label="Vĩ độ (Latitude)">
              <Input type="number" step="any" {...field} placeholder="21.0285" />
            </Form.Item>
          )}
        />

        <Controller
          name="longitude"
          control={control}
          render={({ field }) => (
            <Form.Item label="Kinh độ (Longitude)">
              <Input type="number" step="any" {...field} placeholder="105.8542" />
            </Form.Item>
          )}
        />
      </div>

      <Controller
        name="address"
        control={control}
        render={({ field }) => (
          <Form.Item label="Địa chỉ">
            <Input.TextArea {...field} placeholder="Địa chỉ chi tiết..." rows={3} />
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
