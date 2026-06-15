"use client";

import React from 'react';
import { ChatLayout } from '@/features/chat/components/ChatLayout';
import { MessageSquare } from 'lucide-react';

export default function ChatPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto h-full">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <MessageSquare className="w-6 h-6 text-blue-600" />
          Kênh Giao tiếp Nội bộ (Internal Chat)
        </h1>
        <p className="text-slate-500">
          Trao đổi thông tin công việc nội bộ qua các nhóm chat và tin nhắn cá nhân (Tích hợp Realtime).
        </p>
      </div>

      <ChatLayout />
    </div>
  );
}
