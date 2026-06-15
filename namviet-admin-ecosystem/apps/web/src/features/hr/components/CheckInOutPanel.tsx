import React, { useState, useEffect } from 'react';
import { Button, Card, Spin, Alert } from 'antd';
import { MapPin, LogIn, LogOut, Loader2 } from 'lucide-react';
import { useCreateAttendanceLog, useUpdateAttendanceLog, useAttendanceLogs } from '../hooks';
import dayjs from 'dayjs';
import { toast } from 'sonner';

export const CheckInOutPanel = () => {
  const [location, setLocation] = useState<{ lat: number, lng: number } | null>(null);
  const [ip, setIp] = useState<string>('');
  const [loadingLocation, setLoadingLocation] = useState(false);
  
  const { data: logs, isLoading: loadingLogs } = useAttendanceLogs();
  const createLogMutation = useCreateAttendanceLog();
  const updateLogMutation = useUpdateAttendanceLog();

  // Find today's active log
  const todayStart = dayjs().startOf('day');
  const todayLog = logs?.find((log: any) => 
    dayjs(log.check_in_time).isAfter(todayStart) && !log.check_out_time
  );

  useEffect(() => {
    // Fetch IP
    fetch('https://api.ipify.org?format=json')
      .then(res => res.json())
      .then(data => setIp(data.ip))
      .catch(() => setIp('127.0.0.1'));
  }, []);

  const handleGetLocation = (): Promise<{lat: number, lng: number}> => {
    return new Promise((resolve, reject) => {
      if (!navigator.geolocation) {
        reject(new Error('Geolocation is not supported by your browser'));
      } else {
        navigator.geolocation.getCurrentPosition(
          (position) => {
            resolve({
              lat: position.coords.latitude,
              lng: position.coords.longitude
            });
          },
          () => {
            reject(new Error('Unable to retrieve your location'));
          }
        );
      }
    });
  };

  const handleCheckIn = async () => {
    try {
      setLoadingLocation(true);
      const coords = await handleGetLocation();
      setLocation(coords);
      
      createLogMutation.mutate({
        user_id: 'EMPLOYEE_1', // Mock user ID for demo
        check_in_time: dayjs().toISOString(),
        check_in_ip: ip,
        check_in_lat: coords.lat,
        check_in_lng: coords.lng,
      });
    } catch (error: any) {
      toast.error('Không thể lấy tọa độ. Vui lòng cấp quyền truy cập vị trí.');
    } finally {
      setLoadingLocation(false);
    }
  };

  const handleCheckOut = async () => {
    if (!todayLog?.id) return;
    
    try {
      setLoadingLocation(true);
      const coords = await handleGetLocation();
      setLocation(coords);
      
      updateLogMutation.mutate({
        id: todayLog.id,
        data: {
          check_out_time: dayjs().toISOString(),
          check_out_ip: ip,
          check_out_lat: coords.lat,
          check_out_lng: coords.lng,
          status: 'PRESENT'
        }
      });
    } catch (error: any) {
      toast.error('Không thể lấy tọa độ. Vui lòng cấp quyền truy cập vị trí.');
    } finally {
      setLoadingLocation(false);
    }
  };

  const isWorking = !!todayLog;

  return (
    <Card className="mb-6 shadow-sm border border-slate-100 bg-white">
      <div className="flex flex-col md:flex-row items-center justify-between gap-6">
        <div className="flex-1">
          <h2 className="text-lg font-bold text-slate-800 flex items-center gap-2 mb-2">
            <ClockIcon /> Bảng Điều khiển Chấm công
          </h2>
          <p className="text-slate-500 text-sm">
            Vui lòng cho phép quyền truy cập vị trí để ghi nhận Check-in/Out.
          </p>
          
          {(location || ip) && (
            <div className="mt-4 flex flex-col gap-1 text-xs text-slate-500 bg-slate-50 p-3 rounded-lg w-max">
              <div className="flex items-center gap-2">
                <MapPin className="w-3 h-3 text-blue-500" />
                Vị trí: {location ? `${location.lat.toFixed(5)}, ${location.lng.toFixed(5)}` : 'Đang chờ...'}
              </div>
              <div className="flex items-center gap-2 text-xs">
                <span className="font-mono bg-slate-200 px-1 rounded text-[10px]">IP</span>
                {ip || 'Đang tải...'}
              </div>
            </div>
          )}
        </div>

        <div className="flex flex-col gap-2 min-w-[200px]">
          <div className="text-center mb-2">
            <div className="text-3xl font-black text-slate-800 tracking-tight">
              {dayjs().format('HH:mm')}
            </div>
            <div className="text-xs font-medium text-slate-500 uppercase tracking-widest mt-1">
              {dayjs().format('dddd, DD/MM')}
            </div>
          </div>

          {!isWorking ? (
            <Button
              type="primary"
              size="large"
              className="bg-blue-600 h-14 text-base font-semibold w-full"
              icon={<LogIn className="w-5 h-5" />}
              onClick={handleCheckIn}
              loading={loadingLocation || createLogMutation.isPending}
            >
              Check-in VÀO CA
            </Button>
          ) : (
            <Button
              danger
              type="primary"
              size="large"
              className="h-14 text-base font-semibold w-full"
              icon={<LogOut className="w-5 h-5" />}
              onClick={handleCheckOut}
              loading={loadingLocation || updateLogMutation.isPending}
            >
              Check-out TAN CA
            </Button>
          )}

          {isWorking && (
            <div className="text-xs text-center text-emerald-600 mt-2 font-medium bg-emerald-50 py-1 rounded">
              Đang trong ca làm việc...
            </div>
          )}
        </div>
      </div>
    </Card>
  );
};

const ClockIcon = () => (
  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="text-blue-500">
    <circle cx="12" cy="12" r="10"></circle>
    <polyline points="12 6 12 12 16 14"></polyline>
  </svg>
);
