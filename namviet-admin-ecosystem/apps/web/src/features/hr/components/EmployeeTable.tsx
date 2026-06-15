import React, { useState } from 'react';
import { Table, Button, Tag, Space, Drawer, Input, Select, Divider, Popconfirm } from 'antd';
import { Search, User, Filter, Plus, Pencil, Trash2 } from 'lucide-react';
import { useEmployees, useDeleteEmployee } from '../hooks';
import { Employee } from '../api';
import dayjs from 'dayjs';

export const EmployeeTable = () => {
  const { data: employees = [], isLoading } = useEmployees();
  const deleteMutation = useDeleteEmployee();
  
  const [drawerVisible, setDrawerVisible] = useState(false);
  const [isFilterDrawerOpen, setIsFilterDrawerOpen] = useState(false);
  const [editingEmployee, setEditingEmployee] = useState<Employee | undefined>();
  const [searchText, setSearchText] = useState('');

  const filteredEmployees = employees.filter((emp: any) => 
    emp.full_name?.toLowerCase().includes(searchText.toLowerCase()) ||
    emp.employee_code?.toLowerCase().includes(searchText.toLowerCase())
  );

  const columns = [
    {
      title: 'Mã NV',
      dataIndex: 'employee_code',
      key: 'employee_code',
      render: (code: string) => <span className="font-mono text-blue-600 font-medium">{code}</span>
    },
    {
      title: 'Họ và tên',
      dataIndex: 'full_name',
      key: 'full_name',
      render: (name: string, record: any) => (
        <div>
          <div className="font-semibold text-gray-800">{name}</div>
          <div className="text-xs text-gray-500">{record.email}</div>
        </div>
      )
    },
    {
      title: 'Phòng ban / Chức vụ',
      key: 'department',
      render: (_: any, record: any) => (
        <div>
          <div className="font-medium text-gray-700">{record.department || 'Chưa cập nhật'}</div>
          <div className="text-xs text-gray-500">{record.position || 'Nhân viên'}</div>
        </div>
      )
    },
    {
      title: 'Ngày vào làm',
      dataIndex: 'join_date',
      key: 'join_date',
      render: (date: string) => date ? dayjs(date).format('DD/MM/YYYY') : 'N/A'
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        status === 'ACTIVE' ? <Tag color="success">Đang làm việc</Tag> : <Tag color="default">Đã nghỉ việc</Tag>
      )
    },
    {
      title: 'Hành động',
      key: 'action',
      width: 120,
      render: (_: any, record: any) => (
        <Space size="middle">
          <Button 
            type="text" 
            className="text-blue-600 hover:text-blue-700 hover:bg-blue-50"
            icon={<Pencil className="w-4 h-4" />} 
            onClick={() => {
              setEditingEmployee(record as Employee);
              setDrawerVisible(true);
            }}
          />
          <Popconfirm
            title="Bạn có chắc chắn muốn xóa nhân viên này?"
            onConfirm={() => deleteMutation.mutate(record.id)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button 
              type="text" 
              danger
              className="hover:bg-red-50"
              icon={<Trash2 className="w-4 h-4" />} 
            />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white rounded-lg shadow h-full flex flex-col">
      <div className="p-6 border-b border-gray-100 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h2 className="text-lg font-semibold text-gray-800">Hồ sơ Nhân viên</h2>
          <p className="text-sm text-gray-500">Quản lý danh sách nhân sự, phòng ban</p>
        </div>
        <div className="flex gap-2 w-full sm:w-auto">
          <Input 
            placeholder="Tìm theo tên, mã NV..." 
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
              setEditingEmployee(undefined);
              setDrawerVisible(true);
            }}
            className="bg-blue-600"
          >
            Thêm Nhân viên
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto bg-gray-50/30">
        <div className="hidden md:block p-6">
          <Table 
            columns={columns} 
            dataSource={filteredEmployees} 
            rowKey="id" 
            loading={isLoading}
            pagination={{ pageSize: 15 }}
          />
        </div>

        <div className="block md:hidden p-4 space-y-4">
          {isLoading ? (
            <div className="text-center py-8 text-gray-500">Đang tải...</div>
          ) : filteredEmployees?.length === 0 ? (
            <div className="text-center py-8 text-gray-500">Không có dữ liệu.</div>
          ) : (
            filteredEmployees?.map((emp: any) => (
              <div key={emp.id} className="bg-white p-4 rounded-xl shadow-sm border border-gray-100 flex flex-col gap-3">
                <div className="flex gap-3 items-start">
                  <div className="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0 bg-blue-50 text-blue-500">
                    <User className="w-5 h-5" />
                  </div>
                  <div className="flex-1 min-w-0">
                    <h3 className="font-semibold text-gray-800 text-sm truncate">{emp.full_name}</h3>
                    <div className="text-xs text-gray-500 mt-0.5">{emp.employee_code} • {emp.position || 'Nhân viên'}</div>
                  </div>
                  {emp.status === 'ACTIVE' ? (
                    <Tag color="success" className="m-0">Đang làm</Tag>
                  ) : (
                    <Tag color="default" className="m-0">Đã nghỉ</Tag>
                  )}
                </div>
                
                <div className="bg-gray-50 p-2 rounded text-xs text-gray-500 flex flex-col gap-1">
                  <div className="flex justify-between">
                    <span>Phòng ban:</span>
                    <span className="font-medium text-gray-700">{emp.department || 'Chưa cập nhật'}</span>
                  </div>
                  <div className="flex justify-between">
                    <span>Ngày vào:</span>
                    <span className="font-medium text-gray-700">{emp.join_date ? dayjs(emp.join_date).format('DD/MM/YYYY') : 'Chưa cập nhật'}</span>
                  </div>
                </div>

                <div className="border-t border-gray-50 pt-2 flex justify-end gap-2 mt-1">
                  <Button 
                    type="text" 
                    size="small" 
                    className="bg-blue-50 text-blue-600 hover:bg-blue-100 px-3"
                    icon={<Pencil className="w-3.5 h-3.5" />} 
                    onClick={() => {
                      setEditingEmployee(emp as Employee);
                      setDrawerVisible(true);
                    }}
                  >
                    Sửa
                  </Button>
                </div>
              </div>
            ))
          )}
        </div>
      </div>

      <Drawer
        title={editingEmployee ? 'Sửa thông tin Nhân viên' : 'Thêm Nhân viên mới'}
        placement="right"
        onClose={() => setDrawerVisible(false)}
        open={drawerVisible}
        width={typeof window !== 'undefined' && window.innerWidth < 768 ? '90vw' : 600}
        destroyOnClose
      >
        <div className="p-4 text-center text-gray-500">
          Form cập nhật hồ sơ nhân viên sẽ được đặt tại đây.
        </div>
      </Drawer>
    </div>
  );
};
