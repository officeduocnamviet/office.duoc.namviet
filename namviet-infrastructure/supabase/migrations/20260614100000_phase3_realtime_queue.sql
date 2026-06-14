-- ==========================================
-- 1. BẢNG HÀNG ĐỢI LÂM SÀNG (CLINICAL QUEUES)
-- ==========================================

-- Dữ liệu giao dịch (thay đổi liên tục) nên dùng UUID làm Khóa chính
CREATE TABLE public.clinical_queues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID REFERENCES public.appointments(id) ON DELETE CASCADE,
    customer_id BIGINT NOT NULL REFERENCES public.customers(id),
    doctor_id UUID REFERENCES public.users(id), -- Có thể NULL nếu hàng đợi chung cho nhiều bác sĩ
    
    queue_number INTEGER NOT NULL, -- Số thứ tự trong ngày
    status TEXT NOT NULL CHECK (status IN ('waiting', 'examining', 'completed', 'skipped', 'waiting_vaccination', 'waiting_procedure', 'observing')) DEFAULT 'waiting',
    priority_level TEXT NOT NULL CHECK (priority_level IN ('normal', 'high')) DEFAULT 'normal', -- High cho cấp cứu/trẻ nhỏ
    
    checked_in_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    deleted_at TIMESTAMPTZ
);
COMMENT ON TABLE public.clinical_queues IS 'Quản lý hàng đợi thực tế tại phòng khám/tiêm chủng';
COMMENT ON COLUMN public.clinical_queues.queue_number IS 'Số thứ tự khám, reset mỗi ngày';

-- Đánh Index để tối ưu truy vấn Realtime theo bác sĩ và ngày hiện tại
CREATE INDEX idx_clinical_queues_doctor_status ON public.clinical_queues(doctor_id, status);
CREATE INDEX idx_clinical_queues_date ON public.clinical_queues(checked_in_at);

-- ==========================================
-- 2. KÍCH HOẠT SUPABASE REALTIME
-- ==========================================
-- Đây là "Phép màu". Mở khóa luồng stream từ Database bắn thẳng xuống Web/App
ALTER PUBLICATION supabase_realtime ADD TABLE public.clinical_queues;
ALTER PUBLICATION supabase_realtime ADD TABLE public.appointments;
