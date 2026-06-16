# Database Schema (Lược đồ Cơ sở dữ liệu Nam Việt ERP)

### Table: accounting_journals <== Sổ nhật ký kế toán (Ghi nhận bút toán kép)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **entry_date**: date (Nullable: NO) <== Ngày ghi sổ>
- **doc_type**: text (Nullable: NO) <== Loại chứng từ (VD: Phiếu thu, chi, hóa đơn)>
- **source_ref_id**: text (Nullable: YES) <== ID của chứng từ gốc tham chiếu>
- **description**: text (Nullable: YES) <== Diễn giải bút toán>
- **account_debit**: text (Nullable: NO) <== Tài khoản Nợ>
- **account_credit**: text (Nullable: NO) <== Tài khoản Có>
- **amount**: numeric (Nullable: NO) <== Số tiền>
- **posted_by**: uuid (Nullable: YES) <== ID Kế toán viên ghi sổ>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: agent_workflows <== Quản lý các luồng công việc tự động của AI Agent>
- **intent_code**: text (Nullable: NO) <== Mã ý định (VD: order_status)>
- **description**: text (Nullable: NO) <== Mô tả luồng công việc>
- **required_permission**: text (Nullable: NO) <== Quyền cần thiết để AI gọi action>
- **draft_only**: boolean (Nullable: YES) DEFAULT true <== Chỉ nháp, cần người duyệt>
- **api_endpoint**: text (Nullable: NO) <== Đường dẫn API kích hoạt>
- **is_active**: boolean (Nullable: YES) DEFAULT true <== Trạng thái kích hoạt>

### Table: ai_agent_memories <== Bộ nhớ Vector lưu trữ ngữ cảnh khách hàng cho AI>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **user_id**: uuid (Nullable: YES) <== ID nhân viên (nếu có)>
- **customer_id**: bigint (Nullable: YES) <== ID khách hàng>
- **memory_type**: text (Nullable: NO) <== Loại bộ nhớ (VD: preference, history)>
- **memory_text**: text (Nullable: NO) <== Nội dung bộ nhớ bằng văn bản>
- **embedding**: USER-DEFINED (Nullable: YES) <== Vector nhúng cho Semantic Search>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: appointments <== Ghi Lịch Hẹn với khách hàng sử dụng dịch vụ y tế>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **customer_id**: bigint (Nullable: NO) <== Khách hàng>
- **doctor_id**: uuid (Nullable: YES) <== Bác sĩ thực hiện hoặc đăng ký>
- **room_id**: bigint (Nullable: YES) <== ID của Phòng khám>
- **service_type**: text (Nullable: NO) <== Loại dịch vụ>
- **appointment_time**: timestamp with time zone (Nullable: NO) <== Thời gian hẹn>
- **check_in_time**: timestamp with time zone (Nullable: YES) <== Thời gian khách đến check-in>
- **status**: text (Nullable: NO) DEFAULT 'pending'::text <== Trạng thái (pending, completed, cancelled)>
- **symptoms**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb <== Danh sách triệu chứng ban đầu>
- **note**: text (Nullable: YES) <== Ghi chú thêm>
- **created_by**: uuid (Nullable: YES) <== Người tạo lịch hẹn>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật cuối>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa (Soft delete)>

### Table: approval_requests <== Quản lý Yêu cầu phê duyệt (Luồng duyệt chứng từ)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **requester_id**: uuid (Nullable: NO) <== ID người tạo yêu cầu>
- **module**: text (Nullable: NO) <== Phân hệ cần duyệt (VD: hr, finance)>
- **reference_id**: text (Nullable: YES) <== ID chứng từ liên quan (VD: ID hóa đơn)>
- **payload**: jsonb (Nullable: NO) <== Nội dung/Data cần duyệt (JSON)>
- **status**: text (Nullable: NO) DEFAULT 'pending'::text <== Trạng thái yêu cầu>
- **current_step**: integer (Nullable: YES) DEFAULT 1 <== Bước duyệt hiện tại>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: approval_steps <== Các bước chi tiết trong một Yêu cầu phê duyệt>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **request_id**: uuid (Nullable: NO) <== ID của Yêu cầu phê duyệt gốc>
- **step_order**: integer (Nullable: NO) <== Thứ tự bước duyệt (1, 2, 3...)>
- **approver_id**: uuid (Nullable: YES) <== ID người duyệt cụ thể>
- **approver_role_id**: uuid (Nullable: YES) <== ID nhóm quyền duyệt (nếu không chỉ định đích danh)>
- **status**: text (Nullable: NO) DEFAULT 'pending'::text <== Trạng thái bước này>
- **note**: text (Nullable: YES) <== Ghi chú của người duyệt>
- **processed_at**: timestamp with time zone (Nullable: YES) <== Thời gian đã xử lý>

### Table: assets <== Quản lý tài sản cố định của công ty>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên tài sản>
- **category**: text (Nullable: NO) <== Danh mục tài sản (VD: Máy móc, Điện tử)>
- **purchase_price**: numeric (Nullable: NO) <== Giá mua>
- **purchase_date**: date (Nullable: NO) <== Ngày mua>
- **depreciation_months**: integer (Nullable: NO) <== Số tháng khấu hao>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái tài sản>
- **assigned_to**: uuid (Nullable: YES) <== Bàn giao cho nhân viên nào>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: attendance_logs <== Lịch sử chấm công (Check-in/out) của nhân viên>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **user_id**: uuid (Nullable: NO) <== ID Nhân viên>
- **branch_id**: bigint (Nullable: YES) <== ID Chi nhánh chấm công>
- **check_in_time**: timestamp with time zone (Nullable: NO) DEFAULT now() <== Thời điểm check-in>
- **check_in_ip**: text (Nullable: YES) <== Địa chỉ IP lúc check-in>
- **check_in_lat**: numeric (Nullable: YES) <== Tọa độ Latitude lúc check-in>
- **check_in_lng**: numeric (Nullable: YES) <== Tọa độ Longitude lúc check-in>
- **check_out_time**: timestamp with time zone (Nullable: YES) <== Thời điểm check-out>
- **check_out_ip**: text (Nullable: YES) <== Địa chỉ IP lúc check-out>
- **check_out_lat**: numeric (Nullable: YES) <== Tọa độ Latitude lúc check-out>
- **check_out_lng**: numeric (Nullable: YES) <== Tọa độ Longitude lúc check-out>
- **is_valid**: boolean (Nullable: YES) DEFAULT false <== Hợp lệ (đúng vị trí/IP) hay không>
- **working_hours**: numeric (Nullable: YES) DEFAULT 0 <== Số giờ làm việc tính được>
- **note**: text (Nullable: YES) <== Ghi chú giải trình>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>

