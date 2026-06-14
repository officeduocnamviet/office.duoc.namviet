-- ==============================================================================
-- PHẦN 13: LƯU TRỮ FCM TOKENS CHO PUSH NOTIFICATION
-- ==============================================================================

CREATE TABLE IF NOT EXISTS public.user_fcm_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.users(id) ON DELETE CASCADE,
    token TEXT NOT NULL UNIQUE,
    device_info TEXT, -- Thông tin máy (VD: iPhone 15, Chrome Windows)
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Tạo Index để tối ưu việc tìm kiếm token của một user
CREATE INDEX IF NOT EXISTS idx_user_fcm_tokens_user_id ON public.user_fcm_tokens(user_id);

COMMENT ON TABLE public.user_fcm_tokens IS 'Bảng lưu trữ Firebase Cloud Messaging (FCM) Tokens của người dùng để gửi thông báo đẩy.';
COMMENT ON COLUMN public.user_fcm_tokens.user_id IS 'Mã người dùng sở hữu thiết bị.';
COMMENT ON COLUMN public.user_fcm_tokens.token IS 'Chuỗi Token duy nhất do Firebase cấp cho thiết bị/trình duyệt.';
COMMENT ON COLUMN public.user_fcm_tokens.device_info IS 'Thông tin nhận dạng thiết bị (tùy chọn).';
