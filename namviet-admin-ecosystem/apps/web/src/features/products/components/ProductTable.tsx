import React, { useState } from 'react';
import { Table, Button, Space, Modal, Popconfirm, Tag, Input, Drawer, Select, Divider } from 'antd';
import { Plus, Pencil, Trash2, Search, Filter, Image as ImageIcon, Box } from 'lucide-react';
import { ProductsProduct } from '@namviet/shared-types/src/backend.d';
import { useProducts, useDeleteProduct } from '../hooks/useProducts';
import { toast } from 'sonner';
import { ProductForm } from './ProductForm';

export const ProductTable = () => {
  const { data: products, isLoading } = useProducts();
  const deleteMutation = useDeleteProduct();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingProduct, setEditingProduct] = useState<ProductsProduct | null>(null);
  const [searchText, setSearchText] = useState('');

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa sản phẩm thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const filteredProducts = products?.filter(p => 
    p.name?.toLowerCase().includes(searchText.toLowerCase()) ||
    p.sku?.toLowerCase().includes(searchText.toLowerCase()) ||
    p.barcode?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Ảnh',
      dataIndex: 'image_url',
      key: 'image',
      width: 80,
      render: (url: string) => (
        <div className="w-12 h-12 bg-gray-100 rounded flex items-center justify-center border border-gray-200 overflow-hidden">
          {url ? <img src={url} alt="product" className="w-full h-full object-cover" /> : <ImageIcon className="w-5 h-5 text-gray-400" />}
        </div>
      ),
    },
    {
      title: 'Sản phẩm',
      key: 'product',
      render: (_: any, record: ProductsProduct) => (
        <div>
          <div className="font-medium text-gray-900">{record.name}</div>
          <div className="text-xs text-gray-500">SKU: {record.sku || '---'} | Mã vạch: {record.barcode || '---'}</div>
          {record.active_ingredient && (
            <div className="text-xs text-blue-600 mt-1">Hoạt chất: {record.active_ingredient}</div>
          )}
        </div>
      ),
    },
    {
      title: 'Danh mục',
      dataIndex: 'category_name',
      key: 'category',
    },
    {
      title: 'Nhà SX',
      dataIndex: 'manufacturer_name',
      key: 'manufacturer',
    },
    {
      title: 'Giá bán lẻ',
      dataIndex: 'price_sell',
      key: 'price_sell',
      render: (price: number, record: ProductsProduct) => (
        <div>
          <div className="font-medium text-green-600">{price ? price.toLocaleString('vi-VN') + ' đ' : '---'}</div>
          <div className="text-xs text-gray-500">/{record.retail_unit || 'Đơn vị'}</div>
        </div>
      ),
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'ACTIVE' ? 'success' : 'default'}>
          {status === 'ACTIVE' ? 'Đang bán' : 'Ngừng kinh doanh'}
        </Tag>
      ),
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: ProductsProduct) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingProduct(record);
              setIsModalOpen(true);
            }}
          />
          <Popconfirm
            title="Xóa sản phẩm?"
            description="Bạn có chắc chắn muốn xóa sản phẩm này không?"
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
          <h2 className="text-lg font-semibold text-gray-900">Quản lý Sản phẩm</h2>
          <p className="text-sm text-gray-500">Danh sách các mặt hàng dược phẩm, vật tư y tế</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            prefix={<Search className="w-4 h-4 text-gray-400" />} 
            placeholder="Tìm theo tên, SKU, mã vạch..." 
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
              setEditingProduct(null);
              setIsModalOpen(true);
            }}
          >
            Thêm
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        {/* Desktop View: Table */}
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredProducts} 
            rowKey="id"
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        {/* Mobile View: List Cards */}
        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải dữ liệu...</div>
          ) : filteredProducts?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không tìm thấy sản phẩm nào.</div>
          ) : (
            filteredProducts?.map((product) => (
              <div key={product.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3">
                  <div className="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center border border-gray-200 flex-shrink-0 overflow-hidden">
                    {product.image_url ? (
                      <img src={product.image_url} alt="product" className="w-full h-full object-cover" />
                    ) : (
                      <ImageIcon className="w-6 h-6 text-gray-300" />
                    )}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-medium text-gray-900 text-sm line-clamp-2">{product.name}</h3>
                    <div className="text-xs text-gray-500 mt-1 truncate">SKU: {product.sku || '---'}</div>
                    <div className="mt-1 font-semibold text-green-600 text-sm">
                      {product.price_sell ? product.price_sell.toLocaleString('vi-VN') + 'đ' : '---'}
                      <span className="text-xs text-gray-400 font-normal">/{product.retail_unit || 'ĐV'}</span>
                    </div>
                  </div>
                </div>
                
                {/* Mobile Actions Bottom Bar */}
                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                  <Button 
                    type="text" 
                    size="small"
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingProduct(product);
                      setIsModalOpen(true);
                    }}
                  >
                    Sửa
                  </Button>
                  <Popconfirm
                    title="Xóa SP?"
                    onConfirm={() => product.id && handleDelete(product.id)}
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

      <Modal
        title={editingProduct ? "Cập nhật Sản phẩm" : "Thêm Sản phẩm mới"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
        width={900}
        style={{ top: 20 }}
      >
        <ProductForm 
          initialData={editingProduct} 
          onSuccess={() => setIsModalOpen(false)}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>

      <Drawer
        title={<span className="font-semibold text-gray-800">Lọc & Tìm kiếm</span>}
        placement="right"
        onClose={() => setIsFilterDrawerOpen(false)}
        open={isFilterDrawerOpen}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 500}
        styles={{ body: { padding: '24px' } }}
      >
        <div className="space-y-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả trạng thái" 
              allowClear
              options={[
                { value: 'ACTIVE', label: 'Đang kinh doanh' },
                { value: 'INACTIVE', label: 'Ngừng kinh doanh' }
              ]} 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Danh mục</label>
            <Select 
              className="w-full" 
              placeholder="Chọn danh mục" 
              allowClear 
            />
          </div>
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
