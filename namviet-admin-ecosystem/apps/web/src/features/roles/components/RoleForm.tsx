import React from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, Button, Checkbox, Card } from 'antd';
import { RolesRole } from '@namviet/shared-types/src/backend.d';
import { useCreateRole, useUpdateRole } from '../hooks';
import { toast } from 'sonner';

const schema = z.object({
  name: z.string().min(2, 'Tên vai trò phải có ít nhất 2 ký tự'),
  description: z.string().optional(),
  permissions: z.array(z.string()).min(1, 'Phải chọn ít nhất 1 quyền'),
});

type FormData = z.infer<typeof schema>;

// Danh sách các quyền hệ thống có thể cấp
const PERMISSION_GROUPS = [
  {
    title: 'Quản lý Sản phẩm',
    options: [
      { label: 'Xem Sản phẩm', value: 'product.read' },
      { label: 'Tạo Sản phẩm', value: 'product.create' },
      { label: 'Sửa Sản phẩm', value: 'product.update' },
      { label: 'Xóa Sản phẩm', value: 'product.delete' },
    ]
  },
  {
    title: 'Quản lý Tồn kho',
    options: [
      { label: 'Xem Tồn kho', value: 'inventory.read' },
      { label: 'Nhập/Xuất kho', value: 'inventory.transaction' },
      { label: 'Quản lý Lô', value: 'batch.manage' },
    ]
  },
  {
    title: 'Quản trị Hệ thống',
    options: [
      { label: 'Quản lý Vai trò', value: 'role.manage' },
      { label: 'Quản lý Nhân sự', value: 'user.manage' },
      { label: 'Quản lý Chi nhánh', value: 'warehouse.manage' },
    ]
  }
];

interface RoleFormProps {
  initialData?: RolesRole | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const RoleForm: React.FC<RoleFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const isAdmin = initialData?.name?.toLowerCase() === 'admin';
  const createMutation = useCreateRole();
  const updateMutation = useUpdateRole();

  const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: initialData?.name || '',
      description: initialData?.description || '',
      permissions: initialData?.permissions || [],
    }
  });

  const onSubmit = (data: FormData) => {
    if (isEditing && initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data },
        {
          onSuccess: () => {
            toast.success('Cập nhật vai trò thành công');
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
            toast.success('Thêm vai trò thành công');
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
            label="Tên vai trò" 
            validateStatus={errors.name ? 'error' : ''}
            help={errors.name?.message}
            required
          >
            <Input {...field} placeholder="VD: Quản lý Kho" disabled={isAdmin} />
          </Form.Item>
        )}
      />

      <Controller
        name="description"
        control={control}
        render={({ field }) => (
          <Form.Item label="Mô tả">
            <Input.TextArea {...field} placeholder="Mô tả vai trò này làm được những gì..." rows={2} />
          </Form.Item>
        )}
      />

      <div className="mb-6">
        <label className="block mb-2 font-medium">
          Phân quyền hệ thống <span className="text-red-500">*</span>
        </label>
        {errors.permissions && (
          <p className="text-red-500 text-sm mb-2">{errors.permissions.message}</p>
        )}
        
        <Controller
          name="permissions"
          control={control}
          render={({ field }) => (
            <div className="flex flex-col gap-4">
              {PERMISSION_GROUPS.map((group) => (
                <Card key={group.title} size="small" title={group.title} className="bg-gray-50">
                  <Checkbox.Group 
                    options={group.options} 
                    value={field.value} 
                    onChange={field.onChange}
                    disabled={isAdmin} // Admin luôn có full quyền, không cho sửa
                    className="flex flex-col gap-2 md:flex-row md:flex-wrap"
                  />
                </Card>
              ))}
            </div>
          )}
        />
        {isAdmin && (
          <p className="text-blue-600 text-xs mt-2 italic">* Vai trò Admin là mặc định của hệ thống, không thể thay đổi tên hoặc giới hạn quyền.</p>
        )}
      </div>

      <div className="flex justify-end gap-2 mt-6">
        <Button onClick={onCancel} disabled={isLoading}>Hủy</Button>
        <Button type="primary" htmlType="submit" loading={isLoading} className="bg-blue-600">
          {isEditing ? 'Lưu thay đổi' : 'Tạo mới'}
        </Button>
      </div>
    </Form>
  );
};
