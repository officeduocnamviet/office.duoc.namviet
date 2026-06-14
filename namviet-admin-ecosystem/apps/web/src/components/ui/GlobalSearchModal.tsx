"use client";

import React, { useState, useEffect } from 'react';
import { Modal, Input, List } from 'antd';
import { SearchOutlined, AppstoreOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';

export const GlobalSearchModal = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchText, setSearchText] = useState('');
  const router = useRouter();

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Ctrl + Alt + K (Windows) or Cmd + K / Cmd + Alt + K (Mac)
      if ((e.ctrlKey || e.metaKey) && e.altKey && e.key === 'k') {
        e.preventDefault();
        setIsModalOpen(true);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  const navigationOptions = [
    { title: 'Quản lý Sản phẩm', link: '/products' },
    { title: 'Tạo Đơn hàng POS', link: '/pos' },
    { title: 'Kho bãi / Chi nhánh', link: '/warehouses' },
    { title: 'Quản lý Người dùng', link: '/users' },
  ].filter(opt => opt.title.toLowerCase().includes(searchText.toLowerCase()));

  return (
    <Modal
      title={<Input 
        prefix={<SearchOutlined className="text-gray-400" />} 
        placeholder="Tìm kiếm chức năng... (VD: Sản phẩm)" 
        bordered={false} 
        autoFocus 
        className="text-lg"
        onChange={(e) => setSearchText(e.target.value)}
      />}
      open={isModalOpen}
      onCancel={() => setIsModalOpen(false)}
      footer={null}
      closable={false}
      styles={{ body: { padding: 0 } }}
      style={{ top: 50 }}
    >
      <List
        className="max-h-[60vh] overflow-y-auto"
        dataSource={navigationOptions}
        renderItem={(item) => (
          <List.Item 
            className="hover:bg-blue-50 cursor-pointer px-6 transition-colors"
            onClick={() => {
              setIsModalOpen(false);
              router.push(item.link);
            }}
          >
            <div className="flex items-center gap-3">
              <AppstoreOutlined className="text-blue-500" />
              <span className="font-medium">{item.title}</span>
            </div>
          </List.Item>
        )}
      />
    </Modal>
  );
};
