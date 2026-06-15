import React, { useState } from 'react';
import { Table, Tag, Modal, Spin } from 'antd';
import { useBotSessions, useBotMessages } from '../hooks/useSupervisor';
import dayjs from 'dayjs';
import { Bot, MessageSquare, ShieldAlert } from 'lucide-react';

export const SupervisorTable = () => {
  const { data: sessions, isLoading } = useBotSessions();
  const [selectedSessionId, setSelectedSessionId] = useState<string | undefined>();
  
  const columns = [
    { title: 'Phiên hội thoại', dataIndex: 'id', key: 'id', render: (id: string) => <strong className="text-indigo-600">{id}</strong> },
    { title: 'Khách hàng', dataIndex: 'customer_id', key: 'customer_id' },
    { title: 'Bắt đầu', dataIndex: 'started_at', key: 'started_at', render: (d: string) => dayjs(d).format('HH:mm DD/MM') },
    { title: 'Hoạt động cuối', dataIndex: 'last_activity', key: 'last_activity', render: (d: string) => dayjs(d).format('HH:mm DD/MM') },
    { title: 'Tin nhắn', dataIndex: 'message_count', key: 'message_count' },
    { 
      title: 'Cảm xúc', 
      dataIndex: 'sentiment', 
      key: 'sentiment', 
      render: (s: string) => (
        <Tag color={s === 'POSITIVE' ? 'success' : s === 'NEGATIVE' ? 'error' : 'default'}>{s}</Tag>
      )
    },
    { 
      title: 'Trạng thái', 
      dataIndex: 'status', 
      key: 'status', 
      render: (s: string) => (
        <Tag color={s === 'ACTIVE' ? 'processing' : s === 'HANDED_OVER' ? 'warning' : 'default'}>{s}</Tag>
      )
    },
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2">
          <ShieldAlert className="text-indigo-600 w-5 h-5"/>
          Phiên Chatbot Đang hoạt động
        </h2>
      </div>

      <Table 
        columns={columns} 
        dataSource={sessions} 
        rowKey="id" 
        loading={isLoading}
        onRow={(record) => ({
          onClick: () => setSelectedSessionId(record.id),
          className: 'cursor-pointer hover:bg-slate-50'
        })}
      />

      <Modal
        title={<span className="flex items-center gap-2"><MessageSquare size={18} /> Chi tiết hội thoại {selectedSessionId}</span>}
        open={!!selectedSessionId}
        onCancel={() => setSelectedSessionId(undefined)}
        footer={null}
        width={600}
        destroyOnClose
      >
        <MessageHistory sessionId={selectedSessionId} />
      </Modal>
    </div>
  );
};

const MessageHistory = ({ sessionId }: { sessionId?: string }) => {
  const { data: messages, isLoading } = useBotMessages(sessionId);

  if (isLoading) return <div className="p-8 flex justify-center"><Spin /></div>;

  return (
    <div className="flex flex-col gap-4 mt-4 max-h-[500px] overflow-y-auto p-2">
      {messages?.map(msg => {
        const isBot = msg.sender === 'BOT';
        return (
          <div key={msg.id} className={`flex flex-col ${isBot ? 'items-start' : 'items-end'}`}>
            <div className="flex items-end gap-2 max-w-[80%]">
              {isBot && (
                <div className="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center text-indigo-700 flex-shrink-0">
                  <Bot size={16} />
                </div>
              )}
              <div className={`px-4 py-2 rounded-2xl ${
                isBot 
                  ? 'bg-slate-100 text-slate-800 rounded-bl-sm' 
                  : 'bg-indigo-600 text-white rounded-br-sm'
              }`}>
                <div className="text-[15px]">{msg.content}</div>
                {isBot && msg.confidence_score && (
                  <div className="text-[10px] text-slate-400 mt-1">Độ tự tin: {(msg.confidence_score * 100).toFixed(1)}%</div>
                )}
              </div>
            </div>
            <span className="text-[10px] text-slate-400 mt-1 mx-10">
              {dayjs(msg.created_at).format('HH:mm:ss')}
            </span>
          </div>
        );
      })}
    </div>
  );
};
