"use client";

import React from 'react';
import { Dropdown, Avatar, Typography, Badge } from 'antd';
import { UserOutlined, SettingOutlined, LogoutOutlined } from '@ant-design/icons';
import { Search, MessageCircle, Bell } from 'lucide-react';
import { useAuthStore } from '@/stores/useAuthStore';
import { useRouter } from 'next/navigation';

export const Header = () => {
  const { user, logout } = useAuthStore();
  const router = useRouter();

  const handleMenuClick = (e: any) => {
    if (e.key === 'logout') {
      logout();
      router.push('/login');
    }
  };

  const userMenu = {
    items: [
      { key: 'profile', icon: <UserOutlined />, label: 'Hồ sơ cá nhân' },
      { key: 'settings', icon: <SettingOutlined />, label: 'Cài đặt' },
      { type: 'divider' as const },
      { key: 'logout', icon: <LogoutOutlined />, label: 'Đăng xuất', danger: true },
    ],
    onClick: handleMenuClick,
  };

  return (
    <header className="h-14 bg-white border-b border-gray-200 flex items-center justify-between px-4 sticky top-0 z-40 w-full shadow-sm">
      {/* Left side: Logo */}
      <div className="flex items-center gap-3 cursor-pointer" onClick={() => router.push('/')}>
        <img src="/logo.png" alt="Nam Việt Logo" className="h-8 w-auto object-contain" />
        <span className="font-black text-orange-600 hidden md:block text-lg">NAM VIỆT ERP</span>
      </div>

      {/* Center/Right side: Tools & Avatar */}
      <div className="flex items-center gap-4 md:gap-6">
        
        {/* Global Search Hint */}
        <div 
          className="hidden md:flex items-center gap-2 bg-gray-100 hover:bg-gray-200 cursor-pointer px-3 py-1.5 rounded-full transition-colors"
          onClick={() => {
            // Phát sự kiện bàn phím ảo để mở Modal (hoặc truyền prop, nhưng dispatch event cho nhanh)
            document.dispatchEvent(new KeyboardEvent('keydown', { key: 'k', ctrlKey: true, altKey: true }));
          }}
        >
          <Search size={16} className="text-gray-500" />
          <span className="text-sm text-gray-500">Tìm kiếm...</span>
          <span className="text-xs bg-white border border-gray-300 rounded px-1.5 text-gray-400 font-medium tracking-widest shadow-sm">Ctrl Alt K</span>
        </div>

        {/* Mobile Search Icon */}
        <Search size={20} className="text-gray-600 md:hidden cursor-pointer" />

        {/* Chat */}
        <Badge count={2} size="small" offset={[-2, 2]}>
          <MessageCircle size={20} className="text-gray-600 cursor-pointer hover:text-orange-500 transition-colors" />
        </Badge>

        {/* Notifications */}
        <Badge count={5} size="small" offset={[-2, 2]}>
          <Bell size={20} className="text-gray-600 cursor-pointer hover:text-orange-500 transition-colors" />
        </Badge>

        {/* Divider */}
        <div className="w-px h-6 bg-gray-200 hidden md:block"></div>

        {/* User Profile */}
        <Dropdown menu={userMenu} placement="bottomRight" arrow trigger={['click']}>
          <div className="flex items-center gap-2 cursor-pointer hover:bg-gray-50 p-1 rounded-full pr-2 transition-colors border border-transparent hover:border-gray-200">
            <Avatar style={{ backgroundColor: '#f97316' }} icon={<UserOutlined />} size="small" />
            <Typography.Text strong className="hidden md:block text-sm">{user?.fullName || 'Người dùng'}</Typography.Text>
          </div>
        </Dropdown>
      </div>
    </header>
  );
};
