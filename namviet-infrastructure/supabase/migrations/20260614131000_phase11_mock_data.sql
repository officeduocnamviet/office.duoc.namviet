-- ==============================================================================
-- PHẦN 11: MOCK DATA (DỮ LIỆU MẪU BAN ĐẦU)
-- ==============================================================================

-- 1. Insert Company (Công ty tổng)
INSERT INTO public.companies (id, name, tax_code, address, phone, email, representative_name)
VALUES 
('11111111-1111-1111-1111-111111111111', 'Hệ thống Y tế Nam Việt', '0123456789', 'Hà Nội, Việt Nam', '19001234', 'contact@namviet.com', 'Nguyen Van Giam Doc')
ON CONFLICT (id) DO NOTHING;

-- 2. Insert Roles
INSERT INTO public.roles (id, name, description)
VALUES 
('33333333-3333-3333-3333-333333333331', 'System Admin', 'Toàn quyền'),
('33333333-3333-3333-3333-333333333332', 'Bác Sĩ', 'Khám chữa bệnh'),
('33333333-3333-3333-3333-333333333333', 'Dược Sĩ Bán Hàng', 'Bán hàng')
ON CONFLICT (id) DO NOTHING;

-- 3. Insert Admin User
-- Mật khẩu thật là: namviet123.
-- Do hiện tại Supabase Auth tạo User cần qua Auth API, nên ta chỉ insert vào bảng public.users cho logic ứng dụng.
-- Khi có Supabase auth, Admin cần tự register trên giao diện, nhưng ID này là Mock UUID.
INSERT INTO public.users (id, role_id, company_id, full_name, email, phone, status)
VALUES 
('00000000-0000-0000-0000-000000000001', '33333333-3333-3333-3333-333333333331', '11111111-1111-1111-1111-111111111111', 'Super Admin Nam Việt', 'admin@namviet.com', '0901234567', 'active')
ON CONFLICT (id) DO NOTHING;

-- 5. Insert Categories
INSERT INTO public.categories (id, name, slug)
VALUES 
(1, 'Thuốc kê đơn', 'thuoc-ke-don'),
(2, 'Thực phẩm chức năng', 'thuc-pham-chuc-nang'),
(3, 'Dịch vụ khám bệnh', 'dich-vu-kham-benh')
ON CONFLICT (id) DO NOTHING;

-- 6. Insert Products
INSERT INTO public.products (id, category_id, sku, name, actual_cost, retail_unit, stock_management_type)
VALUES 
(1, 1, 'SP001', 'Panadol Extra Đỏ (Vỉ 10 Viên)', 15000, 'Vỉ', 'lot_date'),
(2, 2, 'SP002', 'Vitamin C 500mg Domesco', 25000, 'Hộp', 'lot_date'),
(3, 3, 'SV001', 'Khám Tổng quát Cơ bản', 200000, 'Lượt', 'none')
ON CONFLICT (id) DO NOTHING;

-- 7. Insert Customers
INSERT INTO public.customers (id, name, phone, loyalty_points)
VALUES 
(1, 'Khách Vãng Lai', '0000000000', 0),
(2, 'Nguyễn Văn A', '0912345678', 500)
ON CONFLICT (id) DO NOTHING;
