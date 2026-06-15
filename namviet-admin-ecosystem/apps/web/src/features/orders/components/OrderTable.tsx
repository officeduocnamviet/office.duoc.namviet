import React, { useState } from 'react';
import { Table, Button, Tag, Space, Input, Typography, Drawer, Select, Divider } from 'antd';
import { Search, Eye, FileText, CheckCircle2, Clock, Filter } from 'lucide-react';
import { useOrders } from '../hooks';
import { OrdersOrder } from '@namviet/shared-types/src/backend.d';
import dayjs from 'dayjs';

const { Text } = Typography;

interface OrderTableProps {
  orderType?: 'RETAIL' | 'B2B'; // Nếu truyền vào sẽ tự filter
}

export const OrderTable: React.FC<OrderTableProps> = ({ orderType }) => {
  const { data: orders = [], isLoading } = useOrders(orderType);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [searchText, setSearchText] = useState('');

  const filteredOrders = orders.filter(o => 
    o.code?.toLowerCase().includes(searchText.toLowerCase()) ||
    o.customer_id?.toString().includes(searchText)
  );

  const columns = [
    {
      title: 'Mã Đơn Hàng',
      dataIndex: 'code',
      key: 'code',
      width: 150,
      render: (code: string) => (
        <div className="flex items-center gap-2">
          <FileText size={16} className="text-blue-500" />
          <Text strong className="text-blue-600">{code}</Text>
        </div>
      )
    },
    {
      title: 'Ngày tạo',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date: string) => dayjs(date).format('DD/MM/YYYY HH:mm')
    },
    {
      title: 'Khách hàng (ID)',
      dataIndex: 'customer_id',
      key: 'customer_id',
      render: (id: number) => <span className="font-medium">KH-{id || 'Khách lẻ'}</span>
    },
    {
      title: 'Loại Đơn',
      dataIndex: 'order_type',
      key: 'order_type',
      render: (type: string) => (
        type === 'B2B' 
          ? <Tag color="blue">Bán sỉ (B2B)</Tag> 
          : <Tag color="green">Bán lẻ</Tag>
      )
    },
    {
      title: 'Tổng tiền',
      dataIndex: 'final_amount',
      key: 'final_amount',
      render: (amount: number) => (
        <span className="font-bold text-gray-800">
          {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(amount || 0)}
        </span>
      )
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => {
        if (status === 'COMPLETED') return <Tag color="success" icon={<CheckCircle2 size={12} className="mr-1 inline" />}>Hoàn thành</Tag>;
        if (status === 'PENDING') return <Tag color="warning" icon={<Clock size={12} className="mr-1 inline" />}>Chờ xử lý</Tag>;
        if (status === 'CANCELLED') return <Tag color="error">Đã hủy</Tag>;
        return <Tag>{status}</Tag>;
      }
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 100,
      render: (_: any, record: OrdersOrder) => (
        <Space size="middle">
          <Button type="text" icon={<Eye size={16} className="text-gray-500" />} />
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">
            {orderType === 'B2B' ? 'Đơn hàng Bán Sỉ (B2B)' : orderType === 'RETAIL' ? 'Đơn hàng Bán Lẻ' : 'Tất cả Đơn hàng'}
          </h2>
          <p className="text-sm text-gray-500">Quản lý và theo dõi trạng thái đơn hàng</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm theo mã đơn..." 
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
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredOrders} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredOrders?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có đơn hàng nào.</div>
          ) : (
            filteredOrders?.map((order) => (
              <div key={order.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex justify-between items-start">
                  <div className="flex items-center gap-2">
                    <FileText className="w-5 h-5 text-blue-500" />
                    <div>
                      <h3 className="font-semibold text-blue-600 text-sm">{order.code}</h3>
                      <div className="text-xs text-gray-500">{dayjs(order.created_at).format('DD/MM/YYYY HH:mm')}</div>
                    </div>
                  </div>
                  {order.status === 'COMPLETED' ? <Tag color="success" className="m-0">Xong</Tag> :
                   order.status === 'PENDING' ? <Tag color="warning" className="m-0">Chờ</Tag> :
                   <Tag color="error" className="m-0">Hủy</Tag>}
                </div>
                
                <div className="bg-gray-50 p-3 rounded-lg flex flex-col gap-1">
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Khách hàng:</span>
                    <span className="font-medium">KH-{order.customer_id || 'Khách lẻ'}</span>
                  </div>
                  <div className="flex justify-between text-sm">
                    <span className="text-gray-500">Loại đơn:</span>
                    <span>{order.order_type === 'B2B' ? 'Bán sỉ' : 'Bán lẻ'}</span>
                  </div>
                  <div className="flex justify-between text-sm mt-1 pt-1 border-t border-gray-200">
                    <span className="text-gray-500">Tổng tiền:</span>
                    <span className="font-bold text-blue-600">
                      {new Intl.NumberFormat('vi-VN', { style: 'currency', currency: 'VND' }).format(order.final_amount || 0)}
                    </span>
                  </div>
                </div>

                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2">
                  <Button 
                    type="text" 
                    size="small"
                    className="bg-gray-100 text-gray-600 hover:bg-gray-200 px-3"
                    icon={<Eye className="w-3.5 h-3.5" />} 
                  >
                    Chi tiết
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

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
            <label className="block text-sm font-medium text-gray-700 mb-1">Trạng thái Đơn hàng</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'PENDING', label: 'Chờ xử lý' },
                { value: 'COMPLETED', label: 'Hoàn thành' },
                { value: 'CANCELLED', label: 'Đã hủy' }
              ]} 
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Loại Đơn hàng</label>
            <Select 
              className="w-full" 
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'RETAIL', label: 'Bán lẻ' },
                { value: 'B2B', label: 'Bán sỉ (B2B)' }
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
