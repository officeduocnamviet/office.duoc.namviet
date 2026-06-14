import { useEffect } from 'react';
import { messaging } from '@/lib/firebase';
import { getToken, onMessage } from 'firebase/messaging';
import { apiClient } from '@/lib/axios';
import { useAuthStore } from '@/stores/useAuthStore';
import { toast } from 'sonner';

export const useFCM = () => {
  const token = useAuthStore((state) => state.token);

  useEffect(() => {
    if (!token) return;

    const requestPermissionAndGetToken = async () => {
      try {
        const permission = await Notification.requestPermission();
        if (permission === 'granted') {
          const msg = await messaging();
          if (msg) {
            // Lấy FCM Token từ Firebase
            const currentToken = await getToken(msg, {
              // VAPID KEY lấy từ Firebase Console (Project Settings > Cloud Messaging > Web Push certs)
              // Tạm thời dùng VAPID trống nếu chưa có, Firebase vẫn tự lấy token cho test
            });

            if (currentToken) {
              // Gửi Token lên Backend
              await apiClient.post('/users/me/fcm-token', {
                token: currentToken,
                device_type: navigator.userAgent
              });
              console.log('FCM Token đã được lưu trữ thành công');
            }
          }
        }
      } catch (error) {
        console.error('Lỗi khi lấy FCM token:', error);
      }
    };

    requestPermissionAndGetToken();

    // Lắng nghe thông báo khi ứng dụng đang mở (Foreground)
    const setupListener = async () => {
      const msg = await messaging();
      if (msg) {
        onMessage(msg, (payload) => {
          console.log('Message received. ', payload);
          toast.success(`${payload.notification?.title}: ${payload.notification?.body}`);
        });
      }
    };
    
    setupListener();

  }, [token]);
};
