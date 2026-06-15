import React from 'react';
import { Tabs, Table, Tag, Button, Space } from 'antd';
import { useConnections, useShippingPartners, useWebhookLogs } from '../hooks/useIntegrations';
import { ThirdPartyConnection, ShippingPartner, WebhookLog } from '../api/integrationApi';
import dayjs from 'dayjs';

export const IntegrationBoard = () => {
  const { data: connections, isLoading: loadingConn } = useConnections();
  const { data: partners, isLoading: loadingPartners } = useShippingPartners();
  const { data: webhooks, isLoading: loadingWebhooks } = useWebhookLogs();

  const partnerColumns = [
    { title: 'Mã đối tác', dataIndex: 'code', key: 'code', render: (text: string) => <strong>{text}</strong> },
    { title: 'Tên đối tác', dataIndex: 'name', key: 'name' },
    { title: 'Loại', dataIndex: 'partner_type', key: 'type' },
    { 
      title: 'Trạng thái', 
      dataIndex: 'status', 
      key: 'status',
      render: (status: string) => <Tag color={status === 'ACTIVE' ? 'success' : 'default'}>{status}</Tag>
    },
    { title: 'Tracking URL', dataIndex: 'tracking_url_template', key: 'url', render: (url: string) => <span className="text-xs text-blue-500">{url}</span> }
  ];

  const connColumns = [
    { title: 'Mã đối tác', dataIndex: 'partner_code', key: 'partner_code' },
    { title: 'Loại kết nối', dataIndex: 'connection_type', key: 'type' },
    { 
      title: 'Trạng thái', 
      dataIndex: 'status', 
      key: 'status',
      render: (status: string) => <Tag color={status === 'ACTIVE' ? 'success' : 'error'}>{status}</Tag>
    },
    { 
      title: 'Ngày tạo', 
      dataIndex: 'created_at', 
      key: 'created_at',
      render: (date: string) => date ? dayjs(date).format('DD/MM/YYYY HH:mm') : '---'
    }
  ];

  const webhookColumns = [
    { title: 'Sự kiện', dataIndex: 'event_type', key: 'event_type', render: (text: string) => <Tag color="blue">{text}</Tag> },
    { title: 'Connection ID', dataIndex: 'connection_id', key: 'conn_id', render: (id: string) => <span className="text-xs">{id}</span> },
    { 
      title: 'Status Code', 
      dataIndex: 'response_status', 
      key: 'status',
      render: (status: number) => <Tag color={status >= 200 && status < 300 ? 'success' : 'error'}>{status}</Tag>
    },
    { 
      title: 'Thời gian', 
      dataIndex: 'created_at', 
      key: 'created_at',
      render: (date: string) => date ? dayjs(date).format('DD/MM/YYYY HH:mm:ss') : '---'
    },
    {
      title: 'Chi tiết',
      key: 'action',
      render: () => <Button type="link" size="small">Xem payload</Button>
    }
  ];

  const items = [
    {
      key: 'partners',
      label: 'Đối tác vận chuyển',
      children: (
        <Table 
          columns={partnerColumns} 
          dataSource={partners} 
          rowKey="id" 
          loading={loadingPartners} 
          pagination={false}
          size="middle"
        />
      )
    },
    {
      key: 'connections',
      label: 'Kết nối API (Connections)',
      children: (
        <Table 
          columns={connColumns} 
          dataSource={connections} 
          rowKey="id" 
          loading={loadingConn} 
          pagination={false}
          size="middle"
        />
      )
    },
    {
      key: 'webhooks',
      label: 'Nhật ký Webhook',
      children: (
        <Table 
          columns={webhookColumns} 
          dataSource={webhooks} 
          rowKey="id" 
          loading={loadingWebhooks} 
          size="middle"
        />
      )
    }
  ];

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <Tabs defaultActiveKey="partners" items={items} />
    </div>
  );
};
