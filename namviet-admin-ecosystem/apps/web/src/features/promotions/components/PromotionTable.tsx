import React, { useState } from 'react';
import { Table, Button, Tag, Space, Popconfirm, Drawer, Input, Typography, Select, Divider } from 'antd';
import { Plus, Pencil, Trash2, Search, Ticket, Filter } from 'lucide-react';
import { usePromotions, useDeletePromotion } from '../hooks';
import { PromotionForm } from './PromotionForm';
import { PromotionsPromotion } from '@namviet/shared-types/src/backend.d';
import dayjs from 'dayjs';

const { Text } = Typography;

export const PromotionTable = () => {
  const { data: promotions = [], isLoading } = usePromotions();
  const deleteMutation = useDeletePromotion();
  
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingPromotion, setEditingPromotion] = useState<PromotionsPromotion | undefined>();
  const [searchText, setSearchText] = useState('');

  const filteredPromotions = promotions.filter(p => 
    p.name?.toLowerCase().includes(searchText.toLowerCase()) ||
    p.code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã Code',
      dataIndex: 'code',
      key: 'code',
      width: 150,
      render: (code: string) => (
        <div className="flex items-center gap-2">
          <Ticket size={16} className="text-orange-500" />
          <Text strong className="text-orange-600 font-mono text-base">{code}</Text>
        </div>
      )
    },
    {
      title: 'Tên chương trình',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => <div className="font-medium text-gray-800">{name}</div>
    },
    {
      title: 'Thời hạn',
      key: 'dates',
      render: (_: any, record: PromotionsPromotion) => {
        const start = dayjs(record.start_date);
        const end = dayjs(record.end_date);
        const isExpired = end.isBefore(dayjs());
        
        return (
          <div>
            <div className="text-sm">Từ: {start.format('DD/MM/YYYY HH:mm')}</div>
            <div className="text-sm">Đến: {end.format('DD/MM/YYYY HH:mm')}</div>
            {isExpired && <Tag color="red" className="mt-1">Đã hết hạn</Tag>}
          </div>
        );
      }
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'ACTIVE' 
          ? <Tag color="success">Đang kích hoạt</Tag>
          : <Tag color="default">Tạm dừng</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 120,
      render: (_: any, record: PromotionsPromotion) => (
        <Space size="middle">
          <Button 
            type="text" 
            icon={<Pencil size={16} className="text-blue-600" />} 
            onClick={() => {
              setEditingPromotion(record);
              setDrawerVisible(true);
            }}
          />
          <Popconfirm
            title="Xóa voucher này?"
            onConfirm={() => deleteMutation.mutate(record.id!)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger icon={<Trash2 size={16} />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Voucher & Khuyến mãi</h2>
          <p className="text-sm text-gray-500">Quản lý các chương trình ưu đãi, giảm giá</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm theo mã, tên..." 
            prefix={<Search className="w-4 h-4 text-gray-400" />}
            value={searchText}
            onChange={e => setSearchText(e.target.value)}
            className="w-full sm:w-64"
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
            onClick={() => {
              setEditingPromotion(undefined);
              setDrawerVisible(true);
            }}
            className="bg-orange-500 hover:bg-orange-600 flex items-center border-none"
          >
            Tạo Voucher
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredPromotions} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredPromotions?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có voucher nào.</div>
          ) : (
            filteredPromotions?.map((promo) => {
              const start = dayjs(promo.start_date);
              const end = dayjs(promo.end_date);
              const isExpired = end.isBefore(dayjs());
              
              return (
                <div key={promo.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                  <div className="flex gap-3 items-center">
                    <div className="w-10 h-10 bg-orange-50 rounded-lg flex items-center justify-center text-orange-500 flex-shrink-0">
                      <Ticket className="w-5 h-5" />
                    </div>
                    <div className="flex-1 min-w-0">
                      <h3 className="font-semibold text-orange-600 text-sm truncate font-mono">{promo.code}</h3>
                      <div className="text-xs text-gray-700 font-medium truncate mt-1">{promo.name}</div>
                    </div>
                    {promo.status === 'ACTIVE' && !isExpired ? (
                      <Tag color="success" className="m-0">HĐ</Tag>
                    ) : isExpired ? (
                      <Tag color="error" className="m-0">Hết hạn</Tag>
                    ) : (
                      <Tag color="default" className="m-0">Tạm dừng</Tag>
                    )}
                  </div>
                  
                  <div className="bg-gray-50 p-2 rounded text-xs text-gray-500 flex flex-col gap-1">
                    <div className="flex justify-between">
                      <span>Từ:</span>
                      <span className="font-medium text-gray-700">{start.format('DD/MM/YYYY HH:mm')}</span>
                    </div>
                    <div className="flex justify-between">
                      <span>Đến:</span>
                      <span className="font-medium text-gray-700">{end.format('DD/MM/YYYY HH:mm')}</span>
                    </div>
                  </div>

                  <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                    <Button 
                      type="text" 
                      size="small"
                      className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                      icon={<Pencil className="w-3.5 h-3.5" />} 
                      onClick={() => {
                        setEditingPromotion(promo);
                        setDrawerVisible(true);
                      }}
                    >
                      Sửa
                    </Button>
                    <Popconfirm
                      title="Xóa voucher này?"
                      onConfirm={() => promo.id && deleteMutation.mutate(promo.id)}
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
              );
            })
          )}
        </div>
      </div>

      <Drawer
        title={editingPromotion ? "Cập nhật Voucher" : "Tạo Voucher mới"}
        width={450}
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        destroyOnClose
      >
        <PromotionForm 
          initialData={editingPromotion} 
          onSuccess={() => setDrawerVisible(false)} 
          onCancel={() => setDrawerVisible(false)} 
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
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái Khuyến mãi</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'ACTIVE', label: 'Đang kích hoạt' },
                { value: 'INACTIVE', label: 'Tạm dừng / Đã kết thúc' }
              ]} 
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
