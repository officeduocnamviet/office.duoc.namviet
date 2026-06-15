import React from 'react';
import { useChannels } from '../hooks/useChat';
import { InternalChannel } from '../api/chatApi';
import { Hash, Lock, Plus } from 'lucide-react';
import { Spin, Button } from 'antd';

interface Props {
  selectedChannelId: number | undefined;
  onSelectChannel: (id: number) => void;
}

export const ChatChannelList: React.FC<Props> = ({ selectedChannelId, onSelectChannel }) => {
  const { data: channels, isLoading } = useChannels();

  if (isLoading) return <div className="p-4 text-center"><Spin /></div>;

  return (
    <div className="w-64 border-r border-slate-200 bg-slate-50 flex flex-col h-full">
      <div className="p-4 border-b border-slate-200 flex justify-between items-center bg-white">
        <h2 className="font-bold text-slate-800">Kênh nhắn tin</h2>
        <Button size="small" type="text" icon={<Plus size={16} />} />
      </div>
      <div className="flex-1 overflow-y-auto p-2">
        {channels?.map(channel => (
          <div 
            key={channel.id}
            onClick={() => channel.id && onSelectChannel(channel.id)}
            className={`flex items-center gap-2 px-3 py-2 rounded-md cursor-pointer transition-colors ${
              selectedChannelId === channel.id 
                ? 'bg-blue-100 text-blue-700 font-medium' 
                : 'text-slate-600 hover:bg-slate-200'
            }`}
          >
            {channel.type === 'PRIVATE' ? <Lock size={16} /> : <Hash size={16} />}
            <span className="truncate flex-1">{channel.name}</span>
          </div>
        ))}
      </div>
    </div>
  );
};
