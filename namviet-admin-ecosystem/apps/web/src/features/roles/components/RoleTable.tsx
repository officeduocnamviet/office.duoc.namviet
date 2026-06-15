import React, { useState } from 'react';
import { Table, Button, Space, Modal, Popconfirm, Tag } from 'antd';
import { Plus, Pencil, Trash2, ShieldCheck } from 'lucide-react';
import { RolesRole } from '@namviet/shared-types/src/backend.d';
import { useRoles, useDeleteRole } from '../hooks';
import { toast } from 'sonner';
import { RoleForm } from './RoleForm';

export const RoleTable = () => {
  const { data: roles, isLoading } = useRoles();
  const deleteMutation = useDeleteRole();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingRole, setEditingRole] = useState<RolesRole | null>(null);

  const handleDelete = (id: string, name?: string) => {
    if (name?.toLowerCase() === 'admin') {
      toast.error('Không thể xóa vai trò Admin của hệ thống!');
      return;
    }
    deleteMutation.mutate(id, {
      onSuccess: () => toast.success('Xóa vai trò thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const columns = [
    {
      title: 'Tên Vai trò',
      dataIndex: 'name',
      key: 'name',
      render: (text: string) => (
        <div className="flex items-center gap-2 font-medium text-gray-900">
          <ShieldCheck className={`w-4 h-4 ${text.toLowerCase() === 'admin' ? 'text-red-500' : 'text-blue-500'}`} />
          {text}
          {text.toLowerCase() === 'admin' && (
            <Tag color="red" className="ml-2 border-0">Hệ thống</Tag>
          )}
        </div>
      ),
    },
    {
      title: 'Mô tả',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'Số lượng Quyền',
      dataIndex: 'permissions',
      key: 'permissions',
      render: (permissions: string[]) => (
        <Tag color="blue">{permissions?.length || 0} quyền</Tag>
      ),
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 120,
      render: (_: any, record: RolesRole) => {
        const isAdmin = record.name?.toLowerCase() === 'admin';
        return (
          <Space size="middle">
            <Button 
              type="text" 
              icon={<Pencil className="w-4 h-4" />} 
              onClick={() => {
                setEditingRole(record);
                setIsModalOpen(true);
              }}
            />
            <Popconfirm
              title="Xóa vai trò?"
              description="Bạn có chắc chắn muốn xóa vai trò này không?"
              onConfirm={() => record.id && handleDelete(record.id, record.name)}
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

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h2 className="text-lg font-semibold text-gray-900">Phân quyền Hệ thống</h2>
          <p className="text-sm text-gray-500">Quản lý các vai trò và quyền hạn của người dùng</p>
        </div>
        <Button 
          type="primary" 
          icon={<Plus className="w-4 h-4" />}
          className="bg-blue-600"
          onClick={() => {
            setEditingRole(null);
            setIsModalOpen(true);
          }}
        >
          Thêm Vai trò
        </Button>
      </div>

      <Table 
        columns={columns} 
        dataSource={roles} 
        rowKey="id"
        loading={isLoading}
        pagination={{ pageSize: 10 }}
      />

      <Modal
        title={editingRole ? "Cập nhật Vai trò" : "Thêm Vai trò mới"}
        open={isModalOpen}
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
        width={700}
      >
        <RoleForm 
          initialData={editingRole} 
          onSuccess={() => setIsModalOpen(false)}
          onCancel={() => setIsModalOpen(false)}
        />
      </Modal>
    </div>
  );
};
