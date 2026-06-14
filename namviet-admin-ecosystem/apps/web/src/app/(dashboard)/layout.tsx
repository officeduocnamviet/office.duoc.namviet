"use client";

import React, { useState } from 'react';
import { Layout } from 'antd';
import { Header } from '@/components/layout/Header';
import { GlobalSearchModal } from '@/components/ui/GlobalSearchModal';
import MobileBottomNav from '@/components/layout/MobileBottomNav';
import { useFCM } from '@/hooks/useFCM';

const { Content } = Layout;

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  // Kích hoạt Firebase Cloud Messaging Hook
  useFCM();

  // State cho Mobile Bottom Nav
  const [activeTab, setActiveTab] = useState('home');
  const [isFabOpen, setIsFabOpen] = useState(false);

  return (
    <Layout className="min-h-screen bg-gray-50">
      <Header />
      <GlobalSearchModal />
      
      {/* Content tràn viền, loại bỏ margin cứng, dùng padding siêu nhỏ (px-2 md:px-4) */}
      <Content className="w-full max-w-7xl mx-auto px-2 sm:px-4 py-4 md:py-6 pb-24 md:pb-6">
        {children}
      </Content>

      <MobileBottomNav 
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        cartCount={3}
        isFabOpen={isFabOpen}
        toggleFabMenu={() => setIsFabOpen(!isFabOpen)}
      />
    </Layout>
  );
}
