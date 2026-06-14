import React, { useEffect } from "react";
import { useForm, Controller } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  WarehouseFormData,
  warehouseSchema,
  Warehouse,
} from "@namviet/shared-types/src/warehouse.types";
import { ChevronLeft, Save, Building, MapPin, Phone, Hash } from "lucide-react";
import { useRouter } from "next/navigation";

interface WarehouseFormProps {
  initialData?: Warehouse;
  onSubmit: (data: WarehouseFormData) => void;
  isLoading: boolean;
}

export function WarehouseForm({
  initialData,
  onSubmit,
  isLoading,
}: WarehouseFormProps) {
  const router = useRouter();
  const isEditing = !!initialData;

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<WarehouseFormData>({
    resolver: zodResolver(warehouseSchema) as any,
    defaultValues: {
      company_id: "",
      key: "",
      code: "",
      name: "",
      type: "retail",
      address: "",
      manager: "",
      phone: "",
      status: "active",
    },
  });

  useEffect(() => {
    if (initialData) {
      reset({
        company_id: initialData.company_id || "",
        key: initialData.key,
        code: initialData.code || "",
        name: initialData.name,
        type: initialData.type,
        address: initialData.address || "",
        manager: initialData.manager || "",
        phone: initialData.phone || "",
        latitude: initialData.latitude || undefined,
        longitude: initialData.longitude || undefined,
        status: initialData.status,
      });
    }
  }, [initialData, reset]);

  return (
    <div className="max-w-4xl mx-auto pb-12">
      {/* Sticky Header */}
      <div className="flex items-center justify-between bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 shadow-sm sticky top-0 z-10 mb-6">
        <div className="flex items-center gap-4">
          <button
            type="button"
            onClick={() => router.push("/warehouses")}
            className="w-10 h-10 flex items-center justify-center rounded-full bg-slate-50 hover:bg-slate-200 text-slate-600 transition-colors"
          >
            <ChevronLeft size={20} />
          </button>
          <div>
            <h2 className="text-xl font-black text-slate-800 flex items-center gap-2">
              <Building className="text-primary" size={24} />
              {isEditing ? "Cập nhật Chi nhánh" : "Thêm mới Chi nhánh"}
            </h2>
            <p className="text-xs font-medium text-slate-500 mt-0.5">
              Điền đầy đủ thông tin để {isEditing ? "cập nhật" : "tạo"} kho/cửa
              hàng
            </p>
          </div>
        </div>
        <button
          onClick={handleSubmit((data) => {
            const formattedData = {
              ...data,
              company_id:
                data.company_id?.trim() === "" ? null : data.company_id,
            };
            onSubmit(formattedData);
          })}
          disabled={isLoading}
          className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white px-6 py-2.5 rounded-xl font-bold transition-all shadow-md active:scale-95 disabled:opacity-50"
        >
          {isLoading ? (
            <span className="animate-spin">⏳</span>
          ) : (
            <Save size={18} />
          )}
          Lưu thông tin
        </button>
      </div>

      <form className="space-y-6">
        <div className="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-lg font-bold text-slate-800 mb-4 border-b pb-2">
            Thông tin Cơ bản
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Mã định danh (Key) *
              </label>
              <input
                {...register("key")}
                className={`w-full p-2.5 bg-slate-50 border rounded-xl outline-none focus:border-primary ${errors.key ? "border-red-500" : "border-slate-200"}`}
                placeholder="VD: CN_HCM_01"
              />
              {errors.key && (
                <p className="text-red-500 text-xs mt-1">
                  {errors.key.message}
                </p>
              )}
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Tên Chi nhánh *
              </label>
              <input
                {...register("name")}
                className={`w-full p-2.5 bg-slate-50 border rounded-xl outline-none focus:border-primary ${errors.name ? "border-red-500" : "border-slate-200"}`}
                placeholder="VD: Cửa hàng Quận 1"
              />
              {errors.name && (
                <p className="text-red-500 text-xs mt-1">
                  {errors.name.message}
                </p>
              )}
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Mã nội bộ (Code)
              </label>
              <input
                {...register("code")}
                className="w-full p-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
                placeholder="VD: WH-001"
              />
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Loại hình
              </label>
              <select
                {...register("type")}
                className="w-full p-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
              >
                <option value="retail">Cửa hàng Bán lẻ</option>
                <option value="wholesale">Kho Sỉ</option>
                <option value="warehouse">Tổng Kho</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Công ty trực thuộc (ID)
              </label>
              <input
                {...register("company_id")}
                className="w-full p-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
                placeholder="UUID của Công ty (Tùy chọn)"
              />
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Trạng thái
              </label>
              <select
                {...register("status")}
                className="w-full p-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
              >
                <option value="active">Đang hoạt động</option>
                <option value="inactive">Tạm ngưng</option>
              </select>
            </div>
          </div>
        </div>

        <div className="bg-white dark:bg-slate-900 rounded-2xl p-6 border border-slate-200 shadow-sm">
          <h3 className="text-lg font-bold text-slate-800 mb-4 border-b pb-2">
            Thông tin Liên hệ
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
            <div className="md:col-span-2">
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Địa chỉ
              </label>
              <div className="relative">
                <MapPin
                  size={18}
                  className="absolute left-3 top-3 text-slate-400"
                />
                <input
                  {...register("address")}
                  className="w-full pl-10 pr-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
                  placeholder="Nhập địa chỉ đầy đủ..."
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Quản lý / Người đại diện
              </label>
              <input
                {...register("manager")}
                className="w-full p-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
                placeholder="Tên quản lý..."
              />
            </div>

            <div>
              <label className="block text-sm font-bold text-slate-700 mb-1">
                Số điện thoại
              </label>
              <div className="relative">
                <Phone
                  size={18}
                  className="absolute left-3 top-3 text-slate-400"
                />
                <input
                  {...register("phone")}
                  className="w-full pl-10 pr-3 py-2.5 bg-slate-50 border border-slate-200 rounded-xl outline-none focus:border-primary"
                  placeholder="0909..."
                />
              </div>
            </div>
          </div>
        </div>
      </form>
    </div>
  );
}