### Table: batches <== Quản lý Lô sản xuất và Hạn sử dụng của sản phẩm>
- **id**: bigint (Nullable: NO) <== Khóa chính (Tự tăng)>
- **product_id**: bigint (Nullable: NO) <== Thuộc sản phẩm nào>
- **batch_code**: text (Nullable: NO) <== Mã lô (VD: L12345)>
- **expiry_date**: date (Nullable: NO) <== Hạn sử dụng (Date)>
- **manufacturing_date**: date (Nullable: YES) <== Ngày sản xuất>
- **inbound_price**: numeric (Nullable: YES) DEFAULT 0 <== Giá nhập của lô này>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: categories <== Danh mục phân loại sản phẩm>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên danh mục>
- **slug**: text (Nullable: NO) <== Đường dẫn URL chuẩn hóa>
- **parent_id**: bigint (Nullable: YES) <== Danh mục cha (nếu có)>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: chart_of_accounts <== Hệ thống Tài khoản kế toán (Cây tài khoản)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **account_code**: text (Nullable: NO) <== Số hiệu tài khoản (VD: 111, 112, 131)>
- **name**: text (Nullable: NO) <== Tên tài khoản>
- **parent_id**: uuid (Nullable: YES) <== Tài khoản cấp cha>
- **type**: text (Nullable: NO) <== Loại tài khoản (Tài sản, Nợ, Vốn...)>
- **balance_type**: text (Nullable: NO) <== Số dư bên (DEBIT, CREDIT, BOTH)>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái>
- **allow_posting**: boolean (Nullable: NO) DEFAULT true <== Cho phép hạch toán trực tiếp không>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: chat_messages <== Lưu trữ nội dung chi tiết từng dòng tin nhắn Chat>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **session_id**: uuid (Nullable: NO) <== Thuộc Phiên chat nào (chat_sessions)>
- **role**: text (Nullable: NO) <== Vai trò người gửi (user, bot, system, agent)>
- **content**: text (Nullable: YES) <== Nội dung văn bản của tin nhắn>
- **intent**: text (Nullable: YES) <== Ý định do AI phân tích được (nếu là bot)>
- **entities**: jsonb (Nullable: YES) <== Các thực thể trích xuất được (VD: mã đơn hàng)>
- **llm_meta**: jsonb (Nullable: YES) <== Siêu dữ liệu từ LLM (số token, model...)>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian gửi>

### Table: chat_sessions <== Quản lý Phiên chat (Phòng chat) giữa User và AI/Agent>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **user_id**: uuid (Nullable: NO) <== ID người dùng mở phiên chat>
- **platform**: text (Nullable: NO) DEFAULT 'web'::text <== Nền tảng (web, zalo, facebook)>
- **status**: text (Nullable: NO) DEFAULT 'bot'::text <== Trạng thái (đang chat với bot, hoặc chuyển cho human)>
- **context**: jsonb (Nullable: NO) DEFAULT '{}'::jsonb <== Ngữ cảnh chung của phiên>
- **assigned_sales_id**: uuid (Nullable: YES) <== ID Sales được phân công (nếu bot bí)>
- **started_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian bắt đầu>
- **last_activity_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Lần tương tác cuối>
- **closed_at**: timestamp with time zone (Nullable: YES) <== Thời gian đóng phiên>

### Table: clinical_queues <== Hàng đợi khám bệnh lâm sàng (Phòng khám)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **appointment_id**: uuid (Nullable: YES) <== Từ Lịch hẹn nào chuyển sang (nếu có)>
- **customer_id**: bigint (Nullable: NO) <== Khách hàng/Bệnh nhân>
- **doctor_id**: uuid (Nullable: YES) <== Bác sĩ chỉ định (nếu có)>
- **queue_number**: integer (Nullable: NO) <== Số thứ tự chờ>
- **status**: text (Nullable: NO) DEFAULT 'waiting'::text <== Trạng thái (đang chờ, đang khám)>
- **priority_level**: text (Nullable: NO) DEFAULT 'normal'::text <== Mức độ ưu tiên>
- **checked_in_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian Check-in lấy số>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: companies <== Quản lý Thông tin Công ty mẹ và Các Công ty con>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **tax_code**: text (Nullable: NO) <== Mã số thuế>
- **name**: text (Nullable: NO) <== Tên đầy đủ công ty>
- **short_name**: text (Nullable: YES) <== Tên viết tắt>
- **address**: text (Nullable: NO) <== Địa chỉ trụ sở>
- **phone**: text (Nullable: NO) <== Số điện thoại>
- **email**: text (Nullable: NO) <== Email liên hệ>
- **logo_url**: text (Nullable: YES) <== Link Logo>
- **representative_name**: text (Nullable: NO) <== Người đại diện pháp luật>
- **business_license_url**: ARRAY (Nullable: YES) <== Danh sách link Giấy ĐKKD>
- **mission**: text (Nullable: YES) <== Sứ mệnh>
- **vision**: text (Nullable: YES) <== Tầm nhìn>
- **status**: text (Nullable: YES) DEFAULT 'active'::text <== Trạng thái>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: customer_vaccination_records <== Hồ sơ Tiêm chủng của khách hàng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **customer_id**: bigint (Nullable: NO) <== Bệnh nhân/Khách hàng>
- **appointment_id**: uuid (Nullable: YES) <== Lịch hẹn tiêm>
- **product_id**: bigint (Nullable: NO) <== Loại Vắc xin (ID Sản phẩm)>
- **dose_number**: integer (Nullable: NO) DEFAULT 1 <== Mũi tiêm số mấy>
- **expected_date**: date (Nullable: NO) <== Ngày dự kiến tiêm>
- **actual_date**: date (Nullable: YES) <== Ngày thực tế tiêm>
- **status**: text (Nullable: NO) DEFAULT 'pending'::text <== Trạng thái mũi tiêm>
- **consulted_by**: uuid (Nullable: YES) <== Bác sĩ khám sàng lọc>
- **administered_by**: uuid (Nullable: YES) <== Điều dưỡng thực hiện tiêm>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: customer_vouchers <== Kho mã giảm giá (Voucher) của khách hàng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **promotion_id**: uuid (Nullable: YES) <== Nguồn từ Chương trình KM nào>
- **customer_id**: bigint (Nullable: YES) <== Khách hàng sở hữu>
- **is_used**: boolean (Nullable: YES) DEFAULT false <== Đã sử dụng chưa>
- **used_at**: timestamp with time zone (Nullable: YES) <== Thời gian sử dụng>
- **order_id**: uuid (Nullable: YES) <== Đã dùng cho Đơn hàng nào>

