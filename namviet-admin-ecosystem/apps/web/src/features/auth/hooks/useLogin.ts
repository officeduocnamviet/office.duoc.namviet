import { useMutation } from '@tanstack/react-query';
import { authApi } from '../api/authApi';
import { LoginRequest, LoginResponse } from '../types';
import { useAuthStore } from '@/stores/useAuthStore';
import { toast } from 'sonner';

export const useLogin = () => {
  const setAuth = useAuthStore((state) => state.setAuth);

  return useMutation<LoginResponse, Error, LoginRequest>({
    mutationFn: authApi.login,
    onSuccess: (data) => {
      setAuth(data.token, data.user);
      toast.success('Đăng nhập thành công!');
    },
    onError: (error) => {
      toast.error(error.message || 'Sai email hoặc mật khẩu!');
    },
  });
};
