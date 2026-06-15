import { useQuery } from '@tanstack/react-query';
import { supervisorApi } from '../api/supervisorApi';

export const useBotSessions = () => {
  return useQuery({
    queryKey: ['bot_sessions'],
    queryFn: supervisorApi.getSessions,
  });
};

export const useBotMessages = (sessionId: string | undefined) => {
  return useQuery({
    queryKey: ['bot_messages', sessionId],
    queryFn: () => supervisorApi.getMessages(sessionId!),
    enabled: !!sessionId,
  });
};