### Table: customers <== Quản lý Danh sách Khách hàng/Bệnh nhân>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **customer_code**: text (Nullable: YES) <== Mã khách hàng>
- **name**: text (Nullable: NO) <== Họ tên/Tên doanh nghiệp>
- **customer_type**: text (Nullable: NO) DEFAULT 'B2C'::text <== Loại khách (B2C, B2B)>
- **phone**: text (Nullable: YES) <== Số điện thoại>
- **email**: text (Nullable: YES) <== Email>
- **address**: text (Nullable: YES) <== Địa chỉ>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái>
- **dob**: date (Nullable: YES) <== Ngày sinh>
- **gender**: text (Nullable: YES) <== Giới tính>
- **cccd**: text (Nullable: YES) <== Số CCCD>
- **loyalty_points**: integer (Nullable: YES) DEFAULT 0 <== Điểm tích lũy>
- **b2b_metadata**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb <== Thông tin riêng cho khách B2B (MST, Công nợ định mức...)>
- **current_debt**: numeric (Nullable: YES) DEFAULT 0 <== Công nợ hiện tại>
- **updated_by**: uuid (Nullable: YES) <== Người cập nhật cuối>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: employment_contracts <== Quản lý Hợp đồng lao động của Nhân sự>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **user_id**: uuid (Nullable: NO) <== ID Nhân viên>
- **contract_code**: text (Nullable: NO) <== Mã hợp đồng>
- **base_salary**: numeric (Nullable: NO) DEFAULT 0 <== Lương cơ bản>
- **standard_working_days**: integer (Nullable: NO) DEFAULT 26 <== Số ngày công chuẩn trong tháng>
- **kpi_conversion_rate**: numeric (Nullable: YES) DEFAULT 0 <== Tỷ lệ quy đổi điểm KPI ra tiền>
- **commission_rate_percent**: numeric (Nullable: YES) DEFAULT 0 <== % Hoa hồng doanh số>
- **tax_deduction_amount**: numeric (Nullable: YES) DEFAULT 0 <== Mức giảm trừ thuế>
- **insurance_deduction_amount**: numeric (Nullable: YES) DEFAULT 0 <== Mức đóng bảo hiểm>
- **valid_from**: date (Nullable: NO) <== Ngày hiệu lực>
- **valid_to**: date (Nullable: YES) <== Ngày hết hiệu lực>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái HĐ>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>

### Table: finance_transactions <== Quản lý Giao dịch Thu/Chi Tài chính>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **code**: text (Nullable: NO) <== Mã giao dịch (Phiếu thu/chi)>
- **transaction_date**: timestamp with time zone (Nullable: NO) DEFAULT now() <== Ngày giao dịch>
- **flow**: text (Nullable: NO) <== Dòng tiền (IN, OUT)>
- **business_type**: text (Nullable: NO) DEFAULT 'other'::text <== Nghiệp vụ (Bán hàng, Trả lương, Nhập hàng...)>
- **category_id**: bigint (Nullable: YES) <== Nhóm thu chi>
- **amount**: numeric (Nullable: NO) <== Số tiền>
- **fund_account_id**: bigint (Nullable: NO) <== Lấy từ Quỹ/Ngân hàng nào>
- **partner_type**: text (Nullable: YES) <== Loại đối tác (Nhân viên, Khách hàng, NCC)>
- **partner_id**: text (Nullable: YES) <== ID Đối tác>
- **partner_name_cache**: text (Nullable: YES) <== Tên đối tác (Lưu cache)>
- **ref_type**: text (Nullable: YES) <== Loại chứng từ gốc>
- **ref_id**: text (Nullable: YES) <== ID Chứng từ gốc (VD: Mã đơn hàng)>
- **description**: text (Nullable: YES) <== Diễn giải>
- **evidence_url**: text (Nullable: YES) <== Ảnh chứng từ đính kèm>
- **status**: text (Nullable: NO) DEFAULT 'pending'::text <== Trạng thái (Hoàn thành, Nháp)>
- **cash_tally**: jsonb (Nullable: YES) <== Kiểm đếm tiền mặt (Bảng kê mệnh giá)>
- **ref_advance_id**: bigint (Nullable: YES) <== Hoàn ứng cho phiếu nào>
- **target_bank_info**: jsonb (Nullable: YES) <== Thông tin ngân hàng đích đến>
- **bank_reference_id**: text (Nullable: YES) <== Mã tham chiếu của Ngân hàng (Webhook)>
- **book_type**: text (Nullable: NO) DEFAULT 'BOTH'::text <== Ghi sổ nội bộ hay sổ thuế>
- **is_posted**: boolean (Nullable: NO) DEFAULT false <== Đã hạch toán vào sổ kế toán chưa>
- **created_by**: uuid (Nullable: YES) <== Kế toán lập phiếu>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: fund_accounts <== Quản lý Sổ Quỹ Tiền mặt và Tài khoản Ngân hàng>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên Quỹ/Tài khoản>
- **type**: text (Nullable: NO) <== Loại (Tiền mặt, Ngân hàng)>
- **location**: text (Nullable: YES) <== Vị trí két sắt (Nếu là quỹ mặt)>
- **account_number**: text (Nullable: YES) <== Số tài khoản>
- **bank_id**: bigint (Nullable: YES) <== Thuộc ngân hàng nào>
- **initial_balance**: numeric (Nullable: NO) DEFAULT 0 <== Số dư đầu kỳ>
- **balance**: numeric (Nullable: NO) DEFAULT 0 <== Số dư hiện tại (Realtime)>
- **currency**: text (Nullable: YES) DEFAULT 'VND'::text <== Đơn vị tiền tệ>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái>
- **bank_info**: jsonb (Nullable: YES) <== Thông tin thêm của Bank>
- **description**: text (Nullable: YES) <== Ghi chú>
- **account_id**: text (Nullable: YES) <== Tài khoản kế toán tương ứng>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: internal_channels <== (CŨ) Kênh chat nội bộ (Đã thay bằng chat_sessions)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **name**: text (Nullable: YES) <== Tên nhóm chat>
- **type**: text (Nullable: NO) <== Loại kênh (1-1, group)>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: internal_messages <== (CŨ) Tin nhắn nội bộ (Đã thay bằng chat_messages)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **channel_id**: uuid (Nullable: NO) <== Thuộc kênh nào>
- **sender_id**: uuid (Nullable: YES) <== Người gửi>
- **content**: text (Nullable: YES) <== Nội dung>
- **attachments**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb <== File đính kèm>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: inventory_batches <== Theo dõi Tồn kho chi tiết theo từng Lô tại từng Kho>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **warehouse_id**: bigint (Nullable: NO) <== Kho hàng>
- **product_id**: bigint (Nullable: NO) <== Sản phẩm>
- **batch_id**: bigint (Nullable: NO) <== Lô hàng>
- **quantity**: numeric (Nullable: NO) DEFAULT 0 <== Số lượng tồn>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Cập nhật cuối>

