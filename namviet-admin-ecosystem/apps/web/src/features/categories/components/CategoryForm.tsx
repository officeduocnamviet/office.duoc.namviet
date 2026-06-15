import React, { useMemo } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, Select, Button, TreeSelect } from 'antd';
import { CategoriesCategory } from '@namviet/shared-types/src/backend.d';
import { useCreateCategory, useUpdateCategory, useCategories } from '../hooks';
import { toast } from 'sonner';

const schema = z.object({
  name: z.string().min(2, 'Tên danh mục phải có ít nhất 2 ký tự'),
  slug: z.string().min(2, 'Đường dẫn (slug) bắt buộc'),
  parent_id: z.number().optional(),
  status: z.enum(['ACTIVE', 'INACTIVE']).optional(),
});

type FormData = z.infer<typeof schema>;

interface CategoryFormProps {
  initialData?: CategoriesCategory | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const CategoryForm: React.FC<CategoryFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const createMutation = useCreateCategory();
  const updateMutation = useUpdateCategory();
  const { data: allCategories = [] } = useCategories();

  // Tạo cây danh mục cho TreeSelect
  const treeData = useMemo(() => {
    // Không cho phép chọn chính nó hoặc con của nó làm cha (chống loop)
    // Để đơn giản, chỉ hiển thị danh mục gốc làm cha, hoặc lọc bỏ chính nó.
    const map = new Map();
    const roots: any[] = [];
    
    allCategories.forEach(cat => {
      // Bỏ qua chính nó khi đang edit
      if (isEditing && cat.id === initialData?.id) return;
      map.set(cat.id, { ...cat, value: cat.id, title: cat.name, children: [] });
    });

    allCategories.forEach(cat => {
      if (isEditing && cat.id === initialData?.id) return;
      const node = map.get(cat.id);
      if (cat.parent_id) {
        const parent = map.get(cat.parent_id);
        if (parent) {
          parent.children.push(node);
        }
      } else {
        roots.push(node);
      }
    });

    return roots;
  }, [allCategories, isEditing, initialData]);

  const { control, handleSubmit, formState: { errors }, watch, setValue } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: initialData?.name || '',
      slug: initialData?.slug || '',
      parent_id: initialData?.parent_id || undefined,
      status: (initialData?.status as any) || 'ACTIVE',
    }
  });

  // Tự động tạo slug từ tên
  const generateSlug = (name: string) => {
    return name.toLowerCase()
      .normalize("NFD").replace(/[\u0300-\u036f]/g, "")
      .replace(/[đĐ]/g, 'd')
      .replace(/([^0-9a-z-\s])/g, '')
      .replace(/(\s+)/g, '-')
      .replace(/-+/g, '-')
      .replace(/^-+|-+$/g, '');
  };

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newName = e.target.value;
    setValue('name', newName);
    if (!isEditing) {
      setValue('slug', generateSlug(newName), { shouldValidate: true });
    }
  };

  const onSubmit = (data: FormData) => {
    if (isEditing && initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data },
        {
          onSuccess: () => {
            toast.success('Cập nhật danh mục thành công');
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
            toast.success('Thêm danh mục thành công');
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
        render={({ field: { onChange, ...field } }) => (
          <Form.Item 
            label="Tên danh mục" 
            validateStatus={errors.name ? 'error' : ''}
            help={errors.name?.message}
            required
          >
            <Input {...field} onChange={(e) => { onChange(e); handleNameChange(e); }} placeholder="VD: Thuốc kháng sinh" />
          </Form.Item>
        )}
      />

      <Controller
        name="slug"
        control={control}
        render={({ field }) => (
          <Form.Item 
            label="Đường dẫn (Slug)" 
            validateStatus={errors.slug ? 'error' : ''}
            help={errors.slug?.message}
            required
          >
            <Input {...field} placeholder="thuoc-khang-sinh" />
          </Form.Item>
        )}
      />

      <Controller
        name="parent_id"
        control={control}
        render={({ field }) => (
          <Form.Item label="Danh mục cha (Để trống nếu là danh mục gốc)">
            <TreeSelect
              {...field}
              treeData={treeData}
              placeholder="Chọn danh mục cha"
              allowClear
              treeDefaultExpandAll
            />
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
              { value: 'INACTIVE', label: 'Tạm ẩn' },
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
