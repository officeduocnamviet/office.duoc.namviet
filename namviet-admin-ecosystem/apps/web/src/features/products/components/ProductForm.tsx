import React, { useState } from 'react';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Form, Input, InputNumber, Select, Button, Tabs } from 'antd';
import { ProductsProduct } from '@namviet/shared-types/src/backend.d';
import { useCreateProduct, useUpdateProduct } from '../hooks/useProducts';
import { useCategories } from '@/features/categories/hooks';
import { useManufacturers } from '@/features/manufacturers/hooks';
import { toast } from 'sonner';
import { ProductUnitSection } from './ProductUnitSection';

const schema = z.object({
  name: z.string().min(2, 'Tên sản phẩm phải có ít nhất 2 ký tự'),
  sku: z.string().optional(),
  barcode: z.string().optional(),
  category_id: z.number().optional(),
  manufacturer_id: z.number().optional(),
  active_ingredient: z.string().optional(),
  price_sell: z.number().min(0, 'Giá bán không được âm').optional(),
  price_cost: z.number().min(0, 'Giá vốn không được âm').optional(),
  retail_unit: z.string().optional(),
  status: z.enum(['ACTIVE', 'INACTIVE']).optional(),
  description: z.string().optional(),
});

type FormData = z.infer<typeof schema>;

interface ProductFormProps {
  initialData?: ProductsProduct | null;
  onSuccess: () => void;
  onCancel: () => void;
}