### Table: inventory_transactions <== Lịch sử Xuất/Nhập/Kiểm kê Tồn kho (Thẻ Kho)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **warehouse_id**: bigint (Nullable: NO) <== Kho thực hiện giao dịch>
- **product_id**: bigint (Nullable: NO) <== Sản phẩm giao dịch>
- **batch_id**: bigint (Nullable: YES) <== Lô sản phẩm giao dịch>
- **type**: text (Nullable: NO) <== Phân loại (IN, OUT)>
- **action_group**: text (Nullable: YES) <== Nhóm hành động (Bán hàng, Hủy hàng, Chuyển kho)>
- **quantity**: numeric (Nullable: NO) <== Số lượng thay đổi (dương hoặc âm)>
- **unit_price**: numeric (Nullable: YES) DEFAULT 0 <== Giá vốn tại thời điểm đó>
- **ref_id**: text (Nullable: YES) <== Chứng từ gốc (Mã hóa đơn)>
- **description**: text (Nullable: YES) <== Diễn giải>
- **partner_id**: bigint (Nullable: YES) <== Khách hàng hoặc NCC>
- **created_by**: uuid (Nullable: YES) <== Thủ kho>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Ngày giao dịch>

### Table: manufacturers <== Quản lý Hãng / Nhà sản xuất>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên hãng>
- **slug**: text (Nullable: NO) <== Đường dẫn URL>
- **country**: text (Nullable: YES) <== Quốc gia>
- **logo_url**: text (Nullable: YES) <== Logo hãng>
- **status**: text (Nullable: YES) DEFAULT 'active'::text <== Trạng thái>
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: NO) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: marketing_campaigns <== Quản lý Chiến dịch Marketing>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên chiến dịch>
- **objective**: text (Nullable: NO) <== Mục tiêu (Tăng nhận diện, Sale...)>
- **start_date**: timestamp with time zone (Nullable: NO) <== Ngày bắt đầu>
- **end_date**: timestamp with time zone (Nullable: NO) <== Ngày kết thúc>
- **budget**: numeric (Nullable: YES) DEFAULT 0 <== Ngân sách dự kiến>
- **actual_cost**: numeric (Nullable: YES) DEFAULT 0 <== Chi phí thực tế>
- **status**: text (Nullable: NO) DEFAULT 'draft'::text <== Trạng thái chiến dịch>
- **target_segment_id**: bigint (Nullable: YES) <== Phân khúc khách hàng mục tiêu>
- **created_by**: uuid (Nullable: YES) <== Người tạo>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: medical_knowledge_vectors <== Cơ sở tri thức Y khoa (Dạng Vector cho AI)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **title**: text (Nullable: NO) <== Tiêu đề tài liệu y khoa>
- **content**: text (Nullable: NO) <== Nội dung chi tiết (Phác đồ, hướng dẫn)>
- **metadata**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb <== Dữ liệu mô tả thêm>
- **embedding**: USER-DEFINED (Nullable: YES) <== Vector nhúng để AI tìm kiếm ngữ nghĩa>
- **created_by**: uuid (Nullable: YES) <== Người nhập liệu>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>

