import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, DatePicker, Select, Divider } from 'antd';
import { Filter, Clock, CheckCircle, XCircle } from 'lucide-react';
import { useAttendanceLogs } from '../hooks';
import dayjs from 'dayjs';

export const TimeAttendanceTable = () => {
  const { data: attendances = [], isLoading } = useAttendanceLogs();
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [selectedDate, setSelectedDate] = useState<dayjs.Dayjs | null>(dayjs());

  const filteredAttendances = attendances.filter((att: any) => 
    !selectedDate || dayjs(att.work_date).format('YYYY-MM-DD') === selectedDate.format('YYYY-MM-DD')
  );

  const columns = [
    {
      title: 'Mã NV',
      dataIndex: 'employee_id',
      key: 'employee_id',
      render: (id: string) => <span className="font-mono text-gray-600 font-medium">{id}</span>
    },
    {
      title: 'Ngày làm việc',
      dataIndex: 'work_date',
      key: 'work_date',
      render: (date: string) => <span className="font-medium text-gray-800">{dayjs(date).format('DD/MM/YYYY')}</span>
    },
    {
      title: 'Giờ vào',
      dataIndex: 'check_in_time',
      key: 'check_in_time',
      render: (time: string) => time ? dayjs(time).format('HH:mm') : '--:--'
    },
    {
      title: 'Giờ ra',
      dataIndex: 'check_out_time',
      key: 'check_out_time',
      render: (time: string) => time ? dayjs(time).format('HH:mm') : '--:--'
    },
    {
      title: 'Tổng giờ',
      dataIndex: 'total_hours',
      key: 'total_hours',
      render: (hours: number) => hours ? <span className="font-semibold">{hours.toFixed(1)}h</span> : '-'
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'PRESENT' ? <Tag color="success">Có mặt</Tag> : 
        status === 'ABSENT' ? <Tag color="error">Vắng mặt</Tag> : 
        status === 'LATE' ? <Tag color="warning">Đi trễ</Tag> : 
        <Tag color="default">{status}</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 120,
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button type="text" className="text-green-600 hover:text-green-700 hover:bg-green-50 p-1" title="Xác nhận đúng giờ">
            <CheckCircle className="w-5 h-5" />
          </Button>
          <Button type="text" className="text-red-500 hover:text-red-600 hover:bg-red-50 p-1" title="Đánh vắng">
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
          <h2 className="text-lg font-semibold text-gray-800">Chấm công (Time Attendance)</h2>
          <p className="text-sm text-gray-500">Theo dõi giờ giấc vào/ra của nhân viên</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <DatePicker 
            className="w-full sm:w-40" 
            value={selectedDate} 
            onChange={(date) => setSelectedDate(date)} 
            format="DD/MM/YYYY"
            allowClear
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
            dataSource={filteredAttendances} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 20 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredAttendances?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có dữ liệu.</div>
          ) : (
            filteredAttendances?.map((att: any) => (
              <div key={att.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-500">
                    <Clock className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm">Mã NV: {att.employee_id}</h3>
                    <div className="text-xs text-gray-500 mt-0.5">{dayjs(att.work_date).format('DD/MM/YYYY')}</div>
                  </div>
                  {att.status === 'PRESENT' ? <Tag color="success" className="m-0">Có mặt</Tag> : 
                   att.status === 'ABSENT' ? <Tag color="error" className="m-0">Vắng</Tag> : 
                   att.status === 'LATE' ? <Tag color="warning" className="m-0">Trễ</Tag> : 
                   <Tag color="default" className="m-0">N/A</Tag>}
                </div>
                
                <div className="bg-gray-50 p-2 rounded text-sm grid grid-cols-2 gap-2 mt-1">
                  <div>
                    <div className="text-xs text-gray-500">Giờ vào:</div>
                    <div className="font-medium text-gray-800">{att.check_in_time ? dayjs(att.check_in_time).format('HH:mm') : '--:--'}</div>
                  </div>
                  <div className="text-right">
                    <div className="text-xs text-gray-500">Giờ ra:</div>
                    <div className="font-medium text-gray-800">{att.check_out_time ? dayjs(att.check_out_time).format('HH:mm') : '--:--'}</div>
                  </div>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      
      <Drawer
        title="Lọc Dữ liệu"
        placement="right"
        onClose={() => setIsFilterDrawerOpen(false)}
        open={isFilterDrawerOpen}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 400}
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
                { value: 'PRESENT', label: 'Có mặt' },
                { value: 'ABSENT', label: 'Vắng mặt' },
                { value: 'LATE', label: 'Đi trễ' }
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
