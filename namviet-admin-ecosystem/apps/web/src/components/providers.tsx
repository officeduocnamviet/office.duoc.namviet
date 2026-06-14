"use client";

import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AntdRegistry } from '@ant-design/nextjs-registry';
import { useState } from 'react';

export function Providers({ children }: { children: React.ReactNode }) {
  const [queryClient] = useState(() => new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60 * 1000, // 1 phút cache mặc định
        refetchOnWindowFocus: false, // Tránh fetch liên tục khi chuyển tab
      },
    },
  }));

  return (
    <QueryClientProvider client={queryClient}>
      <AntdRegistry>
        {children}
      </AntdRegistry>
    </QueryClientProvider>
  );
}