### Table: medical_visits <== Hồ sơ Khám bệnh lâm sàng của khách hàng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **appointment_id**: uuid (Nullable: YES) <== Khám từ Lịch hẹn nào>
- **customer_id**: bigint (Nullable: NO) <== Bệnh nhân>
- **doctor_id**: uuid (Nullable: YES) <== Bác sĩ khám>
- **temperature**: numeric (Nullable: YES) <== Thân nhiệt (Sinh hiệu)>
- **pulse**: integer (Nullable: YES) <== Mạch>
- **sp02**: integer (Nullable: YES) <== Nồng độ Oxy>
- **bp_systolic**: integer (Nullable: YES) <== Huyết áp tâm thu>
- **bp_diastolic**: integer (Nullable: YES) <== Huyết áp tâm trương>
- **weight**: numeric (Nullable: YES) <== Cân nặng>
- **height**: numeric (Nullable: YES) <== Chiều cao>
- **symptoms**: text (Nullable: YES) <== Triệu chứng lâm sàng>
- **examination_summary**: text (Nullable: YES) <== Tóm tắt khám bệnh>
- **diagnosis**: text (Nullable: YES) <== Chẩn đoán xác định>
- **icd_code**: text (Nullable: YES) <== Mã bệnh ICD-10>
- **doctor_notes**: text (Nullable: YES) <== Ghi chú của Bác sĩ (Kê toa/Dặn dò)>
- **red_flags**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb <== Cảnh báo dấu hiệu nguy hiểm>
- **status**: text (Nullable: NO) DEFAULT 'in_progress'::text <== Trạng thái ca khám>
- **created_by**: uuid (Nullable: YES) <== Điều dưỡng tạo>
- **updated_by**: uuid (Nullable: YES) <== Bác sĩ cập nhật>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: order_items <== Chi tiết Sản phẩm trong Đơn hàng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **order_id**: uuid (Nullable: NO) <== Thuộc đơn hàng nào>
- **product_id**: bigint (Nullable: NO) <== Sản phẩm>
- **quantity**: integer (Nullable: NO) <== Số lượng>
- **uom**: text (Nullable: NO) <== Đơn vị tính được chọn (VD: Vỉ)>
- **conversion_factor**: integer (Nullable: YES) <== Tỷ lệ quy đổi so với đơn vị cơ sở>
- **base_quantity**: integer (Nullable: YES) <== Số lượng quy ra đơn vị nhỏ nhất>
- **unit_price**: numeric (Nullable: NO) <== Đơn giá bán>
- **discount**: numeric (Nullable: YES) DEFAULT 0 <== Chiết khấu dòng>
- **is_gift**: boolean (Nullable: YES) DEFAULT false <== Hàng tặng kèm>
- **note**: text (Nullable: YES) <== Ghi chú>
- **batch_no**: text (Nullable: YES) <== Lấy từ lô nào>
- **expiry_date**: date (Nullable: YES) <== Hạn sử dụng của lô đó>
- **total_line**: numeric (Nullable: YES) <== Thành tiền dòng này>
- **quantity_picked**: integer (Nullable: YES) DEFAULT 0 <== Số lượng đã xuất kho thực tế>
- **quantity_returned**: integer (Nullable: YES) DEFAULT 0 <== Số lượng khách trả lại>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: orders <== Quản lý Đơn hàng Tổng (Bán buôn, Bán lẻ, Online)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **code**: text (Nullable: NO) <== Mã hóa đơn (VD: HD123)>
- **customer_id**: bigint (Nullable: YES) <== Khách hàng>
- **creator_id**: uuid (Nullable: YES) <== Người lập đơn (Sales/Thu ngân)>
- **status**: text (Nullable: NO) DEFAULT 'PENDING'::text <== Trạng thái xử lý đơn>
- **order_type**: text (Nullable: NO) DEFAULT 'B2C'::text <== Loại đơn (B2C, B2B)>
- **total_amount**: numeric (Nullable: YES) DEFAULT 0 <== Tổng tiền trước chiết khấu>
- **final_amount**: numeric (Nullable: YES) DEFAULT 0 <== Tổng tiền khách phải trả>
- **payment_status**: text (Nullable: YES) DEFAULT 'unpaid'::text <== Trạng thái thanh toán (đã thu tiền/công nợ)>
- **note**: text (Nullable: YES) <== Ghi chú đơn>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: payroll_items <== (CŨ) Chi tiết Phiếu lương (Đã thay bằng V2)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **payroll_id**: uuid (Nullable: NO)
- **employee_id**: uuid (Nullable: NO)
- **basic_salary**: numeric (Nullable: NO) DEFAULT 0
- **working_days**: numeric (Nullable: NO) DEFAULT 0
- **kpi_score**: numeric (Nullable: YES) DEFAULT 0
- **bonuses**: numeric (Nullable: YES) DEFAULT 0
- **deductions**: numeric (Nullable: YES) DEFAULT 0
- **net_pay**: numeric (Nullable: NO) DEFAULT 0

### Table: payroll_items_v2 <== Chi tiết Phiếu lương cho từng nhân sự (Bản V2)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **payroll_id**: uuid (Nullable: NO) <== Thuộc Kỳ lương nào>
- **employee_id**: uuid (Nullable: NO) <== Nhân viên>
- **base_salary**: numeric (Nullable: NO) <== Lương cơ bản theo hợp đồng>
- **standard_days**: integer (Nullable: NO) <== Số ngày chuẩn>
- **actual_days**: numeric (Nullable: NO) <== Số ngày công thực tế chấm được>
- **prorated_salary**: numeric (Nullable: NO) <== Lương cơ bản theo ngày công>
- **total_kpi_points**: numeric (Nullable: YES) DEFAULT 0 <== Điểm KPI đạt được>
- **kpi_bonus_amount**: numeric (Nullable: YES) DEFAULT 0 <== Tiền thưởng KPI>
- **total_sales_amount**: numeric (Nullable: YES) DEFAULT 0 <== Doanh số đạt được (Nếu là sales)>
- **commission_bonus_amount**: numeric (Nullable: YES) DEFAULT 0 <== Tiền hoa hồng doanh số>
- **other_bonus_amount**: numeric (Nullable: YES) DEFAULT 0 <== Tiền thưởng khác/Phụ cấp>
- **tax_deduction**: numeric (Nullable: YES) DEFAULT 0 <== Trừ Thuế TNCN>
- **insurance_deduction**: numeric (Nullable: YES) DEFAULT 0 <== Trừ Bảo hiểm>
- **net_pay**: numeric (Nullable: NO) <== Thực lĩnh cuối cùng>
- **employee_agreed**: boolean (Nullable: YES) DEFAULT false <== Nhân viên đã xác nhận đúng>
- **employee_note**: text (Nullable: YES) <== Ghi chú/Thắc mắc của nhân viên>
- **accountant_verified**: boolean (Nullable: YES) DEFAULT false <== Kế toán đã rà soát>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>

### Table: payrolls <== Quản lý Kỳ lương / Bảng tổng hợp lương tháng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **code**: text (Nullable: NO) <== Mã kỳ lương>
- **month**: integer (Nullable: NO) <== Tháng tính lương>
- **year**: integer (Nullable: NO) <== Năm tính lương>
- **total_basic_salary**: numeric (Nullable: YES) DEFAULT 0 <== Tổng quỹ lương cơ bản>
- **total_allowance**: numeric (Nullable: YES) DEFAULT 0 <== Tổng quỹ phụ cấp>
- **total_bonus_kpi**: numeric (Nullable: YES) DEFAULT 0 <== Tổng quỹ thưởng KPI>
- **total_deduction**: numeric (Nullable: YES) DEFAULT 0 <== Tổng quỹ phạt/khấu trừ>
- **net_pay**: numeric (Nullable: YES) DEFAULT 0 <== Tổng tiền phải chi cho kỳ này>
- **status**: text (Nullable: NO) DEFAULT 'draft'::text <== Trạng thái (Nháp, Đã duyệt, Đã chi)>
- **created_by**: uuid (Nullable: YES) <== HR tạo bảng lương>
- **approved_by**: uuid (Nullable: YES) <== GĐ duyệt>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>

### Table: product_inventory <== Quản lý Tồn kho tổng hợp tại từng chi nhánh>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **product_id**: bigint (Nullable: YES) <== Sản phẩm>
- **warehouse_id**: bigint (Nullable: YES) <== Kho hàng>
- **stock_quantity**: numeric (Nullable: NO) DEFAULT 0 <== Tổng tồn kho>
- **min_stock**: integer (Nullable: YES) DEFAULT 0 <== Cảnh báo tồn tối thiểu>
- **max_stock**: integer (Nullable: YES) DEFAULT 0 <== Mức tồn tối đa>
- **shelf_location**: text (Nullable: YES) DEFAULT 'Chưa xếp'::text <== Vị trí kệ hàng chung>
- **location_cabinet**: text (Nullable: YES) <== Tủ số mấy>
- **location_row**: text (Nullable: YES) <== Hàng số mấy>
- **location_slot**: text (Nullable: YES) <== Ô số mấy>
- **updated_by**: uuid (Nullable: YES) <== Người cập nhật vị trí>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Lần cập nhật cuối>

