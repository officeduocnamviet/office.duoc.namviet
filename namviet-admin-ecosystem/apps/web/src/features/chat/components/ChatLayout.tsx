"use client";

import React, { useState } from 'react';
import { ChatChannelList } from './ChatChannelList';
import { ChatMessageArea } from './ChatMessageArea';

export const ChatLayout = () => {
  const [selectedChannel, setSelectedChannel] = useState<number | undefined>();

  return (
    <div className="flex h-[calc(100vh-140px)] rounded-2xl overflow-hidden border border-slate-200 shadow-sm bg-white">
      <ChatChannelList 
        selectedChannelId={selectedChannel} 
        onSelectChannel={setSelectedChannel} 
      />
      {selectedChannel ? (
        <ChatMessageArea channelId={selectedChannel} />
      ) : (
        <div className="flex-1 flex flex-col items-center justify-center bg-slate-50 text-slate-400">
          <div className="w-16 h-16 bg-slate-200 rounded-full mb-4 flex items-center justify-center">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-slate-400"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>
          </div>
          <p>Chọn một kênh để bắt đầu trò chuyện</p>
        </div>
      )}
    </div>
  );
};
