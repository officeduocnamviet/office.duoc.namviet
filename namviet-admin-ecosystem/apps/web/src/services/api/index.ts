import { createApiClient } from '@namviet/shared-utils';
import { createClient } from '@/lib/supabase/client';

const supabase = createClient();

const getBaseUrl = () => {
  if (process.env.NEXT_PUBLIC_API_URL) return process.env.NEXT_PUBLIC_API_URL;
  if (process.env.NODE_ENV === 'production') return 'https://namviet-official-150879831872.asia-southeast1.run.app/api/v1';
  return 'http://localhost:8080/api/v1';
};

export const apiClient = createApiClient({
  baseURL: getBaseUrl(),
  getToken: async () => {
    const { data: { session }, error } = await supabase.auth.getSession();
    if (error || !session) return null;
    return session.access_token;
  },
  onUnauthorized: () => {
    // TODO: Bắn event văng ra màn hình đăng nhập hoặc redirect sang /login
    if (typeof window !== 'undefined') {
      // window.location.href = '/login';
    }
  }
});