### Table: product_units <== Quản lý Các đơn vị tính đa cấp (Hộp, Vỉ, Viên) của thuốc>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **product_id**: bigint (Nullable: YES) <== Thuộc sản phẩm nào>
- **unit_name**: text (Nullable: NO) <== Tên ĐVT (VD: Hộp, Viên)>
- **conversion_rate**: integer (Nullable: YES) DEFAULT 1 <== Tỷ lệ quy đổi so với đơn vị nhỏ nhất>
- **barcode**: text (Nullable: YES) <== Mã vạch riêng của ĐVT này>
- **is_base**: boolean (Nullable: YES) DEFAULT false <== Đánh dấu đây là đơn vị cơ sở (nhỏ nhất)>
- **is_direct_sale**: boolean (Nullable: YES) DEFAULT true <== Cho phép bán lẻ ĐVT này không>
- **price_cost**: numeric (Nullable: YES) DEFAULT 0 <== Giá vốn nhập của ĐVT này>
- **price_sell**: numeric (Nullable: YES) DEFAULT 0 <== Giá bán lẻ định mức của ĐVT này>
- **unit_type**: text (Nullable: YES) DEFAULT 'retail'::text <== Phân loại (Sỉ/Lẻ)>
- **price**: numeric (Nullable: YES) DEFAULT 0 <== Giá áp dụng chung>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: product_vectors <== Vector Nhúng của sản phẩm (Hỗ trợ AI tìm kiếm ngữ nghĩa)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **product_id**: bigint (Nullable: NO) <== Sản phẩm>
- **semantic_text**: text (Nullable: NO) <== Nội dung text kết hợp (Tên + HDSD + Công dụng)>
- **embedding**: USER-DEFINED (Nullable: YES) <== Vector Toán học được tạo từ OpenAI/Gemini>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>

### Table: products <== Danh mục Sản phẩm / Từ điển Thuốc>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên thuốc/Sản phẩm>
- **sku**: text (Nullable: YES) <== Mã định danh nội bộ>
- **barcode**: text (Nullable: YES) <== Mã vạch chuẩn>
- **description**: text (Nullable: YES) <== Mô tả ngắn>
- **active_ingredient**: text (Nullable: YES) <== Hoạt chất chính>
- **image_url**: text (Nullable: YES) <== Ảnh đại diện>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái kinh doanh>
- **fts**: tsvector (Nullable: YES) <== Index tìm kiếm Full-text search>
- **category_id**: bigint (Nullable: YES) <== Nhóm ngành hàng>
- **manufacturer_id**: bigint (Nullable: YES) <== Nhà sản xuất>
- **category_name**: text (Nullable: YES) <== Tên nhóm hàng (Cache)>
- **manufacturer_name**: text (Nullable: YES) <== Tên nhà sản xuất (Cache)>
- **distributor_id**: bigint (Nullable: YES) <== Nhà phân phối>
- **invoice_price**: numeric (Nullable: YES) DEFAULT 0 <== Giá in trên hóa đơn VAT>
- **actual_cost**: numeric (Nullable: NO) DEFAULT 0 <== Giá vốn thực tế>
- **wholesale_unit**: text (Nullable: YES) DEFAULT 'Hộp'::text <== Đơn vị bán sỉ>
- **retail_unit**: text (Nullable: YES) DEFAULT 'Vỉ'::text <== Đơn vị bán lẻ>
- **conversion_factor**: integer (Nullable: YES) DEFAULT 1 <== Tỷ lệ quy đổi Sỉ/Lẻ>
- **wholesale_margin_value**: numeric (Nullable: YES) DEFAULT 0 <== Biên độ lợi nhuận bán buôn>
- **wholesale_margin_type**: text (Nullable: YES) DEFAULT '%'::text <== Loại lợi nhuận (% hay VNĐ)>
- **retail_margin_value**: numeric (Nullable: YES) DEFAULT 0 <== Biên độ lợi nhuận bán lẻ>
- **retail_margin_type**: text (Nullable: YES) DEFAULT '%'::text <== Loại lợi nhuận lẻ (% hay VNĐ)>
- **items_per_carton**: integer (Nullable: YES) DEFAULT 1 <== Quy cách đóng thùng>
- **carton_weight**: numeric (Nullable: YES) DEFAULT 0 <== Khối lượng thùng>
- **carton_dimensions**: text (Nullable: YES) <== Kích thước thùng (Dài x Rộng x Cao)>
- **purchasing_policy**: text (Nullable: YES) DEFAULT 'ALLOW_LOOSE'::text <== Chính sách mua hàng>
- **registration_number**: text (Nullable: YES) <== Số đăng ký lưu hành thuốc>
- **packing_spec**: text (Nullable: YES) <== Quy cách đóng gói>
- **stock_management_type**: text (Nullable: YES) DEFAULT 'lot_date'::text <== Quản lý kho theo Lô/Hạn hay FIFO>
- **wholesale_margin_rate**: numeric (Nullable: YES) DEFAULT 0 <== Tỷ suất LN sỉ (%)>
- **retail_margin_rate**: numeric (Nullable: YES) DEFAULT 0 <== Tỷ suất LN lẻ (%)>
- **usage_instructions**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb <== HDSD chi tiết (JSON)>
- **stock_status**: text (Nullable: YES) DEFAULT 'in_stock'::text <== Trạng thái tồn kho (còn hàng/hết)>
- **product_images**: ARRAY (Nullable: YES) DEFAULT '{}'::text[] <== Danh sách thư viện ảnh>
- **updated_by**: uuid (Nullable: YES) <== Người sửa cuối>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: promotions <== Quản lý Chương trình Khuyến mãi / Giảm giá>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **code**: text (Nullable: NO) <== Mã chương trình (VD: KM_TET)>
- **name**: text (Nullable: NO) <== Tên chương trình>
- **rules**: jsonb (Nullable: NO) <== Quy tắc KM dạng JSON (Điều kiện áp dụng, Giá trị giảm)>
- **start_date**: timestamp with time zone (Nullable: NO) <== Thời gian bắt đầu chạy>
- **end_date**: timestamp with time zone (Nullable: NO) <== Thời gian kết thúc>
- **status**: text (Nullable: YES) DEFAULT 'active'::text <== Trạng thái>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>

