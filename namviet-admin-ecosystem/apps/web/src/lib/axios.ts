import axios from 'axios';
import { useAuthStore } from '@/stores/useAuthStore';

// Cấu hình BaseURL, hướng vào Backend Golang
export const apiClient = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor: Trước khi Request bay đi
apiClient.interceptors.request.use(
  (config) => {
    // Lấy Token từ Zustand
    const state = useAuthStore.getState();
    if (state.token) {
      config.headers.Authorization = `Bearer ${state.token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Interceptor: Xử lý Response trả về
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // LUÔN in ra console trước khi thực hiện bất kỳ lệnh redirect nào
    console.error("API Error: ", error.response || error.message);
    
    if (error.response) {
      // CHỈ đá văng ra Login nếu API trả về đúng lỗi 401 Unauthorized
      if (error.response.status === 401) {
        alert("Phiên đăng nhập hết hạn hoặc không hợp lệ. Vui lòng đăng nhập lại!");
        const { logout } = useAuthStore.getState();
        logout(); // Xóa state
        window.location.href = '/login'; 
      }
    } else {
      // Báo lỗi mạng (Ví dụ: Server sập, mất mạng, lỗi CORS)
      alert("Lỗi kết nối tới máy chủ (Network Error)!");
    }
    return Promise.reject(error);
  }
);

// Interceptor: Khi Response trả về
apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    // Nếu lỗi 401 Unauthorized (Token hết hạn hoặc sai)
    if (error.response?.status === 401) {
      useAuthStore.getState().logout();
      if (typeof window !== 'undefined') {
        window.location.href = '/login';
      }
    }
    return Promise.reject(error);
  }
);