export const ProductForm: React.FC<ProductFormProps> = ({ initialData, onSuccess, onCancel }) => {
  const isEditing = !!initialData;
  const [activeTab, setActiveTab] = useState('1');
  
  const createMutation = useCreateProduct();
  const updateMutation = useUpdateProduct();

  // Load select options
  const { data: categories = [], isLoading: loadingCategories } = useCategories();
  const { data: manufacturers = [], isLoading: loadingManufacturers } = useManufacturers();

  const { control, handleSubmit, formState: { errors } } = useForm<FormData>({
    resolver: zodResolver(schema),
    defaultValues: {
      name: initialData?.name || '',
      sku: initialData?.sku || '',
      barcode: initialData?.barcode || '',
      category_id: initialData?.category_id || undefined,
      manufacturer_id: initialData?.manufacturer_id || undefined,
      active_ingredient: initialData?.active_ingredient || '',
      price_sell: initialData?.price_sell || 0,
      price_cost: initialData?.price_cost || 0,
      retail_unit: initialData?.retail_unit || '',
      status: (initialData?.status as any) || 'ACTIVE',
      description: initialData?.description || '',
    }
  });

  const onSubmit = (data: FormData) => {
    if (isEditing && initialData?.id) {
      updateMutation.mutate(
        { id: initialData.id, data },
        {
          onSuccess: () => {
            toast.success('Cập nhật sản phẩm thành công');
            // Cập nhật thành công không cần đóng modal ngay nếu đang sửa tab
            // Nhưng có thể gọi onSuccess để tải lại hoặc đóng modal.
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
            toast.success('Thêm sản phẩm thành công');
            onSuccess();
          },
          onError: (err) => toast.error(`Lỗi: ${err.message}`)
        }
      );
    }
  };

  const isLoading = createMutation.isPending || updateMutation.isPending;

  const basicInfoContent = (
    <Form layout="vertical" onFinish={handleSubmit(onSubmit)} className="mt-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Controller
          name="name"
          control={control}
          render={({ field }) => (
            <Form.Item 
              label="Tên sản phẩm" 
              validateStatus={errors.name ? 'error' : ''}
              help={errors.name?.message}
              required
            >
              <Input {...field} placeholder="Tên hàng hóa, thuốc, vật tư..." />
            </Form.Item>
          )}
        />
        <Controller
          name="active_ingredient"
          control={control}
          render={({ field }) => (
            <Form.Item label="Hoạt chất chính">
              <Input {...field} placeholder="VD: Paracetamol 500mg" />
            </Form.Item>
          )}
        />
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <Controller
          name="sku"
          control={control}
          render={({ field }) => (
            <Form.Item label="Mã SKU nội bộ">
              <Input {...field} placeholder="SKU-XXXX" />
            </Form.Item>
          )}
        />
        <Controller
          name="barcode"
          control={control}
          render={({ field }) => (
            <Form.Item label="Mã vạch (Barcode)">
              <Input {...field} placeholder="Quét mã vạch..." />
            </Form.Item>
          )}
        />
        <Controller
          name="status"
          control={control}
          render={({ field }) => (
            <Form.Item label="Trạng thái">
              <Select {...field} options={[
                { value: 'ACTIVE', label: 'Đang kinh doanh' },
                { value: 'INACTIVE', label: 'Ngừng kinh doanh' },
              ]} />
            </Form.Item>
          )}
        />
      </div>

      <div className="border-t border-gray-100 my-4 pt-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Controller
            name="category_id"
            control={control}
            render={({ field }) => (
              <Form.Item label="Danh mục">
                <Select 
                  {...field} 
                  placeholder="Chọn danh mục"
                  loading={loadingCategories}
                  allowClear
                  options={categories.map(c => ({ value: c.id, label: c.name }))}
                />
              </Form.Item>
            )}
          />
          <Controller
            name="manufacturer_id"
            control={control}
            render={({ field }) => (
              <Form.Item label="Nhà sản xuất">
                <Select 
                  {...field} 
                  placeholder="Chọn NSX"
                  loading={loadingManufacturers}
                  allowClear
                  options={manufacturers.map(m => ({ value: m.id, label: m.name }))}
                />
              </Form.Item>
            )}
          />
        </div>
      </div>

      <div className="border-t border-gray-100 my-4 pt-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Controller
            name="retail_unit"
            control={control}
            render={({ field }) => (
              <Form.Item label="Đơn vị bán lẻ cơ bản">
                <Input {...field} placeholder="Viên, Chai, Gói..." />
              </Form.Item>
            )}
          />
          <Controller
            name="price_sell"
            control={control}
            render={({ field }) => (
              <Form.Item label="Giá bán lẻ (VNĐ)" help="Giá bán cho 1 Đơn vị cơ bản">
                <InputNumber {...field} className="w-full" formatter={value => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')} parser={value => value?.replace(/\$\s?|(,*)/g, '') as any} />
              </Form.Item>
            )}
          />
          <Controller
            name="price_cost"
            control={control}
            render={({ field }) => (
              <Form.Item label="Giá vốn tham khảo">
                <InputNumber {...field} className="w-full" formatter={value => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')} parser={value => value?.replace(/\$\s?|(,*)/g, '') as any} />
              </Form.Item>
            )}
          />
        </div>
      </div>

      <div className="flex justify-end gap-2 mt-6 border-t border-gray-100 pt-4">
        <Button onClick={onCancel} disabled={isLoading}>Hủy</Button>
        <Button type="primary" htmlType="submit" loading={isLoading} className="bg-blue-600">
          {isEditing ? 'Lưu thông tin chung' : 'Tạo sản phẩm'}
        </Button>
      </div>
    </Form>
  );

  return (
    <Tabs 
      activeKey={activeTab} 
      onChange={setActiveTab}
      items={[
        {
          key: '1',
          label: 'Thông tin chung',
          children: basicInfoContent,
        },
        {
          key: '2',
          label: 'Đơn vị tính & Quy đổi',
          disabled: !isEditing, // Chỉ cho thêm đơn vị quy đổi khi đã lưu SP gốc
          children: isEditing && initialData?.id ? <ProductUnitSection productId={initialData.id} baseUnit={initialData.retail_unit || 'Đơn vị'} /> : null,
        }
      ]}
    />
  );
};
