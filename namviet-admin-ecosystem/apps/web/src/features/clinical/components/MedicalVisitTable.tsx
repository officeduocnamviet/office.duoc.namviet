import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input, Select, Divider } from 'antd';
import { Search, Stethoscope, Filter, Plus, FileText } from 'lucide-react';
import { useVisits } from '../hooks';
import dayjs from 'dayjs';

export const MedicalVisitTable = () => {
  const { data: visits = [], isLoading } = useVisits();
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [searchText, setSearchText] = useState('');

  const filteredVisits = visits.filter((v: any) => 
    v.patient_name?.toLowerCase().includes(searchText.toLowerCase()) ||
    v.visit_code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã Phiếu khám',
      dataIndex: 'visit_code',
      key: 'visit_code',
      render: (code: string) => <span className="font-mono text-gray-600 font-medium">{code}</span>
    },
    {
      title: 'Bệnh nhân',
      dataIndex: 'patient_name',
      key: 'patient_name',
      render: (name: string, record: any) => (
        <div>
          <div className="font-semibold text-gray-800">{name}</div>
          <div className="text-xs text-gray-500 font-mono">ID: {record.patient_id}</div>
        </div>
      )
    },
    {
      title: 'Bác sĩ khám',
      dataIndex: 'doctor_name',
      key: 'doctor_name',
      render: (doctor: string) => <div className="font-medium text-gray-800">{doctor || 'Chưa cập nhật'}</div>
    },
    {
      title: 'Thời gian',
      key: 'time',
      render: (_: any, record: any) => (
        <span className="text-blue-700">
          {dayjs(record.visit_date).format('DD/MM/YYYY')} {record.start_time?.slice(0, 5)}
        </span>
      )
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'IN_PROGRESS' ? <Tag color="processing">Đang khám</Tag> : 
        status === 'COMPLETED' ? <Tag color="success">Hoàn thành</Tag> : 
        status === 'CANCELLED' ? <Tag color="error">Đã hủy</Tag> : 
        <Tag color="default">{status}</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 140,
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button type="primary" size="small" className="bg-blue-600 text-xs flex items-center gap-1" icon={<FileText size={14} />}>
            Xem Hồ sơ
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Hồ sơ Khám bệnh (Medical Visits)</h2>
          <p className="text-sm text-gray-500">Quản lý bệnh án, chỉ định lâm sàng và kết quả khám</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm theo Mã PK, Tên BN..." 
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
            Tạo Hồ sơ
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredVisits} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 15 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredVisits?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có hồ sơ nào.</div>
          ) : (
            filteredVisits?.map((visit: any) => (
              <div key={visit.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-10 h-10 rounded-lg flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-500">
                    <Stethoscope className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{visit.patient_name}</h3>
                    <div className="text-xs text-gray-500 font-mono mt-0.5">{visit.visit_code}</div>
                  </div>
                  {visit.status === 'IN_PROGRESS' ? <Tag color="processing" className="m-0">Đang khám</Tag> : 
                   visit.status === 'COMPLETED' ? <Tag color="success" className="m-0">Đã xong</Tag> : 
                   <Tag color="default" className="m-0">{visit.status}</Tag>}
                </div>
                
                <div className="bg-gray-50 p-2 rounded text-sm grid grid-cols-2 gap-2 mt-1">
                  <div>
                    <div className="text-xs text-gray-500">Bác sĩ khám:</div>
                    <div className="font-medium text-gray-800 truncate">{visit.doctor_name || 'Chưa rõ'}</div>
                  </div>
                  <div className="text-right">
                    <div className="text-xs text-gray-500">Ngày khám:</div>
                    <div className="font-medium text-gray-800">{dayjs(visit.visit_date).format('DD/MM/YYYY')}</div>
                  </div>
                </div>

                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2 mt-1">
                  <Button type="primary" size="small" className="bg-blue-600 flex items-center gap-1 w-full justify-center">
                    <FileText className="w-3.5 h-3.5" /> Mở Hồ Sơ
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
      
      <Drawer
        title="Lọc Hồ sơ khám"
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
                { value: 'IN_PROGRESS', label: 'Đang khám' },
                { value: 'COMPLETED', label: 'Hoàn thành' },
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