### Table: roles <== Phân quyền người dùng (Nhóm quyền)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **name**: text (Nullable: NO) <== Tên nhóm quyền (VD: Admin, Kế toán, Thu ngân)>
- **description**: text (Nullable: YES) <== Mô tả quyền>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **permissions**: jsonb (Nullable: YES) DEFAULT '[]'::jsonb <== Danh sách mã quyền chi tiết (VD: ["users.read", "orders.write"])>

### Table: shift_assignments <== Lịch phân ca làm việc cho nhân viên>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **shift_id**: bigint (Nullable: NO) <== ID của Ca làm việc>
- **user_id**: uuid (Nullable: NO) <== Nhân viên được xếp ca>
- **work_date**: date (Nullable: NO) <== Ngày làm việc>
- **status**: text (Nullable: NO) DEFAULT 'scheduled'::text <== Trạng thái (Đã xếp, Nghỉ phép)>
- **is_overtime**: boolean (Nullable: YES) DEFAULT false <== Đánh dấu là ca tăng ca>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian xếp ca>

### Table: shift_handovers <== Phiếu Bàn giao ca (Chốt sổ Quầy Thu Ngân)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **assignment_id**: bigint (Nullable: YES) <== Bàn giao của ca nào>
- **user_id**: uuid (Nullable: NO) <== Người bàn giao>
- **branch_id**: bigint (Nullable: NO) <== Chi nhánh>
- **system_cash_amount**: numeric (Nullable: NO) DEFAULT 0 <== Tiền mặt ghi nhận trên máy>
- **system_cod_amount**: numeric (Nullable: NO) DEFAULT 0 <== Tiền thu hộ COD trên máy>
- **actual_cash_submitted**: numeric (Nullable: NO) <== Số tiền mặt nộp thực tế (Kiểm đếm tay)>
- **status**: text (Nullable: NO) DEFAULT 'pending_finance'::text <== Trạng thái chờ thủ quỹ xác nhận>
- **finance_transaction_id**: bigint (Nullable: YES) <== Phiếu thu nội bộ tương ứng>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian bàn giao>

### Table: system_audit_logs <== Nhật ký hệ thống (Audit Trail) - Lưu mọi thay đổi dữ liệu>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **table_name**: text (Nullable: NO) <== Bảng dữ liệu bị thay đổi (VD: products)>
- **record_id**: text (Nullable: NO) <== Khóa chính của dòng dữ liệu>
- **action**: text (Nullable: NO) <== Hành động (INSERT, UPDATE, DELETE)>
- **old_data**: jsonb (Nullable: YES) <== Dữ liệu trước khi sửa>
- **new_data**: jsonb (Nullable: YES) <== Dữ liệu sau khi sửa>
- **performed_by**: uuid (Nullable: YES) <== ID người thao tác>
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now() <== Thời điểm sửa>

### Table: system_audit_logs_2026_06 <== Bảng phân mảnh của system_audit_logs theo Tháng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **table_name**: text (Nullable: NO)
- **record_id**: text (Nullable: NO)
- **action**: text (Nullable: NO)
- **old_data**: jsonb (Nullable: YES)
- **new_data**: jsonb (Nullable: YES)
- **performed_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now()

### Table: system_audit_logs_2026_07 <== Bảng phân mảnh của system_audit_logs theo Tháng>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid()
- **table_name**: text (Nullable: NO)
- **record_id**: text (Nullable: NO)
- **action**: text (Nullable: NO)
- **old_data**: jsonb (Nullable: YES)
- **new_data**: jsonb (Nullable: YES)
- **performed_by**: uuid (Nullable: YES)
- **created_at**: timestamp with time zone (Nullable: NO) DEFAULT now()

### Table: system_configs <== Cấu hình hệ thống chung (Tham số dùng chung)>
- **key**: text (Nullable: NO) <== Khóa cấu hình (VD: thuế_VAT_mặc_định)>
- **value**: jsonb (Nullable: NO) <== Giá trị cấu hình>
- **description**: text (Nullable: YES) <== Ghi chú>
- **updated_by**: uuid (Nullable: YES) <== Người đổi cấu hình cuối>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Cập nhật lúc>

### Table: third_party_connections <== Cấu hình kết nối API Đối tác thứ 3>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **provider**: text (Nullable: NO) <== Nhà cung cấp (VD: GHTK, Zalo, ViettelPost)>
- **access_token**: text (Nullable: YES) <== Token truy cập API>
- **refresh_token**: text (Nullable: YES) <== Token làm mới>
- **config**: jsonb (Nullable: YES) DEFAULT '{}'::jsonb <== Các tham số môi trường khác>
- **status**: text (Nullable: YES) DEFAULT 'connected'::text <== Trạng thái kết nối>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Lần làm mới token cuối>

### Table: training_courses <== Quản lý Các khóa Đào tạo Nội bộ>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **title**: text (Nullable: NO) <== Tiêu đề bài giảng>
- **content_type**: text (Nullable: NO) <== Loại nội dung (Video, PDF)>
- **content_url**: text (Nullable: YES) <== Link tài liệu>
- **passing_score**: integer (Nullable: YES) <== Điểm thi đậu yêu cầu>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái bài giảng>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian cập nhật>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian xóa>

### Table: user_fcm_tokens <== Lưu trữ mã định danh Thiết bị để nhận Thông báo Push (Firebase)>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **user_id**: uuid (Nullable: NO) <== Của người dùng nào>
- **token**: text (Nullable: NO) <== Mã token FCM của trình duyệt/điện thoại>
- **device_info**: text (Nullable: YES) <== Thông tin thiết bị (iPhone, Chrome)>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Lần kích hoạt cuối>

### Table: user_social_mappings <== Liên kết tài khoản hệ thống với Mạng xã hội>
- **user_id**: uuid (Nullable: NO) <== ID nhân viên>
- **telegram_id**: text (Nullable: YES) <== ID Telegram (để nhận thông báo chatbot)>
- **zalo_id**: text (Nullable: YES) <== ID Zalo>
- **is_verified**: boolean (Nullable: YES) DEFAULT false <== Đã xác thực hay chưa>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Lần liên kết>

