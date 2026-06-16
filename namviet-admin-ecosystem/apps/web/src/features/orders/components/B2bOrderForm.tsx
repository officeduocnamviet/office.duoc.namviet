import React, { useState } from "react";
import {
  Form,
  Input,
  Select,
  Button,
  Row,
  Col,
  Typography,
  Card,
  Space,
  Divider,
  message,
  InputNumber,
  Table,
} from "antd";
import {
  MinusCircleOutlined,
  PlusOutlined,
  ShoppingCartOutlined,
} from "@ant-design/icons";
import { Building2, Truck, Percent } from "lucide-react";
import { useCreateOrder } from "../hooks";
import { OrdersCreateOrderRequest } from "@namviet/shared-types/src/backend.d";
import { useForm, useFieldArray } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import * as z from "zod";
import { Box } from "lucide-react";

const { Text, Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const orderItemSchema = z.object({
  product_id: z.number({ message: "Chọn sản phẩm" }),
  batch_no: z.string().optional(),
  quantity: z.number().min(1, "Số lượng tối thiểu là 1"),
  unit_price: z.number().min(0, "Giá không hợp lệ"),
  uom: z.string().min(1, "Đơn vị tính"),
  discount: z.number().optional(),
});

const orderSchema = z.object({
  customer_id: z.number({ message: "Bắt buộc chọn Đối tác/Doanh nghiệp" }),
  warehouse_id: z.number({ message: "Bắt buộc chọn Chi nhánh xuất kho" }),
  note: z.string().optional(),
  items: z.array(orderItemSchema).min(1, "Đơn hàng phải có ít nhất 1 sản phẩm"),
});

type OrderFormValues = z.infer<typeof orderSchema>;

export const B2bOrderForm = () => {
  const createMutation = useCreateOrder();

  const {
    control,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm<OrderFormValues>({
    resolver: zodResolver(orderSchema),
    defaultValues: {
      items: [
        {
          product_id: 1,
          quantity: 10,
          unit_price: 45000,
          uom: "Thùng",
          discount: 0,
        },
      ],
    },
  });

  const { fields, append, remove } = useFieldArray({
    control,
    name: "items",
  });

  const items = watch("items");

  // Tính tổng tiền realtime
  const totalAmount = items.reduce(
    (sum, item) => sum + item.quantity * item.unit_price - (item.discount || 0),
    0,
  );

  const onSubmit = (data: OrderFormValues) => {
    const payload: OrdersCreateOrderRequest = {
      code: `B2B-${Date.now()}`,
      order_type: "B2B",
      customer_id: data.customer_id,
      note: data.note,
      items: data.items ? data.items.map((i) => ({
        product_id: i.product_id,
        batch_no: i.batch_no,
        quantity: i.quantity,
        unit_price: i.unit_price,
        uom: i.uom,
        discount: i.discount || 0,
      })) : [],
    };

    createMutation.mutate(payload, {
      onSuccess: () => {
        message.success("Tạo đơn hàng B2B thành công!");
      },
      onError: () => {
        message.error("Lỗi khi tạo đơn hàng");
      },
    });
  };

  return (
    <div className="flex flex-col xl:flex-row gap-6">
      {/* Cột trái: Thông tin Đối tác */}
      <div className="w-full xl:w-1/4 space-y-6">
        <Card
          title={
            <>
              <Building2 className="inline mr-2" />
              Thông tin B2B
            </>
          }
          className="shadow-sm border-blue-100"
        >
          <div className="mb-4">
            <Text className="text-gray-500 mb-1 block font-medium">
              Khách hàng Doanh nghiệp <span className="text-red-500">*</span>
            </Text>
            <Select
              className="w-full"
              placeholder="Chọn hoặc tìm kiếm..."
              onChange={(val) => setValue("customer_id", val)}
              status={errors.customer_id ? "error" : ""}
              options={[
                { label: "Công ty Dược phẩm A (MST: 010101)", value: 101 },
                { label: "Nhà thuốc An Khang (MST: 020202)", value: 102 },
              ]}
            />
            {errors.customer_id && (
              <Text type="danger" className="text-xs">
                {errors.customer_id.message}
              </Text>
            )}
          </div>
          <div className="mb-4">
            <Text className="text-gray-500 mb-1 block font-medium">
              Chi nhánh Xuất kho <span className="text-red-500">*</span>
            </Text>
            <Select
              className="w-full"
              placeholder="Chọn kho xuất..."
              onChange={(val) => setValue("warehouse_id", val)}
              status={errors.warehouse_id ? "error" : ""}
              options={[
                { label: "Kho Tổng Miền Bắc", value: 1 },
                { label: "Kho Trung Chuyển HN", value: 2 },
              ]}
            />
            {errors.warehouse_id && (
              <Text type="danger" className="text-xs">
                {errors.warehouse_id.message}
              </Text>
            )}
          </div>
          <div>
            <Text className="text-gray-500 mb-1 block font-medium">
              Ghi chú (Vận đơn / Thanh toán)
            </Text>
            <TextArea
              rows={4}
              placeholder="VD: Giao hàng vào thứ 7, thu COD..."
              onChange={(e) => setValue("note", e.target.value)}
            />
          </div>
        </Card>

        <Card className="shadow-sm bg-gradient-to-br from-slate-800 to-slate-900 border-none text-white">
          <div className="mb-2 opacity-80 flex items-center gap-2">
            <Box size={16} /> Tổng giá trị đơn hàng (B2B)
          </div>
          <div className="text-3xl font-bold mb-6 text-emerald-400">
            {new Intl.NumberFormat("vi-VN", {
              style: "currency",
              currency: "VND",
            }).format(totalAmount)}
          </div>
          <Button
            size="large"
            className="w-full bg-blue-500 hover:bg-blue-400 border-none text-white font-semibold"
            onClick={handleSubmit(onSubmit)}
            loading={createMutation.isPending}
          >
            Lưu & Duyệt Đơn
          </Button>
        </Card>
      </div>

      {/* Cột phải: Danh sách Sản phẩm */}
      <div className="w-full xl:w-3/4">
        <Card
          title={
            <>
              <ShoppingCartOutlined className="mr-2" />
              Giỏ hàng Sỉ (Order Items)
            </>
          }
          className="shadow-sm min-h-[500px]"
        >
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-slate-50 border-y border-gray-200 text-sm text-gray-600 font-medium">
                  <th className="p-3 w-10">#</th>
                  <th className="p-3 min-w-[200px]">Sản phẩm & Lô</th>
                  <th className="p-3 w-28">Số lượng</th>
                  <th className="p-3 w-24">ĐVT (Sỉ)</th>
                  <th className="p-3 w-32">Đơn giá</th>
                  <th className="p-3 w-32">Ck/Tặng</th>
                  <th className="p-3 w-32 text-right">Thành tiền</th>
                  <th className="p-3 w-10 text-center">
                    <PlusOutlined />
                  </th>
                </tr>
              </thead>
              <tbody>
                {fields.map((field, index) => (
                  <tr
                    key={field.id}
                    className="border-b border-gray-100 last:border-0 hover:bg-blue-50/30 transition-colors"
                  >
                    <td className="p-3 text-gray-400">{index + 1}</td>
                    <td className="p-3">
                      <Select
                        className="w-full mb-2"
                        placeholder="Chọn SP..."
                        value={watch(`items.${index}.product_id`)}
                        onChange={(val) =>
                          setValue(`items.${index}.product_id`, val)
                        }
                        options={[
                          { label: "Viên Uống ABC 500mg", value: 1 },
                          { label: "Siro Trẻ em XYZ", value: 2 },
                        ]}
                      />
                      <Select
                        className="w-full"
                        placeholder="Chọn Lô (Tùy chọn)"
                        size="small"
                        allowClear
                        value={watch(`items.${index}.batch_no`)}
                        onChange={(val) =>
                          setValue(`items.${index}.batch_no`, val)
                        }
                        options={[
                          {
                            label: "Lô: BATCH001 (HSD: 12/2026)",
                            value: "BATCH001",
                          },
                          {
                            label: "Lô: BATCH002 (HSD: 01/2027)",
                            value: "BATCH002",
                          },
                        ]}
                      />
                    </td>
                    <td className="p-3">
                      <InputNumber
                        min={1}
                        className="w-full"
                        value={watch(`items.${index}.quantity`)}
                        onChange={(val) =>
                          setValue(`items.${index}.quantity`, val || 1)
                        }
                      />
                    </td>
                    <td className="p-3">
                      <Input
                        value={watch(`items.${index}.uom`)}
                        onChange={(e) =>
                          setValue(`items.${index}.uom`, e.target.value)
                        }
                      />
                    </td>
                    <td className="p-3">
                      <InputNumber
                        className="w-full"
                        formatter={(value) =>
                          `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ",")
                        }
                        value={watch(`items.${index}.unit_price`)}
                        onChange={(val) =>
                          setValue(`items.${index}.unit_price`, val || 0)
                        }
                      />
                    </td>
                    <td className="p-3">
                      <InputNumber
                        className="w-full"
                        prefix={<Percent size={14} className="text-gray-400" />}
                        formatter={(value) =>
                          `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ",")
                        }
                        value={watch(`items.${index}.discount`)}
                        onChange={(val) =>
                          setValue(`items.${index}.discount`, val || 0)
                        }
                      />
                    </td>
                    <td className="p-3 text-right font-semibold text-blue-700">
                      {new Intl.NumberFormat("vi-VN", {
                        style: "currency",
                        currency: "VND",
                      }).format(
                        watch(`items.${index}.quantity`) *
                          watch(`items.${index}.unit_price`) -
                          (watch(`items.${index}.discount`) || 0),
                      )}
                    </td>
                    <td className="p-3 text-center">
                      <Button
                        type="text"
                        danger
                        icon={<MinusCircleOutlined />}
                        onClick={() => remove(index)}
                      />
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <Button
            type="dashed"
            onClick={() =>
              append({
                product_id: 1,
                quantity: 1,
                unit_price: 0,
                uom: "Thùng",
                discount: 0,
              })
            }
            block
            icon={<PlusOutlined />}
            className="mt-6 border-blue-200 text-blue-500 hover:text-blue-600 hover:border-blue-500 bg-blue-50/50"
          >
            Thêm sản phẩm
          </Button>
          {errors.items && (
            <Text type="danger" className="block mt-2 font-medium">
              {errors.items.message}
            </Text>
          )}
        </Card>
      </div>
    </div>
  );
};
