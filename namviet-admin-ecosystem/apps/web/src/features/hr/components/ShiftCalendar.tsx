import React, { useState } from 'react';
import { Calendar, Badge, Modal, Form, Select, Button } from 'antd';
import { useWorkShifts, useShiftAssignments, useCreateShiftAssignment } from '../hooks';
import dayjs, { Dayjs } from 'dayjs';

export const ShiftCalendar = () => {
  const { data: shifts } = useWorkShifts();
  const { data: assignments } = useShiftAssignments();
  const createAssignmentMutation = useCreateShiftAssignment();
  
  const [selectedDate, setSelectedDate] = useState<Dayjs | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [form] = Form.useForm();

  const getListData = (value: Dayjs) => {
    const dateStr = value.format('YYYY-MM-DD');
    return assignments?.filter((a: any) => dayjs(a.work_date).format('YYYY-MM-DD') === dateStr) || [];
  };

  const dateCellRender = (value: Dayjs) => {
    const listData = getListData(value);
    return (
      <ul className="events p-0 m-0 list-none">
        {listData.map((item: any) => {
          const shift = shifts?.find((s: any) => s.id === item.shift_id);
          return (
            <li key={item.id} className="text-xs mb-1">
              <Badge 
                status="success" 
                text={`${shift?.name || 'Ca làm'} - NV: ${item.user_id}`} 
                className="truncate w-full block text-xs" 
              />
            </li>
          );
        })}
      </ul>
    );
  };

  const onSelect = (newValue: Dayjs) => {
    setSelectedDate(newValue);
    form.setFieldsValue({ work_date: newValue });
    setIsModalOpen(true);
  };

  const handleFinish = (values: any) => {
    createAssignmentMutation.mutate({
      shift_id: values.shift_id,
      user_id: values.user_id,
      work_date: values.work_date.toISOString(),
    }, {
      onSuccess: () => setIsModalOpen(false)
    });
  };

  return (
    <div className="bg-white p-6 rounded-2xl shadow-sm border border-slate-100">
      <Calendar cellRender={dateCellRender} onSelect={onSelect} />

      <Modal 
        title={`Phân ca - ${selectedDate?.format('DD/MM/YYYY')}`}
        open={isModalOpen} 
        onCancel={() => setIsModalOpen(false)}
        footer={null}
        destroyOnClose
      >
        <Form form={form} layout="vertical" onFinish={handleFinish} className="mt-4">
          <Form.Item name="work_date" hidden>
            <Select />
          </Form.Item>
          
          <Form.Item 
            name="shift_id" 
            label="Chọn ca làm việc"
            rules={[{ required: true, message: 'Vui lòng chọn ca' }]}
          >
            <Select>
              {shifts?.map((s: any) => (
                <Select.Option key={s.id} value={s.id}>
                  {s.name} ({s.start_time} - {s.end_time})
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item 
            name="user_id" 
            label="Mã nhân viên"
            rules={[{ required: true, message: 'Vui lòng nhập nhân viên' }]}
          >
            {/* In a real app this would be an Employee Select */}
            <Select mode="tags" placeholder="Nhập mã nhân viên (Enter để thêm)" />
          </Form.Item>

          <div className="flex justify-end gap-2 mt-6">
            <Button onClick={() => setIsModalOpen(false)}>Hủy</Button>
            <Button type="primary" htmlType="submit" className="bg-blue-600" loading={createAssignmentMutation.isPending}>
              Phân ca
            </Button>
          </div>
        </Form>
      </Modal>
    </div>
  );
};
