"use client";

import React from 'react';
import { Card, Row, Col, Typography, List, Badge, Calendar, Tag } from 'antd';
import { NotificationOutlined, TrophyOutlined, CalendarOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

export const CommonBoard = () => {
  const news = [
    { title: "Chào đón 2 bác sĩ mới tham gia phòng khám", time: "2 giờ trước" },
    { title: "Cập nhật chính sách chiết khấu T8/2026", time: "Hôm qua" },
  ];

  const honors = [
    { name: "Dược sĩ Nguyễn Thị Lan", reason: "Nhận 5 sao tư vấn", date: "Hôm nay" },
    { name: "Team Marketing", reason: "Đạt 200 lượt đặt lịch tuần", date: "Hôm qua" },
  ];

  return (
    <div className="mb-6">
      <Title level={4} className="mb-4">🌟 Quảng trường Chung (Toàn Công ty)</Title>
      <Row gutter={[16, 16]}>
        {/* Tin tức & Thông báo */}
        <Col xs={24} md={8}>
          <Card 
            title={<><NotificationOutlined className="mr-2 text-blue-500" />Tin tức nội bộ</>}
            className="h-full border-2 border-gray-200"
            styles={{ body: { padding: '12px 24px' } }}
          >
            <List
              dataSource={news}
              renderItem={(item) => (
                <List.Item className="px-0">
                  <List.Item.Meta
                    title={<a href="#">{item.title}</a>}
                    description={item.time}
                  />
                </List.Item>
              )}
            />
          </Card>
        </Col>

        {/* Vinh danh */}
        <Col xs={24} md={8}>
          <Card 
            title={<><TrophyOutlined className="mr-2 text-yellow-500" />Vinh danh cá nhân xuất sắc</>}
            className="h-full border-2 border-gray-200 bg-yellow-50"
            styles={{ body: { padding: '12px 24px' } }}
          >
            <List
              dataSource={honors}
              renderItem={(item) => (
                <List.Item className="px-0">
                  <List.Item.Meta
                    title={<Text strong>{item.name}</Text>}
                    description={<Text type="secondary">{item.reason}</Text>}
                  />
                  <Tag color="gold">{item.date}</Tag>
                </List.Item>
              )}
            />
          </Card>
        </Col>

        {/* Lịch trực nhanh */}
        <Col xs={24} md={8}>
          <Card 
            title={<><CalendarOutlined className="mr-2 text-green-500" />Lịch trực hôm nay</>}
            className="h-full border-2 border-gray-200"
          >
            <div className="flex flex-col gap-2">
              <div className="flex justify-between border-b pb-2">
                <Text strong>Khám Nhi:</Text>
                <Badge status="success" text="BS. Trần Thị Lan" />
              </div>
              <div className="flex justify-between border-b pb-2">
                <Text strong>Khám Nội:</Text>
                <Badge status="success" text="BS. Nguyễn Văn Minh" />
              </div>
              <div className="flex justify-between">
                <Text strong>Nhà thuốc:</Text>
                <Badge status="processing" text="DS. Lê Văn C (Ca Sáng)" />
              </div>
            </div>
          </Card>
        </Col>
      </Row>
    </div>
  );
};
