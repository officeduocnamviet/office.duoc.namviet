import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input, Select, Divider } from 'antd';
import { Search, ListOrdered, Filter, PlayCircle, SkipForward } from 'lucide-react';
import { useQueues } from '../hooks';
import dayjs from 'dayjs';

export const ClinicalQueueTable = () => {
  const { data: queues = [], isLoading } = useQueues();
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [searchText, setSearchText] = useState('');

  const filteredQueues = queues.filter((q: any) => 
    q.patient_name?.toLowerCase().includes(searchText.toLowerCase()) ||
    q.queue_number?.toString().includes(searchText)
  );

  const columns = [
    {
      title: 'Số TT',
      dataIndex: 'queue_number',
      key: 'queue_number',
      render: (num: number) => <span className="font-bold text-blue-600 text-lg">{num}</span>
    },
    {
      title: 'Bệnh nhân',
      dataIndex: 'patient_name',
      key: 'patient_name',
      render: (name: string, record: any) => (
        <div>
          <div className="font-semibold text-gray-800">{name}</div>
          <div className="text-xs text-gray-500">Mã: {record.patient_id}</div>
        </div>
      )
    },
    {
      title: 'Phòng khám / Bác sĩ',
      key: 'room',
      render: (_: any, record: any) => (
        <div>
          <div className="font-medium text-gray-700">{record.room_id || 'Chưa xếp phòng'}</div>
          <div className="text-xs text-gray-500">{record.doctor_name || 'Chờ xếp BS'}</div>
        </div>
      )
    },
    {
      title: 'Loại khám',
      dataIndex: 'visit_type',
      key: 'visit_type',
      render: (type: string) => <Tag color="blue">{type || 'Khám bệnh'}</Tag>
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'WAITING' ? <Tag color="warning">Đang chờ</Tag> : 
        status === 'IN_PROGRESS' ? <Tag color="processing">Đang khám</Tag> : 
        status === 'COMPLETED' ? <Tag color="success">Đã khám xong</Tag> : 
        status === 'SKIPPED' ? <Tag color="default">Bỏ qua</Tag> : 
        <Tag color="default">{status}</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 140,
      render: (_: any, record: any) => (
        <Space size="middle">
          {record.status === 'WAITING' ? (
            <>
              <Button type="primary" size="small" className="bg-blue-600 text-xs px-2" icon={<PlayCircle size={14} />}>
                Gọi vào
              </Button>
              <Button type="default" size="small" className="text-xs px-2" icon={<SkipForward size={14} />}>
                Bỏ qua
              </Button>
            </>
          ) : record.status === 'IN_PROGRESS' ? (
            <Button type="primary" size="small" className="bg-green-600 text-xs px-2">
              Hoàn thành
            </Button>
          ) : null}
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Hàng đợi Lâm sàng</h2>
          <p className="text-sm text-gray-500">Quản lý thứ tự bệnh nhân vào phòng khám</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm bệnh nhân..." 
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
            dataSource={filteredQueues} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 20 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredQueues?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Hàng đợi trống.</div>
          ) : (
            filteredQueues?.map((queue: any) => (
              <div key={queue.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-12 h-12 rounded-lg flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-600 font-bold text-xl border border-blue-100">
                    {queue.queue_number}
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{queue.patient_name}</h3>
                    <div className="text-xs text-gray-500 font-mono mt-0.5">Khám: {queue.visit_type || 'Tổng quát'}</div>
                  </div>
                  {queue.status === 'WAITING' ? <Tag color="warning" className="m-0">Chờ khám</Tag> : 
                   queue.status === 'IN_PROGRESS' ? <Tag color="processing" className="m-0">Đang khám</Tag> : 
                   queue.status === 'COMPLETED' ? <Tag color="success" className="m-0">Đã khám</Tag> : 
                   <Tag color="default" className="m-0">{queue.status}</Tag>}
                </div>
                
                <div className="bg-gray-50 p-2 rounded text-sm flex justify-between items-center mt-1">
                  <span className="text-gray-600">Phòng khám:</span>
                  <span className="font-medium text-gray-800">{queue.room_id || 'Chưa xếp'}</span>
                </div>

                <div className="border-t border-gray-50 pt-3 flex justify-end gap-2 mt-1">
                  {queue.status === 'WAITING' && (
                    <Button type="primary" className="bg-blue-600 w-full flex items-center justify-center">
                      <PlayCircle className="w-4 h-4 mr-2" /> Gọi vào khám
                    </Button>
                  )}
                  {queue.status === 'IN_PROGRESS' && (
                    <Button type="primary" className="bg-green-600 w-full">Hoàn thành khám</Button>
                  )}
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      
      <Drawer
        title="Lọc Hàng đợi"
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
              placeholder="Tất cả" 
              allowClear
              options={[
                { value: 'WAITING', label: 'Đang chờ' },
                { value: 'IN_PROGRESS', label: 'Đang khám' },
                { value: 'COMPLETED', label: 'Đã khám xong' }
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
