import React, { useState } from 'react';
import { Table, Button, Space, Modal, Popconfirm, Tag } from 'antd';
import { Plus, Pencil, Trash2, Users as UsersIcon, MapPin, Briefcase } from 'lucide-react';
import { UsersUser } from '@namviet/shared-types/src/backend.d';
import { useUsers, useDeleteUser } from '../hooks';
import { toast } from 'sonner';
import { UserForm } from './UserForm';
import { useRoles } from '@/features/roles/hooks';
import { useWarehouses } from '@/features/warehouses/hooks';

export const UserTable = () => {
  const { data: users, isLoading: loadingUsers } = useUsers();
  const { data: roles, isLoading: loadingRoles } = useRoles();
  const { data: warehouses, isLoading: loadingWarehouses } = useWarehouses();
  
  const deleteMutation = useDeleteUser();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingUser, setEditingUser] = useState<UsersUser | null>(null);

  const handleDelete = (id: string, email?: string) => {
    if (email?.includes('admin')) {
      toast.error('Không thể xóa tài khoản Admin gốc!');
      return;
    }
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa nhân sự thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const getRoleName = (roleId?: string) => {
    if (!roleId) return 'Chưa phân quyền';
    const role = roles?.find(r => r.id === roleId);
    return role ? role.name : 'Không xác định';
  };

  const getWarehouseName = (warehouseId?: number) => {
    if (!warehouseId) return '---';
    const warehouse = warehouses?.find(w => w.id === warehouseId);
    return warehouse ? warehouse.name : 'Không xác định';
  };

  const columns = [
    {
      title: 'Nhân viên',
      key: 'user',
      render: (_: any, record: UsersUser) => (
        <div className="flex items-center gap-3">
          <div className="w-8 h-8 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center font-bold">
            {record.full_name ? record.full_name.charAt(0).toUpperCase() : 'U'}
          </div>
          <div>
            <div className="font-medium text-gray-900">{record.full_name || 'Chưa cập nhật'}</div>
            <div className="text-xs text-gray-500">{record.email} • {record.phone || 'Chưa có SĐT'}</div>
          </div>
        </div>
      )
    },
    {
      title: 'Vai trò',
      dataIndex: 'role_id',
      key: 'role',
      render: (roleId: string) => (
        <Tag color="indigo">{getRoleName(roleId)}</Tag>
      ),
    },
    {
      title: 'Công tác',
      key: 'work',
      render: (_: any, record: UsersUser) => (
        <div className="space-y-1">
          <div className="flex items-center gap-1.5 text-xs text-gray-600">
            <Briefcase className="w-3.5 h-3.5" />
            Công ty Nam Việt
          </div>
          <div className="flex items-center gap-1.5 text-xs text-gray-600">
            <MapPin className="w-3.5 h-3.5" />
            {getWarehouseName(record.warehouse_id)}
          </div>
        </div>
      ),
    },
    {
      title: 'Trạng thái',
      dataIndex: 'status',
      key: 'status',
      render: (status: string) => (
        <Tag color={status === 'active' || status === 'working' ? 'success' : 'warning'}>
          {status === 'active' || status === 'working' ? 'Đang làm việc' : 'Đã nghỉ / Khóa'}
        </Tag>
      ),
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: UsersUser) => {
        const isAdmin = record.email?.includes('admin');
        return (
          <Space size="middle">
            <Button 
              type="text" 
              icon={<Pencil className="w-4 h-4" />} 
              onClick={() => {
                setEditingUser(record);
                setIsModalOpen(true);
              }}
            />
            <Popconfirm
              title="Xóa nhân sự?"
              description="Bạn có chắc chắn muốn xóa nhân sự này khỏi hệ thống?"
              onConfirm={() => record.id && handleDelete(record.id, record.email)}
              okText="Xóa"
              cancelText="Hủy"
              okButtonProps={{ danger: true }}
              disabled={isAdmin}
            >
              <Button type="text" danger icon={<Trash2 className="w-4 h-4" />} disabled={isAdmin} />
            </Popconfirm>
          </Space>
        );
      },
    },
  ];

  const isLoading = loadingUsers || loadingRoles || loadingWarehouses;

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Quản lý Nhân sự</h2>
          <p className="text-sm text-gray-500">Quản lý tài khoản và phân công công việc nhân viên</p>
        </div>
        <Button 
          type="primary" 
          icon={<Plus className="w-4 h-4" />}
          className="bg-blue-600"
          onClick={() => {
            setEditingUser(null);
            setIsModalOpen(true);
          }}
        >
          Thêm Nhân sự
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={users} 
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editingUser ? "Cập nhật Hồ sơ Nhân sự" : "Thêm Nhân sự mới"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
        width={700}
      >
        <UserForm 
          initialData={editingUser} 
          onSuccess={() => setIsModalOpen(false)}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>
    </div>
  );
};
