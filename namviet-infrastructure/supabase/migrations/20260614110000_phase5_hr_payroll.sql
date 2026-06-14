-- ==========================================
-- 1. HỢP ĐỒNG LAO ĐỘNG (NƠI LƯU TRỮ "LUẬT" CHO TỪNG NGƯỜI)
-- ==========================================
CREATE TABLE public.employment_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    contract_code TEXT UNIQUE NOT NULL,
    
    -- Lương & Ngày công chuẩn
    base_salary NUMERIC NOT NULL DEFAULT 0,
    standard_working_days INTEGER NOT NULL DEFAULT 26, -- Kinh doanh: 24, Kế toán: 26
    
    -- Các tỷ lệ quy đổi (Cá nhân hóa)
    kpi_conversion_rate NUMERIC DEFAULT 0, -- Số tiền / 1 điểm KPI (VD: 50000)
    commission_rate_percent NUMERIC DEFAULT 0, -- % Hoa hồng doanh số (VD: 2.5)
    
    -- Các khoản giảm trừ cố định (Kế toán/HR cài đặt)
    tax_deduction_amount NUMERIC DEFAULT 0, -- Thuế TNCN cố định (nếu có)
    insurance_deduction_amount NUMERIC DEFAULT 0, -- BHXH tự đóng
    
    valid_from DATE NOT NULL,
    valid_to DATE,
    status TEXT NOT NULL CHECK (status IN ('active', 'expired', 'terminated')) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
COMMENT ON TABLE public.employment_contracts IS 'Hợp đồng lao động, nơi cấu hình công thức tính lương động cho từng nhân sự';

-- ==========================================
-- 2. NHẬT KÝ ĐIỂM DANH (CHỐNG GIAN LẬN GPS/IP)
-- ==========================================
CREATE TABLE public.attendance_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.users(id),
    branch_id BIGINT REFERENCES public.warehouses(id), -- Điểm danh tại Cơ sở/Kho nào
    
    -- Dữ liệu Check-in
    check_in_time TIMESTAMPTZ NOT NULL DEFAULT now(),
    check_in_ip TEXT,
    check_in_lat NUMERIC,
    check_in_lng NUMERIC,
    
    -- Dữ liệu Check-out
    check_out_time TIMESTAMPTZ,
    check_out_ip TEXT,
    check_out_lat NUMERIC,
    check_out_lng NUMERIC,
    
    -- Trạng thái hợp lệ (Backend tự động tính toán dựa trên sai số GPS)
    is_valid BOOLEAN DEFAULT false, 
    working_hours NUMERIC DEFAULT 0, -- Số giờ làm thực tế trong ngày
    note TEXT,
    
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
COMMENT ON TABLE public.attendance_logs IS 'Dữ liệu thô điểm danh hàng ngày, khóa chặt vị trí và IP mạng';

-- Tối ưu Index cho Cronjob quét cuối tháng
CREATE INDEX idx_attendance_user_time ON public.attendance_logs(user_id, check_in_time);

-- ==========================================
-- 3. CHI TIẾT BẢNG LƯONG (PAYROLL ITEMS - NÂNG CẤP)
-- ==========================================
-- (Bảng public.payrolls tổng đã tạo ở phần trước, đây là bảng chi tiết cho từng người)
CREATE TABLE public.payroll_items_v2 (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    payroll_id UUID NOT NULL REFERENCES public.payrolls(id) ON DELETE CASCADE,
    employee_id UUID NOT NULL REFERENCES public.users(id),
    
    -- 1. Lương tỷ lệ
    base_salary NUMERIC NOT NULL,
    standard_days INTEGER NOT NULL,
    actual_days NUMERIC NOT NULL,
    prorated_salary NUMERIC NOT NULL, -- (base / standard) * actual
    
    -- 2. Thưởng KPIs & Tasks
    total_kpi_points NUMERIC DEFAULT 0,
    kpi_bonus_amount NUMERIC DEFAULT 0,
    
    -- 3. Thưởng Doanh số (Hoa hồng)
    total_sales_amount NUMERIC DEFAULT 0,
    commission_bonus_amount NUMERIC DEFAULT 0,
    
    -- 4. Thưởng khác (Đào tạo, Lễ tết)
    other_bonus_amount NUMERIC DEFAULT 0,
    
    -- 5. Giảm trừ
    tax_deduction NUMERIC DEFAULT 0,
    insurance_deduction NUMERIC DEFAULT 0,
    
    -- THỰC LÃNH
    net_pay NUMERIC NOT NULL,
    
    -- Luồng Kiểm tra chéo (Cross-check)
    employee_agreed BOOLEAN DEFAULT false,
    employee_note TEXT,
    accountant_verified BOOLEAN DEFAULT false,
    
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    UNIQUE(payroll_id, employee_id)
);
