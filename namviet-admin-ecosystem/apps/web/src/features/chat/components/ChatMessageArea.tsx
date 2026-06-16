import React, { useState, useEffect, useRef } from 'react';
import { useMessages, useSendMessage } from '../hooks/useChat';
import { InternalMessage } from '../api/chatApi';
import { supabase } from '@/lib/supabase';
import { useQueryClient } from '@tanstack/react-query';
import { Spin, Input, Button } from 'antd';
import { Send } from 'lucide-react';
import dayjs from 'dayjs';
import { toast } from 'sonner';
import { useAuthStore } from '@/stores/useAuthStore';

interface Props {
  channelId: number;
}

export const ChatMessageArea: React.FC<Props> = ({ channelId }) => {
  const queryClient = useQueryClient();
  const { data: messages = [], isLoading } = useMessages(channelId);
  const sendMutation = useSendMessage();
  const { user } = useAuthStore();
  const currentUserId = user?.id || 'CURRENT_USER';
  
  const [content, setContent] = useState('');
  const messagesEndRef = useRef<HTMLDivElement>(null);

  // Setup Supabase Realtime Listener
  useEffect(() => {
    if (!channelId) return;

    const channel = supabase.channel(`chat_messages_${channelId}`)
      .on('postgres_changes', {
        event: 'INSERT',
        schema: 'public',
        table: 'chat_messages',
        filter: `channel_id=eq.${channelId}`
      }, (payload) => {
        const newMessage = payload.new as InternalMessage;
        
        // Push notification if not from me
        if (newMessage.sender_id !== currentUserId) {
          toast.info(`Tin nhắn mới: ${newMessage.content}`, {
            description: 'Từ cuộc trò chuyện nội bộ',
          });
        }

        // Append new message to cache
        queryClient.setQueryData(['chat_messages', channelId], (oldData: InternalMessage[] | undefined) => {
          if (!oldData) return [newMessage];
          return [...oldData, newMessage];
        });
      })
      .subscribe();

    return () => {
      supabase.removeChannel(channel);
    };
  }, [channelId, queryClient]);

  // Scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  }, [messages]);

  const handleSend = () => {
    if (!content.trim()) return;
    
    sendMutation.mutate({
      channel_id: channelId,
      content,
      sender_id: currentUserId
    }, {
      onSuccess: () => {
        setContent('');
      }
    });
  };

  if (isLoading) return <div className="flex-1 flex items-center justify-center bg-white"><Spin /></div>;

  return (
    <div className="flex-1 flex flex-col bg-white h-full relative">
      <div className="flex-1 overflow-y-auto p-6 space-y-4">
        {messages.map((msg, idx) => {
          const isMe = msg.sender_id === currentUserId;
          return (
            <div key={msg.id || idx} className={`flex flex-col ${isMe ? 'items-end' : 'items-start'}`}>
              <div className="flex items-end gap-2 max-w-[70%]">
                {!isMe && (
                  <div className="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-700 font-bold text-xs flex-shrink-0">
                    {msg.sender_id?.substring(0, 2).toUpperCase()}
                  </div>
                )}
                <div className={`px-4 py-2 rounded-2xl ${
                  isMe 
                    ? 'bg-blue-600 text-white rounded-br-sm' 
                    : 'bg-slate-100 text-slate-800 rounded-bl-sm'
                }`}>
                  <div className="text-[15px]">{msg.content}</div>
                </div>
              </div>
              <span className="text-[10px] text-slate-400 mt-1 mx-10">
                {dayjs(msg.created_at).format('HH:mm')}
              </span>
            </div>
          );
        })}
        <div ref={messagesEndRef} />
      </div>

      <div className="p-4 bg-white border-t border-slate-100">
        <div className="flex items-center gap-2">
          <Input 
            value={content}
            onChange={(e) => setContent(e.target.value)}
            onPressEnter={handleSend}
            placeholder="Nhập tin nhắn..."
            className="rounded-full px-4 py-2"
          />
          <Button 
            type="primary" 
            shape="circle" 
            icon={<Send size={16} />} 
            onClick={handleSend}
            className="bg-blue-600"
            loading={sendMutation.isPending}
          />
        </div>
      </div>
    </div>
  );
};
