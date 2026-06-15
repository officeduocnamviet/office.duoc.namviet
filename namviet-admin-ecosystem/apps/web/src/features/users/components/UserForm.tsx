import React from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, Select, Button } from 'antd';
import { UsersUser } from '@namviet/shared-types/src/backend.d';
import { useCreateUser, useUpdateUser } from '../hooks';
import { useRoles } from '@/features/roles/hooks';
import { useWarehouses } from '@/features/warehouses/hooks';
import { toast } from 'sonner';

const schema = z.object({
  full_name: z.string().min(2, 'Tên nhân viên phải có ít nhất 2 ký tự'),
  email: z.string().email('Email không hợp lệ'),
  phone: z.string().optional(),
  password: z.string().min(6, 'Mật khẩu phải từ 6 ký tự').optional(),
  role_id: z.string().min(1, 'Vui lòng chọn vai trò'),
  warehouse_id: z.number().optional().nullable(),
  status: z.enum(['active', 'inactive', 'working']).optional(),
});

type FormData = z.infer<typeof schema>;

interface UserFormProps {
  initialData?: UsersUser | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const UserForm: React.FC<UserFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const createMutation = useCreateUser();
  const updateMutation = useUpdateUser();

  // Fetch dropdown data
  const { data: roles = [], isLoading: loadingRoles } = useRoles();
  const { data: warehouses = [], isLoading: loadingWarehouses } = useWarehouses();

  const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      full_name: initialData?.full_name || '',
      email: initialData?.email || '',
      phone: initialData?.phone || '',
      password: '', // Không hiện password cũ
      role_id: initialData?.role_id || '',
      warehouse_id: initialData?.warehouse_id || null,
      status: (initialData?.status as any) || 'active',
    }
  });

  const onSubmit = (data: FormData) => {
    if (isEditing && initialData?.id) {
      // Bỏ qua password nếu không nhập khi update
      const updateData = { ...data };
      if (!updateData.password) delete updateData.password;

      updateMutation.mutate(
        { id: initialData.id, data: updateData as any },
        {
          onSuccess: () => {
            toast.success('Cập nhật nhân sự thành công');
            onSuccess();
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    } else {
      if (!data.password) {
        toast.error('Vui lòng nhập mật khẩu cho nhân sự mới');
        return;
      }
      createMutation.mutate(
        { ...data, password: data.password, company_id: '1' } as any,
        {
          onSuccess: () => {
            toast.success('Thêm nhân sự thành công');
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
      <div className="grid grid-cols-2 gap-4">
        <Controller
          name="full_name"
          control={control}
          render={({ field }) => (
            <Form.Item 
              label="Họ và tên" 
              validateStatus={errors.full_name ? 'error' : ''}
              help={errors.full_name?.message}
              required
            >
              <Input {...field} placeholder="VD: Nguyễn Văn A" />
            </Form.Item>
          )}
        />

        <Controller
          name="email"
          control={control}
          render={({ field }) => (
            <Form.Item 
              label="Email đăng nhập" 
              validateStatus={errors.email ? 'error' : ''}
              help={errors.email?.message}
              required
            >
              <Input {...field} placeholder="nguyenvana@namviet.com" disabled={isEditing} />
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
              <Input {...field} placeholder="09xxxx..." />
            </Form.Item>
          )}
        />

        <Controller
          name="password"
          control={control}
          render={({ field }) => (
            <Form.Item 
              label={isEditing ? "Mật khẩu mới (Để trống nếu không đổi)" : "Mật khẩu khởi tạo"}
              validateStatus={errors.password ? 'error' : ''}
              help={errors.password?.message}
              required={!isEditing}
            >
              <Input.Password {...field} placeholder="******" />
            </Form.Item>
          )}
        />
      </div>

      <div className="border-t border-gray-100 my-4 pt-4">
        <h4 className="text-sm font-semibold mb-4 text-gray-700">Phân công công tác</h4>
        
        <div className="grid grid-cols-2 gap-4">
          <Controller
            name="role_id"
            control={control}
            render={({ field }) => (
              <Form.Item 
                label="Vai trò (Phân quyền)"
                validateStatus={errors.role_id ? 'error' : ''}
                help={errors.role_id?.message}
                required
              >
                <Select 
                  {...field} 
                  placeholder="Chọn vai trò"
                  loading={loadingRoles}
                  options={roles.map(r => ({ value: r.id, label: r.name }))}
                />
              </Form.Item>
            )}
          />

          <Controller
            name="warehouse_id"
            control={control}
            render={({ field }) => (
              <Form.Item label="Chi nhánh làm việc (Kho)">
                <Select 
                  {...field} 
                  placeholder="Chọn chi nhánh"
                  allowClear
                  loading={loadingWarehouses}
                  options={warehouses.map(w => ({ value: w.id, label: w.name }))}
                />
              </Form.Item>
            )}
          />
        </div>
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
