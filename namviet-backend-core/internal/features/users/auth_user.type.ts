/**
 * DTOs cho quá trình xác thực (Authentication)
 */

export interface LoginRequest {
  email: string;      // Bắt buộc. Ví dụ: admin@namviet.com
  password: string;   // Bắt buộc. Ví dụ: namviet123
}

export interface LoginResponse {
  token: string;      // Chuỗi JWT hoặc Master Token
  user: User;         // Thông tin chi tiết của User sau khi đăng nhập
}

/**
 * Interface biểu diễn bảng Users trong Database.
 * File Backend tương ứng: internal/features/users/models.go
 */
export interface User {
  id: string;         // UUID (Dùng cho Supabase Auth)
  email: string;
  full_name: string;
  phone: string;
  status: string;     // 'pending_approval', 'active', 'inactive'
  role_id: string;    // UUID của Role
  company_id: string; // UUID của Company
  created_at: string; // ISO 8601 Date string
  updated_at: string; // ISO 8601 Date string
}
