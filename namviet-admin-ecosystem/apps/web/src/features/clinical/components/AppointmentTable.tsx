import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input, Select, Divider } from 'antd';
import { Search, Calendar, Plus, Filter, CheckCircle, XCircle } from 'lucide-react';
import { useAppointments } from '../hooks';
import dayjs from 'dayjs';

export const AppointmentTable = () => {
  const { data: appointments = [], isLoading } = useAppointments();
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [searchText, setSearchText] = useState('');

  const filteredAppointments = appointments.filter((a: any) => 
    a.patient_name?.toLowerCase().includes(searchText.toLowerCase()) ||
    a.phone_number?.includes(searchText)
  );

  const columns = [
    {
      title: 'Khách hàng / Bệnh nhân',
      dataIndex: 'patient_name',
      key: 'patient_name',
      render: (name: string, record: any) => (
        <div>
          <div className="font-semibold text-gray-800">{name}</div>
          <div className="text-xs text-gray-500 font-mono">{record.phone_number}</div>
        </div>
      )
    },
    {
      title: 'Dịch vụ',
      dataIndex: 'service_type',
      key: 'service_type',
      render: (service: string) => <Tag color="blue">{service || 'Khám Tổng quát'}</Tag>
    },
    {
      title: 'Thời gian hẹn',
      key: 'time',
      render: (_: any, record: any) => (
        <span className="font-medium text-blue-700">
          {dayjs(record.appointment_date).format('DD/MM/YYYY')} {record.appointment_time?.slice(0, 5)}
        </span>
      )
    },
    {
      title: 'Bác sĩ / Nguồn',
      dataIndex: 'doctor_name',
      key: 'doctor_name',
      render: (doctor: string) => doctor || 'Chưa xếp'
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'SCHEDULED' ? <Tag color="blue">Đã lên lịch</Tag> : 
        status === 'CONFIRMED' ? <Tag color="cyan">Đã xác nhận</Tag> : 
        status === 'CHECKED_IN' ? <Tag color="success">Đã đến (Check-in)</Tag> : 
        status === 'CANCELLED' ? <Tag color="error">Đã hủy</Tag> : 
        status === 'NO_SHOW' ? <Tag color="default">Không đến</Tag> : 
        <Tag color="default">{status}</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 140,
      render: (_: any, record: any) => (
        <Space size="middle">
          {record.status === 'SCHEDULED' || record.status === 'CONFIRMED' ? (
            <Button type="text" className="text-green-600 hover:text-green-700 hover:bg-green-50 p-1" title="Check-in">
              <CheckCircle className="w-5 h-5" />
            </Button>
          ) : null}
          <Button type="text" className="text-red-500 hover:text-red-600 hover:bg-red-50 p-1" title="Hủy lịch">
            <XCircle className="w-5 h-5" />
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Lịch Hẹn Khám</h2>
          <p className="text-sm text-gray-500">Quản lý đặt lịch, check-in bệnh nhân</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm tên, SĐT..." 
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
            className="bg-blue-600 flex items-center"
          >
            Thêm Lịch hẹn
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredAppointments} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 10 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredAppointments?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có lịch hẹn nào.</div>
          ) : (
            filteredAppointments?.map((apt: any) => (
              <div key={apt.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-500">
                    <Calendar className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{apt.patient_name}</h3>
                    <div className="text-xs text-gray-500 font-mono mt-0.5">{apt.phone_number}</div>
                  </div>
                  {apt.status === 'SCHEDULED' ? <Tag color="blue" className="m-0">Lên lịch</Tag> : 
                   apt.status === 'CONFIRMED' ? <Tag color="cyan" className="m-0">Đã XN</Tag> : 
                   apt.status === 'CHECKED_IN' ? <Tag color="success" className="m-0">Đã đến</Tag> : 
                   <Tag color="default" className="m-0">{apt.status}</Tag>}
                </div>
                
                <div className="bg-gray-50 p-2 rounded text-sm flex justify-between items-center mt-1">
                  <span className="text-gray-600">Thời gian:</span>
                  <span className="font-bold text-blue-600">{dayjs(apt.appointment_date).format('DD/MM')} {apt.appointment_time?.slice(0, 5)}</span>
                </div>

                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2 mt-1">
                  <Button type="text" size="small" className="text-green-600 bg-green-50 px-3">Check-in</Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      
      <Drawer
        title="Lọc Lịch hẹn"
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
                { value: 'SCHEDULED', label: 'Đã lên lịch' },
                { value: 'CHECKED_IN', label: 'Đã đến (Check-in)' },
                { value: 'CANCELLED', label: 'Đã hủy' }
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
