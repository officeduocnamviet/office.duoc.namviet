"use client";

import React, { useState, useEffect, useRef, useMemo } from 'react';
import { Modal, Input, List, Typography, Divider, Badge } from 'antd';
import { SearchOutlined, ArrowRightOutlined, FolderOpenOutlined } from '@ant-design/icons';
import { useRouter } from 'next/navigation';
import { MENU_DATA } from '@/config/menu.config';

const { Text } = Typography;

export const GlobalSearchModal = () => {
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [searchText, setSearchText] = useState('');
  const inputRef = useRef<any>(null);
  const router = useRouter();

  // Mở modal khi bấm phím tắt
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Bắt phím Ctrl + Alt + K hoặc Cmd + K
      if ((e.ctrlKey || e.metaKey) && (e.altKey || e.key === 'k') && e.key === 'k') {
        e.preventDefault();
        setIsModalOpen(true);
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  // Tự động focus vào ô input khi modal vừa mở lên
  useEffect(() => {
    if (isModalOpen) {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 100);
    } else {
      setSearchText('');
    }
  }, [isModalOpen]);

  // Làm phẳng toàn bộ MENU_DATA thành 1 mảng dễ search
  const allMenus = useMemo(() => {
    let result: { group: string; title: string; link: string; icon: any }[] = [];
    MENU_DATA.forEach(group => {
      group.items.forEach(item => {
        if (item.children) {
          item.children.forEach(child => {
            result.push({
              group: item.name,
              title: child.name,
              link: child.href,
              icon: item.icon,
            });
          });
        }
      });
    });
    return result;
  }, []);

  // Lọc dữ liệu theo từ khoá
  const filteredMenus = useMemo(() => {
    if (!searchText.trim()) return allMenus;
    const lowerSearch = searchText.toLowerCase();
    return allMenus.filter(m => 
      m.title.toLowerCase().includes(lowerSearch) || 
      m.group.toLowerCase().includes(lowerSearch)
    );
  }, [searchText, allMenus]);

  return (
    <Modal
      title={
        <Input 
          ref={inputRef}
          prefix={<SearchOutlined className="text-blue-500 text-xl mr-2" />} 
          placeholder="Bạn muốn tìm giao diện nào? (VD: Khách hàng, Sản phẩm...)" 
          bordered={false} 
          className="text-lg py-4 border-b border-gray-100"
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
        />
      }
      open={isModalOpen}
      onCancel={() => setIsModalOpen(false)}
      footer={null}
      closable={false}
      styles={{ 
        body: { padding: 0 }, 
        header: { marginBottom: 0, padding: 0 } 
      }}
      width={600}
      style={{ top: 80 }}
    >
      <div className="bg-slate-50 p-3 text-xs text-slate-500 flex justify-between items-center border-b border-slate-100">
        <span>Gợi ý: Gõ từ khoá để lọc nhanh các tính năng</span>
        <div className="flex gap-2">
          <kbd className="bg-white border rounded px-1.5 shadow-sm font-mono">↑↓</kbd> Lên xuống
          <kbd className="bg-white border rounded px-1.5 shadow-sm font-mono ml-2">Enter</kbd> Chọn
        </div>
      </div>

      <List
        className="max-h-[50vh] overflow-y-auto"
        dataSource={filteredMenus}
        renderItem={(item) => {
          const Icon = item.icon || FolderOpenOutlined;
          return (
            <List.Item 
              className="hover:bg-blue-50 cursor-pointer px-6 py-4 transition-colors group border-b-0"
              onClick={() => {
                setIsModalOpen(false);
                router.push(item.link);
              }}
            >
              <div className="flex items-center justify-between w-full">
                <div className="flex items-center gap-4">
                  <div className="w-10 h-10 rounded-full bg-slate-100 group-hover:bg-blue-100 flex items-center justify-center transition-colors">
                    <Icon className="w-5 h-5 text-slate-500 group-hover:text-blue-600" />
                  </div>
                  <div>
                    <div className="font-semibold text-gray-800 group-hover:text-blue-700 text-base">{item.title}</div>
                    <div className="text-xs text-gray-400 flex items-center gap-1 mt-0.5">
                      {item.group} <ArrowRightOutlined className="text-[10px]" /> {item.link}
                    </div>
                  </div>
                </div>
                <Badge status="processing" className="opacity-0 group-hover:opacity-100 transition-opacity" />
              </div>
            </List.Item>
          );
        }}
        locale={{ emptyText: <div className="p-8 text-gray-400">Không tìm thấy chức năng phù hợp</div> }}
      />
    </Modal>
  );
};
