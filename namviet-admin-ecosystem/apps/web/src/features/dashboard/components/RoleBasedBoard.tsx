"use client";

import React, { useState } from 'react';
import { Card, Row, Col, Typography, Segmented, Table, Tag, Button, Badge } from 'antd';
import { MedicineBoxOutlined, UserOutlined, ClockCircleOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

export const RoleBasedBoard = () => {
  const [role, setRole] = useState<string>('Lễ tân');

  const renderContent = () => {
    switch (role) {
      case 'Lễ tân':
        return (
          <Card title="Trung tâm Điều phối Lịch Hẹn" className="border-2 border-blue-200">
            <Row gutter={[16, 16]}>
              <Col span={8}>
                <Card type="inner" title="BS. Nguyễn Văn Minh (Khám Tổng Quát)">
                  <div className="p-2 mb-2 bg-green-100 border border-green-300 rounded">
                    <Text strong>Nguyễn Văn Bình (09:00)</Text><br/>
                    <Tag color="green" className="mt-1">Đã Check-in</Tag>
                  </div>
                  <div className="p-2 bg-gray-100 border border-gray-300 rounded">
                    <Text strong>Lê Thị Hoa (09:30)</Text><br/>
                    <Tag className="mt-1">Chưa xác nhận</Tag>
                  </div>
                </Card>
              </Col>
              <Col span={8}>
                <Card type="inner" title="BS. Trần Thị Lan (Khám Nhi)">
                  <div className="p-2 mb-2 bg-yellow-100 border border-yellow-300 rounded">
                    <Text strong>Bé An (08:30)</Text><br/>
                    <Tag color="orange" className="mt-1">Đang khám</Tag>
                  </div>
                </Card>
              </Col>
              <Col span={8}>
                <Card type="inner" title="Phòng Tiêm Chủng">
                  <div className="p-2 bg-blue-100 border border-blue-300 rounded">
                    <Text strong>Trần Văn C (10:00)</Text><br/>
                    <Tag color="blue" className="mt-1">Đã xác nhận</Tag>
                  </div>
                </Card>
              </Col>
            </Row>
          </Card>
        );
      case 'Dược sĩ':
        return (
          <Card title="Buồng lái Tư vấn & Bán hàng" className="border-2 border-green-200">
             <Row gutter={[16, 16]}>
              <Col span={12}>
                <Card type="inner" title="Đơn thuốc chờ xử lý từ Phòng Khám">
                  <div className="p-3 mb-2 bg-purple-50 border border-purple-200 rounded flex justify-between items-center">
                    <div>
                      <Text strong className="text-lg">Bệnh nhân: Nguyễn Văn Bình</Text><br/>
                      <Text type="secondary">BS. Minh chuyển đến - 2 Phút trước</Text>
                    </div>
                    <Button type="primary">Lấy Thuốc Ngay</Button>
                  </div>
                </Card>
              </Col>
              <Col span={12}>
                <Card type="inner" title="Nhiệm vụ Chăm sóc Khách hàng Hôm nay">
                  <ul className="pl-4">
                    <li className="mb-2"><strong>[Gọi điện]</strong> Chị Mai - Hỏi thăm tình hình bé uống Siro Prospan (Bán 3 ngày trước)</li>
                    <li><strong>[Gửi Zalo]</strong> Anh Hưng - Nhắc lịch tiêm nhắc lại Vaccine Viêm gan B</li>
                  </ul>
                </Card>
              </Col>
            </Row>
          </Card>
        );
      case 'Bác sĩ':
        return (
          <Card title="Trung tâm Chỉ huy Lâm sàng" className="border-2 border-red-200">
             <Row gutter={[16, 16]}>
              <Col span={12}>
                <Card type="inner" title="Hàng đợi Bệnh nhân (Đã Check-in)">
                  <Table 
                    dataSource={[
                      { key: '1', name: 'Nguyễn Văn Bình', time: '09:00', status: 'Đang chờ' }
                    ]}
                    columns={[
                      { title: 'Tên Bệnh nhân', dataIndex: 'name', key: 'name' },
                      { title: 'Giờ hẹn', dataIndex: 'time', key: 'time' },
                      { title: 'Trạng thái', dataIndex: 'status', key: 'status', render: (text) => <Tag color="green">{text}</Tag> },
                      { title: 'Thao tác', key: 'action', render: () => <Button type="link">Gọi vào khám</Button> }
                    ]}
                    pagination={false}
                  />
                </Card>
              </Col>
              <Col span={12}>
                <Card type="inner" title="Hộp Thư Kết Quả Cận Lâm Sàng" extra={<Badge count={1} status="error" />}>
                  <div className="p-3 bg-red-50 border-l-4 border-red-500 rounded mb-2">
                    <Text strong className="text-red-600">🔴 KHẨN CẤP: Glucose 150 mg/dL (Cao)</Text><br/>
                    <Text type="secondary">BN Nguyễn Văn A - Lấy mẫu 30p trước</Text>
                  </div>
                  <div className="p-3 bg-gray-50 border-l-4 border-gray-300 rounded">
                    <Text strong>⚪ BÌNH THƯỜNG: Siêu âm ổ bụng</Text><br/>
                    <Text type="secondary">BN Trần Thị B - Lấy mẫu 1h trước</Text>
                  </div>
                </Card>
              </Col>
             </Row>
          </Card>
        );
      default:
        return null;
    }
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-4">
        <Title level={4} className="!mb-0">🚀 Buồng Lái Cá Nhân (Role-based)</Title>
        <Segmented 
          options={['Lễ tân', 'Dược sĩ', 'Bác sĩ']} 
          value={role} 
          onChange={(value) => setRole(value as string)} 
        />
      </div>
      {renderContent()}
    </div>
  );
};
