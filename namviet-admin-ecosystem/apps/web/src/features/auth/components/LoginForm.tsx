"use client";

import React from 'react';
import { Form, Input, Button, Card, Typography } from 'antd';
import { MailOutlined, LockOutlined } from '@ant-design/icons'; // Dùng tạm icon mặc định cho form login
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { useLogin } from '../hooks/useLogin';
import { useRouter } from 'next/navigation';

const loginSchema = z.object({
  email: z.string().email('Email không hợp lệ'),
  password: z.string().min(6, 'Mật khẩu phải có ít nhất 6 ký tự'),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export const LoginForm = () => {
  const { control, handleSubmit, formState: { errors } } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { email: '', password: '' },
  });

  const loginMutation = useLogin();
  const router = useRouter();

  const onSubmit = (data: LoginFormValues) => {
    loginMutation.mutate(data, {
      onSuccess: () => {
        router.push('/');
      }
    });
  };

  return (
    <Card className="w-full max-w-md shadow-2xl border-2 border-gray-100 rounded-xl">
      <div className="text-center mb-8 flex flex-col items-center">
        <img src="/logo.png" alt="Nam Việt Logo" className="h-16 mb-2 object-contain" />
        <Typography.Title level={2} className="!mb-1 text-orange-600">Nam Việt ERP</Typography.Title>
        <Typography.Text type="secondary">Đăng nhập hệ thống quản trị</Typography.Text>
      </div>

      <Form layout="vertical" onFinish={handleSubmit(onSubmit)}>
        <Form.Item 
          validateStatus={errors.email ? 'error' : ''} 
          help={errors.email?.message}
        >
          <Controller
            name="email"
            control={control}
            render={({ field }) => (
              <Input 
                {...field} 
                size="large" 
                prefix={<MailOutlined className="text-gray-400" />} 
                placeholder="Email (vd: admin@namviet.com)" 
              />
            )}
          />
        </Form.Item>

        <Form.Item 
          validateStatus={errors.password ? 'error' : ''} 
          help={errors.password?.message}
        >
          <Controller
            name="password"
            control={control}
            render={({ field }) => (
              <Input.Password 
                {...field} 
                size="large" 
                prefix={<LockOutlined className="text-gray-400" />} 
                placeholder="Mật khẩu" 
              />
            )}
          />
        </Form.Item>

        <Form.Item className="mt-6 mb-0">
          <Button 
            type="primary" 
            htmlType="submit" 
            size="large" 
            block 
            loading={loginMutation.isPending}
            className="bg-blue-600 hover:bg-blue-700"
          >
            Đăng nhập
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
};
