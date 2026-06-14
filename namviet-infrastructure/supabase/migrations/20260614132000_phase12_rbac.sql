-- ==============================================================================
-- PHẦN 12: BỔ SUNG RBAC (ROLE-BASED ACCESS CONTROL)
-- ==============================================================================

-- Thêm cột permissions vào bảng roles để phân quyền chi tiết
ALTER TABLE public.roles ADD COLUMN IF NOT EXISTS permissions JSONB DEFAULT '[]'::jsonb;

COMMENT ON COLUMN public.roles.permissions IS 'Danh sách quyền hạn dưới dạng JSONB (VD: ["users.read", "users.write", "orders.create"]).';

-- Cập nhật quyền cho Mock Roles (từ phase 11)
UPDATE public.roles 
SET permissions = '["*"]'::jsonb 
WHERE id = '33333333-3333-3333-3333-333333333331'; -- System Admin

UPDATE public.roles 
SET permissions = '["patients.read", "patients.write", "orders.create"]'::jsonb 
WHERE id = '33333333-3333-3333-3333-333333333332'; -- Bác Sĩ

UPDATE public.roles 
SET permissions = '["products.read", "inventory.read", "orders.create", "orders.read"]'::jsonb 
WHERE id = '33333333-3333-3333-3333-333333333333'; -- Dược Sĩ Bán Hàng
