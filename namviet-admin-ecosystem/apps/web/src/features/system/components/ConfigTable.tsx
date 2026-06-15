import React, { useState } from 'react';
import { Table, Button, Space, Drawer, Input, Popconfirm } from 'antd';
import { Plus, Pencil, Trash2, Search, Settings } from 'lucide-react';
import { useSystemConfigs, useDeleteSystemConfig } from '../hooks/useConfigs';
import { SystemConfig } from '../api/configApi';
import { ConfigForm } from './ConfigForm';
import { toast } from 'sonner';
import dayjs from 'dayjs';

export const ConfigTable = () => {
  const { data: configs, isLoading } = useSystemConfigs();
  const deleteMutation = useDeleteSystemConfig();
  
  const [searchTerm, setSearchTerm] = useState('');
  const [isDrawerOpen, setIsDrawerOpen] = useState(false);
  const [editingRecord, setEditingRecord] = useState<SystemConfig | null>(null);

  const handleDelete = (key: string) => {
    deleteMutation.mutate(key, {
      onSuccess: () => toast.success('Xóa cấu hình thành công'),
      onError: (err) => toast.error(`Xóa thất bại: ${err.message}`)
    });
  };

  const openDrawer = (record?: SystemConfig) => {
    setEditingRecord(record || null);
    setIsDrawerOpen(true);
  };

  const filteredConfigs = configs?.filter(c => 
    c.config_key?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    c.description?.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const columns = [
    {
      title: 'Khóa cấu hình (Key)',
      dataIndex: 'config_key',
      key: 'config_key',
      render: (text: string) => <strong className="text-blue-600">{text}</strong>
    },
    {
      title: 'Mô tả',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'Giá trị',
      key: 'config_value',
      render: (_: any, record: SystemConfig) => (
        <div className="max-w-xs truncate text-xs text-gray-500 font-mono bg-gray-50 p-1 rounded">
          {typeof record.config_value === 'object' ? JSON.stringify(record.config_value) : String(record.config_value)}
        </div>
      ),
    },
    {
      title: 'Cập nhật lần cuối',
      dataIndex: 'updated_at',
      key: 'updated_at',
      render: (date: string) => date ? dayjs(date).format('DD/MM/YYYY HH:mm') : '---'
    },
    {
      title: 'Thao tác',
      key: 'action',
      width: 100,
      render: (_: any, record: SystemConfig) => (
        <Space size="small">
          <Button 
            type="text" 
            size="small" 
            icon={<Pencil className="w-4 h-4 text-gray-500" />} 
            onClick={() => openDrawer(record)} 
          />
          <Popconfirm
            title="Xóa cấu hình?"
            description="Bạn có chắc chắn muốn xóa khóa cấu hình này không?"
            onConfirm={() => record.config_key && handleDelete(record.config_key)}
            okText="Xóa"
            cancelText="Hủy"
            okButtonProps={{ danger: true }}
          >
            <Button type="text" danger size="small" icon={<Trash2 className="w-4 h-4" />} />
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex flex-col md:flex-row justify-between items-start md:items-center gap-4 mb-6">
        <div>
          <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
            <Settings className="w-5 h-5 text-blue-600" />
            Cấu hình Hệ thống
          </h2>
          <p className="text-sm text-slate-500 mt-1">Quản lý các tham số, API Key và cấu hình động cho toàn bộ hệ thống</p>
        </div>
        <div className="flex items-center gap-3 w-full md:w-auto">
          <Input
            placeholder="Tìm kiếm cấu hình..."
            prefix={<Search className="w-4 h-4 text-slate-400" />}
            value={searchTerm}
            onChange={e => setSearchTerm(e.target.value)}
            className="w-full md:w-64"
          />
          <Button 
            type="primary" 
            icon={<Plus className="w-4 h-4" />} 
            className="bg-blue-600"
            onClick={() => openDrawer()}
          >
            Thêm mới
          </Button>
        </div>
      </div>

      <Table 
        columns={columns} 
        dataSource={filteredConfigs} 
        rowKey="config_key"
        loading={isLoading}
        pagination={{ pageSize: 15 }}
        scroll={{ x: 'max-content' }}
        size="middle"
      />

      <Drawer
        title={editingRecord ? "Cập nhật Cấu hình" : "Thêm mới Cấu hình"}
        width={typeof window !== 'undefined' && window.innerWidth > 768 ? 500 : '90vw'}
        onClose={() => setIsDrawerOpen(false)}
        open={isDrawerOpen}
        destroyOnClose
      >
        <ConfigForm 
          initialData={editingRecord || undefined}
          onSuccess={() => {
            setIsDrawerOpen(false);
            toast.success(editingRecord ? 'Cập nhật thành công' : 'Tạo mới thành công');
          }}
          onCancel={() => setIsDrawerOpen(false)}
        />
      </Drawer>
    </div>
  );
};
