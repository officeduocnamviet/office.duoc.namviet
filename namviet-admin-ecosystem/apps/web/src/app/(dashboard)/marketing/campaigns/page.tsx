"use client";

import React from 'react';
import { CampaignTable } from '@/features/marketing/components/CampaignTable';
import { Megaphone } from 'lucide-react';

export default function CampaignsPage() {
  return (
    <div className="p-4 md:p-8 space-y-6 max-w-7xl mx-auto">
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold text-slate-800 flex items-center gap-2">
          <Megaphone className="w-6 h-6 text-blue-600" />
          Chiến dịch Marketing (Campaigns)
        </h1>
        <p className="text-slate-500">
          Quản lý các chiến dịch khuyến mãi, SMS marketing, email marketing đến tập khách hàng mục tiêu.
        </p>
      </div>

      <CampaignTable />
    </div>
  );
}