### Table: users <== Quản lý Hồ sơ Nhân sự và Người dùng đăng nhập hệ thống>
- **id**: uuid (Nullable: NO) <== Khóa chính (Trùng với ID của Supabase Auth)>
- **email**: text (Nullable: YES) <== Email đăng nhập>
- **full_name**: text (Nullable: YES) <== Họ tên>
- **avatar_url**: text (Nullable: YES) <== Link ảnh đại diện>
- **employee_code**: text (Nullable: YES) <== Mã nhân viên>
- **position**: text (Nullable: YES) <== Chức vụ (Trưởng phòng, Nhân viên)>
- **status**: text (Nullable: NO) DEFAULT 'pending_approval'::text <== Trạng thái (Đang làm, Đã nghỉ)>
- **dob**: date (Nullable: YES) <== Ngày sinh>
- **phone**: text (Nullable: YES) <== Số điện thoại>
- **gender**: text (Nullable: YES) <== Giới tính>
- **cccd**: text (Nullable: YES) <== CCCD>
- **cccd_issue_date**: date (Nullable: YES) <== Ngày cấp CCCD>
- **address**: text (Nullable: YES) <== Địa chỉ thường trú>
- **marital_status**: text (Nullable: YES) <== Tình trạng hôn nhân>
- **cccd_front_url**: text (Nullable: YES) <== Ảnh mặt trước CCCD>
- **cccd_back_url**: text (Nullable: YES) <== Ảnh mặt sau CCCD>
- **education_level**: text (Nullable: YES) <== Trình độ học vấn>
- **specialization**: text (Nullable: YES) <== Chuyên môn>
- **bank_name**: text (Nullable: YES) <== Ngân hàng nhận lương>
- **bank_account_number**: text (Nullable: YES) <== Số tài khoản nhận lương>
- **bank_account_name**: text (Nullable: YES) <== Tên chủ thẻ>
- **hobbies**: text (Nullable: YES) <== Sở thích>
- **limitations**: text (Nullable: YES) <== Điểm yếu>
- **strengths**: text (Nullable: YES) <== Điểm mạnh>
- **needs**: text (Nullable: YES) <== Nhu cầu đào tạo>
- **work_state**: text (Nullable: YES) DEFAULT 'working'::text <== Trạng thái công việc>
- **role_id**: uuid (Nullable: YES) <== Nhóm quyền truy cập phần mềm>
- **company_id**: uuid (Nullable: YES) <== Trực thuộc pháp nhân công ty nào>
- **warehouse_id**: bigint (Nullable: YES) <== Chi nhánh/Kho trực thuộc>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **updated_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Cập nhật thông tin công việc>
- **profile_updated_at**: timestamp with time zone (Nullable: YES) <== Cập nhật hồ sơ cá nhân>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Nghỉ việc (Soft delete)>

### Table: vendor_product_mappings <== Ánh xạ mã sản phẩm của Nhà cung cấp với Mã Nội bộ>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **vendor_tax_code**: text (Nullable: NO) <== MST Nhà cung cấp>
- **vendor_product_name**: text (Nullable: NO) <== Tên thuốc trên hóa đơn NCC xuất>
- **vendor_unit**: text (Nullable: YES) <== ĐVT của NCC>
- **internal_product_id**: bigint (Nullable: NO) <== Map vào Mã thuốc nội bộ của Nam Việt>
- **internal_unit**: text (Nullable: YES) <== Map vào ĐVT nội bộ>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Ngày ánh xạ>
- **last_used_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Lần dùng đối chiếu cuối>
- **updated_by**: uuid (Nullable: YES) <== Kế toán thực hiện map>

### Table: warehouses <== Quản lý Danh sách Kho hàng / Chi nhánh / Cửa hàng bán lẻ>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **key**: text (Nullable: NO) <== Khóa định danh>
- **name**: text (Nullable: NO) <== Tên kho / Tên chi nhánh>
- **unit**: text (Nullable: NO) DEFAULT 'Hộp'::text <== Đơn vị quản lý tồn kho gốc>
- **address**: text (Nullable: YES) <== Địa chỉ kho>
- **type**: text (Nullable: NO) DEFAULT 'retail'::text <== Phân loại (Sỉ, Lẻ, Kho tổng)>
- **latitude**: numeric (Nullable: YES) <== Tọa độ Latitude>
- **longitude**: numeric (Nullable: YES) <== Tọa độ Longitude>
- **code**: text (Nullable: YES) <== Mã chi nhánh>
- **manager**: text (Nullable: YES) <== Tên Quản lý / Cửa hàng trưởng>
- **phone**: text (Nullable: YES) <== SĐT chi nhánh>
- **status**: text (Nullable: NO) DEFAULT 'active'::text <== Trạng thái>
- **company_id**: uuid (Nullable: YES) <== Trực thuộc công ty con nào>
- **outlet_type**: text (Nullable: YES) <== Mô hình kinh doanh>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
- **deleted_at**: timestamp with time zone (Nullable: YES) <== Thời gian đóng cửa kho>

### Table: webhook_logs <== Lịch sử nhận dữ liệu tự động (Webhook) từ Đối tác>
- **id**: uuid (Nullable: NO) DEFAULT gen_random_uuid() <== Khóa chính>
- **provider**: text (Nullable: NO) <== Đối tác gửi (Ngân hàng, GHTK)>
- **payload**: jsonb (Nullable: NO) <== Dữ liệu thô gửi đến>
- **status**: text (Nullable: NO) DEFAULT 'processing'::text <== Trạng thái xử lý (Lỗi, Thành công)>
- **error_message**: text (Nullable: YES) <== Lỗi nếu phân tích tịt>
- **processed_at**: timestamp with time zone (Nullable: YES) <== Thời gian đã xử lý>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian nhận>

### Table: work_shifts <== Danh mục Ca làm việc trong ngày (Khai báo ca)>
- **id**: bigint (Nullable: NO) <== Khóa chính>
- **branch_id**: bigint (Nullable: NO) <== Áp dụng cho chi nhánh nào>
- **name**: text (Nullable: NO) <== Tên ca (Ca Sáng, Ca Chiều)>
- **start_time**: time without time zone (Nullable: NO) <== Giờ bắt đầu (VD: 08:00)>
- **end_time**: time without time zone (Nullable: NO) <== Giờ kết thúc (VD: 17:00)>
- **is_active**: boolean (Nullable: YES) DEFAULT true <== Trạng thái áp dụng>
- **created_at**: timestamp with time zone (Nullable: YES) DEFAULT now() <== Thời gian tạo>
