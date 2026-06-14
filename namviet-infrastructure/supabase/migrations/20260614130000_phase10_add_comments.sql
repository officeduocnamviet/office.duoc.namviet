-- ==============================================================================
-- PHẦN 10: TÀI LIỆU HÓA CƠ SỞ DỮ LIỆU (DATABASE COMMENTS)
-- ==============================================================================

-- Bảng users
COMMENT ON TABLE public.users IS 'Bảng lưu trữ thông tin nhân viên trong công ty.';
COMMENT ON COLUMN public.users.id IS 'ID định danh nhân viên (Trùng với UUID của Supabase Auth).';
COMMENT ON COLUMN public.users.role_id IS 'Chức danh của nhân viên (Trưởng phòng, Kế toán...).';
COMMENT ON COLUMN public.users.company_id IS 'Công ty mà nhân viên này đang làm việc.';
COMMENT ON COLUMN public.users.full_name IS 'Họ và tên đầy đủ của nhân viên.';
COMMENT ON COLUMN public.users.phone IS 'Số điện thoại liên hệ.';
COMMENT ON COLUMN public.users.email IS 'Email làm việc, dùng để đăng nhập.';
COMMENT ON COLUMN public.users.status IS 'Trạng thái tài khoản (active/inactive).';

-- Bảng roles
COMMENT ON TABLE public.roles IS 'Bảng danh mục các vai trò, chức vụ trong hệ thống.';
COMMENT ON COLUMN public.roles.name IS 'Tên vai trò (VD: Quản trị viên, Bác sĩ, Dược sĩ).';

-- Bảng products
COMMENT ON TABLE public.products IS 'Bảng danh mục sản phẩm, thuốc, dịch vụ.';
COMMENT ON COLUMN public.products.category_id IS 'ID danh mục nhóm sản phẩm.';
COMMENT ON COLUMN public.products.sku IS 'Mã vạch hoặc mã sản phẩm duy nhất (Stock Keeping Unit).';
COMMENT ON COLUMN public.products.name IS 'Tên hiển thị của sản phẩm/dịch vụ.';
COMMENT ON COLUMN public.products.description IS 'Mô tả chi tiết về sản phẩm.';
COMMENT ON COLUMN public.products.actual_cost IS 'Giá vốn hoặc giá bán lẻ cơ bản của sản phẩm.';
COMMENT ON COLUMN public.products.retail_unit IS 'Đơn vị tính bán lẻ (Viên, Vỉ, Hộp).';

-- Bảng customers
COMMENT ON TABLE public.customers IS 'Bảng lưu trữ thông tin Khách hàng và Bệnh nhân.';
COMMENT ON COLUMN public.customers.name IS 'Họ tên khách hàng.';
COMMENT ON COLUMN public.customers.phone IS 'Số điện thoại liên lạc chính.';
COMMENT ON COLUMN public.customers.loyalty_points IS 'Số điểm tích lũy hiện tại của khách hàng.';

-- Bảng orders
COMMENT ON TABLE public.orders IS 'Bảng lưu trữ Đơn hàng / Hóa đơn bán hàng.';
COMMENT ON COLUMN public.orders.code IS 'Mã đơn hàng hiển thị cho khách xem.';
COMMENT ON COLUMN public.orders.customer_id IS 'Người mua (nếu có). Trống nghĩa là khách vãng lai.';
COMMENT ON COLUMN public.orders.creator_id IS 'Nhân viên lập đơn hàng.';
COMMENT ON COLUMN public.orders.total_amount IS 'Tổng tiền hàng trước khi giảm giá.';
COMMENT ON COLUMN public.orders.final_amount IS 'Tổng tiền thực tế khách phải thanh toán (subtotal - discount).';
COMMENT ON COLUMN public.orders.payment_status IS 'Trạng thái thanh toán (unpaid, partial, paid).';
COMMENT ON COLUMN public.orders.status IS 'Trạng thái giao hàng (PENDING, COMPLETED...).';

-- Bảng order_items
COMMENT ON TABLE public.order_items IS 'Chi tiết các sản phẩm trong Đơn hàng.';
COMMENT ON COLUMN public.order_items.order_id IS 'Thuộc đơn hàng nào.';
COMMENT ON COLUMN public.order_items.product_id IS 'Sản phẩm nào được mua.';
COMMENT ON COLUMN public.order_items.quantity IS 'Số lượng mua.';
COMMENT ON COLUMN public.order_items.unit_price IS 'Đơn giá lúc mua (có thể khác giá gốc của sản phẩm).';
COMMENT ON COLUMN public.order_items.total_line IS 'Thành tiền (quantity * unit_price).';
