import { apiClient } from '@/lib/axios';

export interface BotSession {
  id: string;
  customer_id: string;
  started_at: string;
  last_activity: string;
  message_count: number;
  status: 'ACTIVE' | 'CLOSED' | 'HANDED_OVER';
  sentiment: 'POSITIVE' | 'NEUTRAL' | 'NEGATIVE';
}

export interface BotMessageLog {
  id: string;
  session_id: string;
  sender: 'BOT' | 'CUSTOMER';
  content: string;
  created_at: string;
  confidence_score?: number;
}

export const supervisorApi = {
  getSessions: async (): Promise<BotSession[]> => {
    // const { data } = await apiClient.get<BotSession[]>('/ai/bot-sessions');
    // return data;
    
    // Mock data for UI presentation
    return [
      {
        id: 'SES-001',
        customer_id: 'CUST-1002',
        started_at: new Date(Date.now() - 3600000).toISOString(),
        last_activity: new Date(Date.now() - 300000).toISOString(),
        message_count: 12,
        status: 'ACTIVE',
        sentiment: 'NEUTRAL'
      },
      {
        id: 'SES-002',
        customer_id: 'CUST-1005',
        started_at: new Date(Date.now() - 7200000).toISOString(),
        last_activity: new Date(Date.now() - 3600000).toISOString(),
        message_count: 8,
        status: 'HANDED_OVER',
        sentiment: 'NEGATIVE'
      }
    ];
  },

  getMessages: async (sessionId: string): Promise<BotMessageLog[]> => {
    // const { data } = await apiClient.get<BotMessageLog[]>(`/ai/bot-sessions/${sessionId}/messages`);
    // return data;
    
    return [
      {
        id: 'MSG-001',
        session_id: sessionId,
        sender: 'CUSTOMER',
        content: 'Tôi muốn hỏi về giá xét nghiệm máu',
        created_at: new Date(Date.now() - 300000).toISOString()
      },
      {
        id: 'MSG-002',
        session_id: sessionId,
        sender: 'BOT',
        content: 'Chào bạn, giá xét nghiệm máu cơ bản tại phòng khám là 150.000đ. Bạn có muốn tôi đặt lịch khám giúp bạn không?',
        created_at: new Date(Date.now() - 290000).toISOString(),
        confidence_score: 0.95
      }
    ];
  }
};
