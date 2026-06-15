import React, { useState, useMemo } from 'react';
import { Table, Button, Space, Drawer, Popconfirm, Tag, Input, Divider } from 'antd';
import { Plus, Pencil, Trash2, FolderTree, Search, Filter } from 'lucide-react';
import { CategoriesCategory } from '@namviet/shared-types/src/backend.d';
import { useCategories, useDeleteCategory } from '../hooks';
import { toast } from 'sonner';
import { CategoryForm } from './CategoryForm';

export const CategoryTable = () => {
  const { data: categories, isLoading } = useCategories();
  const deleteMutation = useDeleteCategory();

  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingCategory, setEditingCategory] = useState<CategoriesCategory | null>(null);
  const [searchText, setSearchText] = useState('');

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa danh mục thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  // Convert flat array to tree structure for Antd Table
  const treeData = useMemo(() => {
    if (!categories) return [];
    const map = new Map();
    const roots: any[] = [];
    
    categories.forEach(cat => {
      map.set(cat.id, { ...cat, key: cat.id, children: [] });
    });

    categories.forEach(cat => {
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

    // Remove empty children arrays so Antd doesn't show expand icon unnecessarily
    const cleanEmptyChildren = (nodes: any[]) => {
      nodes.forEach(node => {
        if (node.children.length === 0) {
          delete node.children;
        } else {
          cleanEmptyChildren(node.children);
        }
      });
    };
    cleanEmptyChildren(roots);

    return roots;
  }, [categories]);

  const columns = [
    {
      title: 'Tên danh mục',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <div className="flex items-center gap-2 font-medium text-gray-900">
          <FolderTree className="w-4 h-4 text-blue-500" />
          {text}
        </div>
      ),
    },
    {
      title: 'Đường dẫn (Slug)',
      dataIndex: 'slug',
      key: 'slug',
      render: (slug: string) => <span className="text-gray-500">{slug}</span>
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'ACTIVE' ? 'success' : 'default'}>
          {status === 'ACTIVE' ? 'Hoạt động' : 'Tạm ẩn'}
        </Tag>
      ),
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: CategoriesCategory) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingCategory(record);
              setIsDrawerOpen(true);
            }}
          />
          <Popconfirm
            title="Xóa danh mục?"
            description="Lưu ý: Bạn không thể xóa danh mục nếu đang có danh mục con hoặc sản phẩm phụ thuộc."
            onConfirm={() => record.id && handleDelete(record.id)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger icon={<Trash2 className="w-4 h-4" />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Danh mục Sản phẩm</h2>
          <p className="text-sm text-gray-500">Quản lý phân loại, nhóm các mặt hàng</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            prefix={<Search className="w-4 h-4 text-gray-400" />} 
            placeholder="Tìm danh mục..." 
            className="w-full sm:w-64"
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
          />
          <Button 
            icon={<Filter className="w-4 h-4" />}
            onClick={() => setIsFilterDrawerOpen(true)}
          >
            Lọc
          </Button>
          <Button 
            type="primary" 
            icon={<Plus className="w-4 h-4" />}
            className="bg-blue-600"
            onClick={() => {
              setEditingCategory(null);
              setIsDrawerOpen(true);
            }}
          >
            Thêm
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={treeData} 
            rowKey="id"
            loading={isLoading}
            pagination={false}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : categories?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có danh mục nào.</div>
          ) : (
            categories?.filter(c => c.name?.toLowerCase().includes(searchText.toLowerCase())).map((cat) => (
              <div key={cat.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-center">
                  <div className="w-10 h-10 bg-blue-50 rounded-lg flex items-center justify-center text-blue-500 flex-shrink-0">
                    <FolderTree className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-gray-900 text-sm truncate">{cat.name}</h3>
                    <div className="text-xs text-gray-500 mt-1 truncate">Slug: {cat.slug}</div>
                  </div>
                  <Tag color={cat.status === 'ACTIVE' ? 'success' : 'default'} className="m-0">
                    {cat.status === 'ACTIVE' ? 'HĐ' : 'Ẩn'}
                  </Tag>
                </div>
                
                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                  <Button 
                    type="text" 
                    size="small"
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingCategory(cat);
                      setIsDrawerOpen(true);
                    }}
                  >
                    Sửa
                  </Button>
                  <Popconfirm
                    title="Xóa?"
                    onConfirm={() => cat.id && handleDelete(cat.id)}
                    okText="Xóa"
                    cancelText="Hủy"
                  >
                    <Button 
                      type="text" 
                      size="small"
                      danger
                      className="bg-red-50 hover:bg-red-100 px-3"
                      icon={<Trash2 className="w-3.5 h-3.5" />} 
                    >
                      Xóa
                    </Button>
                  </Popconfirm>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      <Drawer
        title={editingCategory ? "Cập nhật Danh mục" : "Thêm Danh mục mới"}
        width={450}
        onClose={() => setIsDrawerOpen(false)}
        open={isDrawerOpen}
        destroyOnClose
      >
        <CategoryForm 
          initialData={editingCategory} 
          onSuccess={() => setIsDrawerOpen(false)}
          onCancel={() => setIsDrawerOpen(false)}
        />
      </Drawer>

      <Drawer
        title={<span className="font-semibold text-gray-800">Lọc & Tìm kiếm</span>}
        placement="right"
        onClose={() => setIsFilterDrawerOpen(false)}
        open={isFilterDrawerOpen}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 500}
        styles={{ body: { padding: '24px' } }}
      >
        <div className="space-y-6">
          <p className="text-gray-500">Bộ lọc danh mục sẽ được cập nhật thêm.</p>
          <Divider />
          <div className="flex justify-end gap-2">
            <Button onClick={() => setIsFilterDrawerOpen(false)}>Hủy</Button>
            <Button type="primary" className="bg-blue-600" onClick={() => setIsFilterDrawerOpen(false)}>Áp dụng</Button>
          </div>
        </div>
      </Drawer>
    </div>
  );
};
