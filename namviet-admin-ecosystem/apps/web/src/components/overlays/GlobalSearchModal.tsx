"use client";

import React, { useEffect, useRef, useState } from 'react';
import { Search, X, Package, Users, FileText, Settings, ArrowRight, LayoutDashboard, Boxes, Store, Building2 } from 'lucide-react';

interface GlobalSearchModalProps {
  isOpen: boolean;
  onClose: () => void;
}

import { useRouter } from 'next/navigation';

const SEARCH_ITEMS = [
  { id: 'companies', title: 'Quản lý Công ty', category: 'Hệ thống', icon: Building2, color: 'text-indigo-500', bg: 'bg-indigo-100', href: '/companies' },
  { id: 1, title: 'Quản lý kho Tổng', category: 'Kho hàng', icon: Boxes, color: 'text-orange-500', bg: 'bg-orange-100', href: '/warehouses' },
  { id: 2, title: 'Danh sách khách hàng B2B', category: 'Khách hàng', icon: Users, color: 'text-blue-500', bg: 'bg-blue-100', href: '/customers' },
  { id: 3, title: 'Danh sách nhân sự', category: 'Hệ thống', icon: Users, color: 'text-emerald-500', bg: 'bg-emerald-100', href: '/users' },
  { id: 4, title: 'Danh sách Cửa hàng bán lẻ', category: 'Đơn hàng', icon: Store, color: 'text-purple-500', bg: 'bg-purple-100', href: '/retail/list' },
  { id: 5, title: 'Tạo đơn hàng B2B', category: 'Đơn hàng', icon: Package, color: 'text-amber-500', bg: 'bg-amber-100', href: '/b2b/create' },
  { id: 6, title: 'Cấu hình phân quyền hệ thống', category: 'Hệ thống', icon: Settings, color: 'text-slate-500', bg: 'bg-slate-100', href: '/roles' },
];

export default function GlobalSearchModal({ isOpen, onClose }: GlobalSearchModalProps) {
  const router = useRouter();
  const [query, setQuery] = useState('');
  const inputRef = useRef<HTMLInputElement>(null);

  // Focus input khi mở modal
  useEffect(() => {
    if (isOpen) {
      setTimeout(() => {
        inputRef.current?.focus();
      }, 100);
    } else {
      setQuery('');
    }
  }, [isOpen]);

  // Đóng modal khi bấm phím ESC
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <>
      <div 
        className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[200] animate-in fade-in duration-200"
        onClick={onClose}
      />
      
      <div className="fixed top-[10%] left-1/2 -translate-x-1/2 w-[90%] max-w-2xl bg-white dark:bg-slate-900 rounded-2xl shadow-2xl z-[210] overflow-hidden flex flex-col max-h-[80vh] animate-in slide-in-from-top-10 zoom-in-95 duration-200 border border-slate-200 dark:border-slate-700">
        
        {/* Search Header */}
        <div className="flex items-center px-4 py-3 border-b border-slate-100 dark:border-slate-800 shrink-0 bg-white dark:bg-slate-900">
          <Search size={20} className="text-primary mr-3 shrink-0" />
          <input
            ref={inputRef}
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            placeholder="Tìm kiếm tính năng, đơn hàng, khách hàng..."
            className="flex-1 bg-transparent text-slate-800 dark:text-white text-[15px] outline-none"
          />
          <button 
            onClick={onClose}
            className="p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors ml-2"
          >
            <X size={20} />
          </button>
        </div>

        {/* Results Body */}
        <div className="flex-1 overflow-y-auto p-4 sm:p-5 custom-scrollbar bg-slate-50/50 dark:bg-slate-900">
          <div className="mb-3 flex items-center justify-between">
            <h3 className="text-xs font-bold text-slate-500 uppercase tracking-wider">Kết quả tìm kiếm nổi bật</h3>
          </div>
          
          {/* Card Grid Layout */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
            {SEARCH_ITEMS.filter(item => item.title.toLowerCase().includes(query.toLowerCase()) || item.category.toLowerCase().includes(query.toLowerCase())).map((item) => (
              <button 
                key={item.id}
                onClick={() => {
                  router.push(item.href);
                  onClose();
                }}
                className="flex items-center p-3 bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl hover:border-primary dark:hover:border-primary-500 hover:shadow-md transition-all group text-left"
              >
                <div className={`w-10 h-10 rounded-lg flex items-center justify-center shrink-0 mr-3 ${item.bg} dark:bg-opacity-20 ${item.color}`}>
                  <item.icon size={20} />
                </div>
                <div className="flex-1 min-w-0">
                  <h4 className="text-sm font-bold text-slate-800 dark:text-slate-100 truncate group-hover:text-primary transition-colors">
                    {item.title}
                  </h4>
                  <p className="text-[11px] font-medium text-slate-500 truncate mt-0.5">
                    Phân loại: {item.category}
                  </p>
                </div>
                <div className="text-slate-300 dark:text-slate-600 group-hover:text-primary transition-colors shrink-0">
                  <ArrowRight size={16} />
                </div>
              </button>
            ))}
          </div>
        </div>
        
        {/* Footer info */}
        <div className="px-5 py-3 border-t border-slate-100 dark:border-slate-800 bg-slate-50 dark:bg-slate-900/50 shrink-0 flex items-center gap-4 text-[11px] text-slate-500 font-medium">
          <span className="flex items-center gap-1">
            <kbd className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded px-1.5 py-0.5 shadow-sm">↑↓</kbd> để di chuyển
          </span>
          <span className="flex items-center gap-1">
            <kbd className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded px-1.5 py-0.5 shadow-sm">Enter</kbd> để chọn
          </span>
          <span className="flex items-center gap-1">
            <kbd className="bg-white dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded px-1.5 py-0.5 shadow-sm">ESC</kbd> để đóng
          </span>
        </div>

      </div>
    </>
  );
}
