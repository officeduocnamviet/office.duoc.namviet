"use client";

import React, { useState, useEffect } from "react";
import {
  Building2,
  Plus,
  Search,
  Edit2,
  Trash2,
  ChevronLeft,
  Upload,
  Save,
  Image as ImageIcon,
  MapPin,
  Phone,
  Mail,
  Loader2,
  Eye,
  X,
} from "lucide-react";
import { ColumnDef } from "@tanstack/react-table";
import { DataTable } from "../ui/DataTable";

// Interface
interface Company {
  id: string;
  tax_code: string;
  name: string;
  short_name?: string;
  address: string;
  phone: string;
  email?: string;
  logo_url?: string;
  representative_name?: string;
  business_license_url?: string[];
  mission?: string;
  vision?: string;
  status: string;
}

import { apiClient } from "@/lib/axios";

export default function CompanyInfoView() {
  const [viewMode, setViewMode] = useState<"list" | "form">("list");
  const [editingId, setEditingId] = useState<string | null>(null);
  const [companies, setCompanies] = useState<Company[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  const fetchCompanies = async () => {
    setIsLoading(true);
    try {
      const res = await apiClient.get('/v1/companies');
      if (res.data) {
        setCompanies(res.data || []);
      }
    } catch (error) {
      console.error("Lỗi khi tải danh sách công ty:", error);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    if (viewMode === "list") {
      fetchCompanies();
    }
  }, [viewMode]);

  const handleCreate = () => {
    setEditingId(null);
    setViewMode("form");
  };

  const handleEdit = (id: string) => {
    setEditingId(id);
    setViewMode("form");
  };

  const handleDelete = async (id: string) => {
    if (!confirm("Bạn có chắc chắn muốn xóa công ty này?")) return;
    try {
      const res = await apiClient.delete(`/v1/companies/${id}`);
      if (res.status >= 200 && res.status < 300) {
        fetchCompanies();
      } else {
        alert("Lỗi khi xóa");
      }
    } catch (e) {
      console.error(e);
    }
  };

  return (
    <main className="flex-1 overflow-y-auto p-4 lg:p-6 custom-scrollbar bg-transparent h-full">
      {viewMode === "list" ? (
        <CompanyList
          companies={companies}
          isLoading={isLoading}
          onCreate={handleCreate}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      ) : (
        <CompanyForm id={editingId} onBack={() => setViewMode("list")} />
      )}
    </main>
  );
}

function CompanyList({
  companies,
  isLoading,
  onCreate,
  onEdit,
  onDelete,
}: {
  companies: Company[];
  isLoading: boolean;
  onCreate: () => void;
  onEdit: (id: string) => void;
  onDelete: (id: string) => void;
}) {
  const [globalFilter, setGlobalFilter] = useState("");

  const columns: ColumnDef<Company>[] = [
    {
      accessorKey: "logo_url",
      header: "Logo",
      cell: ({ row }) => {
        const logo = row.original.logo_url;
        return (
          <div className="w-10 h-10 rounded-lg bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 flex items-center justify-center overflow-hidden shrink-0">
            {logo ? (
              <img
                src={logo}
                alt="Logo"
                className="w-full h-full object-contain p-1"
              />
            ) : (
              <Building2 size={16} className="text-slate-400" />
            )}
          </div>
        );
      },
    },
    {
      accessorKey: "name",
      header: "Tên Công ty",
      cell: ({ row }) => (
        <div>
          <div
            className="font-bold text-sm text-primary-600 dark:text-blue-400 hover:underline cursor-pointer"
            onClick={() => onEdit(row.original.id)}
          >
            {row.original.name}
          </div>
          <div className="text-xs text-slate-500 font-medium mt-1">
            MST: {row.original.tax_code}
          </div>
        </div>
      ),
    },
    {
      accessorKey: "phone",
      header: "Liên hệ",
      cell: ({ row }) => (
        <div className="text-sm">
          <div className="flex items-center gap-1.5 text-slate-600 dark:text-slate-400 mb-1">
            <Phone size={12} /> {row.original.phone}
          </div>
          {row.original.email && (
            <div className="flex items-center gap-1.5 text-slate-600 dark:text-slate-400">
              <Mail size={12} /> {row.original.email}
            </div>
          )}
        </div>
      ),
    },
    {
      accessorKey: "representative_name",
      header: "Người đại diện",
      cell: ({ row }) => (
        <span className="text-sm font-medium text-slate-700 dark:text-slate-300">
          {row.original.representative_name || "-"}
        </span>
      ),
    },
    {
      accessorKey: "status",
      header: "Trạng thái",
      cell: ({ row }) => (
        <span
          className={`inline-flex items-center px-2 py-1 rounded-md text-xs font-bold border ${row.original.status === "active" ? "bg-emerald-50 text-emerald-600 border-emerald-100" : "bg-slate-50 text-slate-600 border-slate-200"}`}
        >
          {row.original.status === "active" ? "Hoạt động" : "Tạm ngưng"}
        </span>
      ),
    },
    {
      id: "actions",
      header: "Thao tác",
      cell: ({ row }) => (
        <div className="flex items-center gap-2">
          <button
            onClick={() => onEdit(row.original.id)}
            className="p-1.5 text-slate-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
            title="Chỉnh sửa"
          >
            <Edit2 size={16} />
          </button>
          <button
            onClick={() => onDelete(row.original.id)}
            className="p-1.5 text-slate-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
            title="Xóa"
          >
            <Trash2 size={16} />
          </button>
        </div>
      ),
    },
  ];

  return (
    <div className="space-y-6 animate-in fade-in slide-in-from-bottom-4 duration-300 h-full flex flex-col">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 shrink-0">
        <div>
          <h2 className="text-2xl font-black text-slate-800 dark:text-white tracking-tight flex items-center gap-2">
            <Building2 className="text-primary" /> Thông tin Công ty
          </h2>
          <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">
            Quản lý thông tin công ty mẹ và các công ty thành viên.
          </p>
        </div>
        <button
          onClick={onCreate}
          className="flex items-center gap-2 bg-primary hover:bg-primary-600 text-white px-4 py-2 rounded-xl font-bold transition-all shadow-md shadow-primary/20 active:scale-95"
        >
          <Plus size={18} />
          <span>Thêm Công ty</span>
        </button>
      </div>

      {/* Toolbar */}
      <div className="bg-white dark:bg-slate-900 p-4 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm flex gap-4 shrink-0">
        <div className="relative flex-1 max-w-md">
          <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <Search size={16} className="text-slate-400" />
          </div>
          <input
            type="text"
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder="Tìm kiếm công ty..."
            className="w-full pl-10 pr-4 py-2 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl text-sm outline-none focus:border-primary transition-colors text-slate-700 dark:text-slate-200"
          />
        </div>
      </div>

      {/* Table */}
      <div className="flex-1 overflow-hidden">
        {isLoading ? (
          <div className="h-full flex items-center justify-center text-slate-400">
            <Loader2 className="animate-spin" size={32} />
          </div>
        ) : (
          <DataTable
            columns={columns}
            data={companies}
            globalFilter={globalFilter}
            setGlobalFilter={setGlobalFilter}
          />
        )}
      </div>
    </div>
  );
}

import { uploadToCloudinary } from "@namviet/shared-utils";

function CompanyForm({
  id,
  onBack,
}: {
  id: string | null;
  onBack: () => void;
}) {
  const isEditing = !!id;
  const [formData, setFormData] = useState<Partial<Company>>({
    status: "active",
  });
  const [isSaving, setIsSaving] = useState(false);
  const [isUploading, setIsUploading] = useState(false);
  const [isUploadingLicense, setIsUploadingLicense] = useState(false);
  const [viewImageUrl, setViewImageUrl] = useState<string | null>(null);

  useEffect(() => {
    if (isEditing) {
      // Fetch details
      apiClient.get('/v1/companies')
        .then((res) => {
          const comp = res.data?.find((c: Company) => c.id === id);
          if (comp) setFormData(comp);
        });
    }
  }, [id, isEditing]);

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >,
  ) => {
    setFormData((prev) => ({ ...prev, [e.target.name]: e.target.value }));
  };

  const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setIsUploading(true);
      try {
        const file = e.target.files[0];
        const url = await uploadToCloudinary(file);
        setFormData((prev) => ({ ...prev, logo_url: url }));
      } catch (error) {
        console.error("Lỗi upload ảnh:", error);
        alert("Có lỗi xảy ra khi upload ảnh.");
      } finally {
        setIsUploading(false);
      }
    }
  };

  const handleLicenseUpload = async (
    e: React.ChangeEvent<HTMLInputElement>,
  ) => {
    if (e.target.files && e.target.files.length > 0) {
      setIsUploadingLicense(true);
      try {
        const newUrls = [...(formData.business_license_url || [])];
        for (let i = 0; i < e.target.files.length; i++) {
          const file = e.target.files[i];
          const url = await uploadToCloudinary(file);
          newUrls.push(url);
        }
        setFormData((prev) => ({ ...prev, business_license_url: newUrls }));
      } catch (error) {
        console.error("Lỗi upload giấy phép:", error);
        alert("Có lỗi xảy ra khi upload giấy phép.");
      } finally {
        setIsUploadingLicense(false);
      }
    }
  };

  const removeLicense = (index: number) => {
    const newUrls = [...(formData.business_license_url || [])];
    newUrls.splice(index, 1);
    setFormData((prev) => ({ ...prev, business_license_url: newUrls }));
  };

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const url = isEditing ? `/v1/companies/${id}` : '/v1/companies';
      
      const res = isEditing 
        ? await apiClient.put(url, formData)
        : await apiClient.post(url, formData);

      if (res.status >= 200 && res.status < 300) {
        alert("Lưu thành công!");
        // Không tự động onBack() để user tiếp tục sửa hoặc quay lại thủ công
      } else {
        alert("Lỗi khi lưu");
      }
    } catch (e) {
      console.error(e);
      alert("Lỗi kết nối");
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <div className="max-w-5xl mx-auto pb-10">
      {/* Form Header (Sticky) */}
      <div className="sticky top-0 z-50 pt-1 pb-4">
        <div className="flex items-center justify-between bg-white/20 dark:bg-slate-900/20 backdrop-blur-2xl p-4 rounded-2xl border border-white/40 dark:border-slate-700/40 shadow-sm transition-all duration-300">
          <div className="flex items-center gap-4">
            <button 
              onClick={onBack}
              className="w-10 h-10 flex items-center justify-center rounded-full bg-slate-50/80 dark:bg-slate-800/80 hover:bg-white dark:hover:bg-slate-700 text-slate-600 dark:text-slate-300 transition-colors backdrop-blur-md shadow-sm"
            >
              <ChevronLeft size={20} />
            </button>
            <div>
              <h2 className="text-xl font-black text-slate-800 dark:text-white">
                {isEditing ? 'Cập nhật Công ty' : 'Thêm mới Công ty'}
              </h2>
              <p className="text-xs font-medium text-slate-600 dark:text-slate-400 mt-0.5">
                Vui lòng điền đầy đủ thông tin pháp lý
              </p>
            </div>
          </div>
          <button 
            onClick={handleSave}
            disabled={isSaving || isUploading || isUploadingLicense}
            className="flex items-center gap-2 bg-primary hover:bg-primary-600 disabled:opacity-50 text-white px-5 py-2 rounded-xl font-bold transition-all shadow-md shadow-primary/20 active:scale-95"
          >
            {isSaving ? <Loader2 size={18} className="animate-spin" /> : <Save size={18} />}
            <span>Lưu thông tin</span>
          </button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6 animate-in fade-in slide-in-from-bottom-8 duration-500 mt-2">
        {/* Left Column - Main Info */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm p-6">
            <h3 className="text-base font-bold text-slate-800 dark:text-white mb-4 border-b border-slate-100 dark:border-slate-800 pb-2">
              Thông tin Cơ bản
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div className="md:col-span-2">
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Tên Công ty đầy đủ <span className="text-red-500">*</span>
                </label>
                <input
                  name="name"
                  value={formData.name || ""}
                  onChange={handleChange}
                  type="text"
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Tên viết tắt
                </label>
                <input
                  name="short_name"
                  value={formData.short_name || ""}
                  onChange={handleChange}
                  type="text"
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Mã số thuế <span className="text-red-500">*</span>
                </label>
                <input
                  name="tax_code"
                  value={formData.tax_code || ""}
                  onChange={handleChange}
                  type="text"
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Người đại diện pháp luật{" "}
                  <span className="text-red-500">*</span>
                </label>
                <input
                  name="representative_name"
                  value={formData.representative_name || ""}
                  onChange={handleChange}
                  type="text"
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                />
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Trạng thái
                </label>
                <select
                  name="status"
                  value={formData.status || "active"}
                  onChange={handleChange}
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                >
                  <option value="active">Hoạt động</option>
                  <option value="inactive">Tạm ngưng</option>
                </select>
              </div>
            </div>
          </div>

          <div className="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm p-6">
            <h3 className="text-base font-bold text-slate-800 dark:text-white mb-4 border-b border-slate-100 dark:border-slate-800 pb-2">
              Thông tin Liên hệ
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-5">
              <div className="md:col-span-2">
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Địa chỉ trụ sở <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <MapPin
                    size={18}
                    className="absolute left-3 top-3 text-slate-400"
                  />
                  <textarea
                    name="address"
                    value={formData.address || ""}
                    onChange={handleChange}
                    rows={2}
                    className="w-full pl-10 pr-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white resize-none"
                  ></textarea>
                </div>
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Số điện thoại <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <Phone
                    size={18}
                    className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
                  />
                  <input
                    name="phone"
                    value={formData.phone || ""}
                    onChange={handleChange}
                    type="text"
                    className="w-full pl-10 pr-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                  />
                </div>
              </div>

              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Email <span className="text-red-500">*</span>
                </label>
                <div className="relative">
                  <Mail
                    size={18}
                    className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-400"
                  />
                  <input
                    name="email"
                    value={formData.email || ""}
                    onChange={handleChange}
                    type="email"
                    className="w-full pl-10 pr-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white"
                  />
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm p-6">
            <h3 className="text-base font-bold text-slate-800 dark:text-white mb-4 border-b border-slate-100 dark:border-slate-800 pb-2">
              Tầm nhìn & Sứ mệnh
            </h3>
            <div className="grid grid-cols-1 gap-5">
              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Sứ mệnh
                </label>
                <textarea
                  name="mission"
                  value={formData.mission || ""}
                  onChange={handleChange}
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white resize-y min-h-[100px]"
                ></textarea>
              </div>
              <div>
                <label className="block text-sm font-bold text-slate-700 dark:text-slate-300 mb-1">
                  Tầm nhìn
                </label>
                <textarea
                  name="vision"
                  value={formData.vision || ""}
                  onChange={handleChange}
                  className="w-full px-4 py-2.5 bg-slate-50 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-xl outline-none focus:border-primary text-slate-800 dark:text-white resize-y min-h-[100px]"
                ></textarea>
              </div>
            </div>
          </div>
        </div>

        {/* Right Column - Images/Documents */}
        <div className="space-y-6">
          <div className="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm p-6">
            <h3 className="text-base font-bold text-slate-800 dark:text-white mb-4 border-b border-slate-100 dark:border-slate-800 pb-2">
              Logo Công ty
            </h3>
            <label className="flex flex-col items-center justify-center p-6 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-xl bg-slate-50 dark:bg-slate-800/50 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer group">
              <input
                type="file"
                accept="image/*"
                className="hidden"
                onChange={handleFileChange}
                disabled={isUploading}
              />

              {isUploading ? (
                <div className="flex flex-col items-center py-8">
                  <Loader2
                    size={32}
                    className="animate-spin text-primary mb-2"
                  />
                  <p className="text-sm font-medium text-slate-600">
                    Đang tải ảnh lên...
                  </p>
                </div>
              ) : formData.logo_url ? (
                <div className="relative group/logo">
                  <img
                    src={formData.logo_url}
                    alt="Logo"
                    className="w-32 h-32 object-contain mb-4"
                  />
                  <div className="absolute inset-0 bg-slate-900/60 rounded-xl opacity-0 group-hover/logo:opacity-100 transition-opacity flex items-center justify-center gap-2 mb-4">
                    <button
                      type="button"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        setViewImageUrl(formData.logo_url!);
                      }}
                      className="p-2 bg-white/20 hover:bg-white/40 text-white rounded-lg transition-colors backdrop-blur-sm"
                      title="Xem ảnh"
                    >
                      <Eye size={18} />
                    </button>
                    <button
                      type="button"
                      onClick={(e) => {
                        e.preventDefault();
                        e.stopPropagation();
                        setFormData((prev) => ({ ...prev, logo_url: "" }));
                      }}
                      className="p-2 bg-red-500/80 hover:bg-red-500 text-white rounded-lg transition-colors backdrop-blur-sm"
                      title="Xóa ảnh"
                    >
                      <Trash2 size={18} />
                    </button>
                  </div>
                </div>
              ) : (
                <div className="w-16 h-16 rounded-full bg-slate-200 dark:bg-slate-700 flex items-center justify-center mb-4 group-hover:scale-110 transition-transform">
                  <ImageIcon size={24} className="text-slate-400" />
                </div>
              )}
              {!isUploading && (
                <>
                  <p className="text-sm font-bold text-primary mb-1">
                    Nhấn để tải ảnh lên
                  </p>
                  <p className="text-xs text-slate-500 text-center">
                    Tải lên sẽ tự động gọi optimize_cloudinary_url()
                  </p>
                </>
              )}
            </label>
          </div>

          <div className="bg-white dark:bg-slate-900 rounded-2xl border border-slate-200 dark:border-slate-800 shadow-sm p-6">
            <h3 className="text-base font-bold text-slate-800 dark:text-white mb-4 border-b border-slate-100 dark:border-slate-800 pb-2">
              Giấy phép Kinh doanh
            </h3>
            <div className="space-y-4">
              <label className="flex flex-col items-center justify-center p-4 border-2 border-dashed border-slate-200 dark:border-slate-700 rounded-xl bg-slate-50 dark:bg-slate-800/50 hover:bg-slate-100 dark:hover:bg-slate-800 transition-colors cursor-pointer group">
                <input
                  type="file"
                  accept="image/*"
                  multiple
                  className="hidden"
                  onChange={handleLicenseUpload}
                  disabled={isUploadingLicense}
                />
                {isUploadingLicense ? (
                  <div className="flex flex-col items-center">
                    <Loader2
                      size={24}
                      className="animate-spin text-primary mb-2"
                    />
                    <p className="text-sm font-medium text-slate-600">
                      Đang tải lên...
                    </p>
                  </div>
                ) : (
                  <>
                    <Upload
                      size={24}
                      className="text-slate-400 group-hover:scale-110 transition-transform mb-2"
                    />
                    <p className="text-sm font-bold text-primary mb-1">
                      Thêm ảnh giấy phép
                    </p>
                    <p className="text-xs text-slate-500">
                      Hỗ trợ chọn nhiều ảnh
                    </p>
                  </>
                )}
              </label>

              {formData.business_license_url &&
                formData.business_license_url.length > 0 && (
                  <div className="grid grid-cols-2 gap-3">
                    {formData.business_license_url.map((url, idx) => (
                      <div
                        key={idx}
                        className="relative group/img rounded-lg overflow-hidden border border-slate-200 dark:border-slate-700 bg-slate-100 flex items-center justify-center"
                      >
                        <img
                          src={url}
                          alt={`Giấy phép ${idx + 1}`}
                          className="w-full h-28 object-cover"
                        />
                        <div className="absolute inset-0 bg-slate-900/60 opacity-0 group-hover/img:opacity-100 transition-opacity flex items-center justify-center gap-2">
                          <button
                            type="button"
                            onClick={() => setViewImageUrl(url)}
                            className="p-2 bg-white/20 hover:bg-white/40 text-white rounded-lg transition-colors backdrop-blur-sm"
                            title="Xem ảnh"
                          >
                            <Eye size={18} />
                          </button>
                          <button
                            type="button"
                            onClick={() => removeLicense(idx)}
                            className="p-2 bg-red-500/80 hover:bg-red-500 text-white rounded-lg transition-colors backdrop-blur-sm"
                            title="Xóa"
                          >
                            <Trash2 size={18} />
                          </button>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
            </div>
          </div>
        </div>
      </div>

      {/* View Image Modal */}
      {viewImageUrl && (
        <div
          className="fixed inset-0 z-[100] flex items-center justify-center bg-slate-900/90 backdrop-blur-sm p-4"
          onClick={() => setViewImageUrl(null)}
        >
          <div className="relative max-w-4xl max-h-[90vh] w-full h-full flex flex-col items-center justify-center">
            <button
              className="absolute top-0 right-0 p-3 bg-white/10 hover:bg-white/20 text-white rounded-full transition-colors backdrop-blur-md z-10"
              onClick={() => setViewImageUrl(null)}
            >
              <X size={24} />
            </button>
            <img
              src={viewImageUrl}
              alt="Full view"
              className="max-w-full max-h-[85vh] object-contain rounded-lg shadow-2xl"
              onClick={(e) => e.stopPropagation()}
            />
          </div>
        </div>
      )}
    </div>
  );
}
