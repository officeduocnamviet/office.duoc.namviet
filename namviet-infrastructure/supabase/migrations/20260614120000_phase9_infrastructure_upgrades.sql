-- ==============================================================================
-- PHẦN 9: HẠ TẦNG NÂNG CAO (AI VECTORS & PARTITIONING)
-- ==============================================================================

-- 1. Kích hoạt Extension Vector
CREATE EXTENSION IF NOT EXISTS vector;

-- 2. BẢNG VECTOR: CƠ SỞ TRI THỨC Y KHOA (Medical Knowledge)
-- Phục vụ AI tư vấn bệnh, kiểm tra tương tác thuốc (tránh Hallucination)
CREATE TABLE public.medical_knowledge_vectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    metadata JSONB DEFAULT '{}'::jsonb, -- Chứa nguồn gốc, ngày cập nhật, mã ICD...
    embedding vector(768), -- Kích thước vector 768 cho text-embedding-004 (Google Gemini) hoặc text-embedding-ada-002
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
COMMENT ON TABLE public.medical_knowledge_vectors IS 'Cơ sở tri thức (Phác đồ, dược thư) cho AI truy vấn bằng kỹ thuật RAG';

-- Tối ưu hóa tìm kiếm Vector bằng chỉ mục HNSW
CREATE INDEX idx_medical_knowledge_embedding ON public.medical_knowledge_vectors USING hnsw (embedding vector_cosine_ops);

-- 3. BẢNG VECTOR: SẢN PHẨM (Product Embeddings)
-- Phục vụ AI tìm kiếm sản phẩm theo semantic search, tư vấn sản phẩm thay thế
CREATE TABLE public.product_vectors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id BIGINT NOT NULL REFERENCES public.products(id) ON DELETE CASCADE,
    semantic_text TEXT NOT NULL, -- Nội dung ghép từ Tên, Hoạt chất, Công dụng để nhúng
    embedding vector(768),
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
COMMENT ON TABLE public.product_vectors IS 'Vector nhúng của Sản phẩm để AI tìm kiếm theo ngữ nghĩa (Semantic Search) và Gợi ý Mua kèm';
CREATE INDEX idx_product_vectors_embedding ON public.product_vectors USING hnsw (embedding vector_cosine_ops);

-- 4. BẢNG MEMORY: BỘ NHỚ HỘI THOẠI CỦA TRỢ LÝ AI (Agent Memory)
-- Phục vụ AI nhớ ngữ cảnh dài hạn của người dùng, thói quen mua hàng
CREATE TABLE public.ai_agent_memories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES public.users(id),
    customer_id BIGINT REFERENCES public.customers(id),
    memory_type TEXT NOT NULL CHECK (memory_type IN ('preference', 'past_interaction', 'health_condition')),
    memory_text TEXT NOT NULL,
    embedding vector(768),
    created_at TIMESTAMPTZ DEFAULT now()
);
COMMENT ON TABLE public.ai_agent_memories IS 'Bộ nhớ dài hạn (Long-term Memory) của AI Agent cho từng Khách hàng và Nhân viên';
CREATE INDEX idx_ai_agent_memories_embedding ON public.ai_agent_memories USING hnsw (embedding vector_cosine_ops);

-- ==============================================================================
-- PARTITIONING CHO DỮ LIỆU LỚN (TABLE PARTITIONING STRATEGY)
-- Lưu ý: Postgres không cho đổi bảng thường thành Partitioned trực tiếp.
-- Do đây là giai đoạn Init, ta sẽ giả lập cấu trúc Log cho Audit_Logs (ví dụ mẫu).
-- ==============================================================================
CREATE TABLE public.system_audit_logs (
    id UUID DEFAULT gen_random_uuid(),
    table_name TEXT NOT NULL,
    record_id TEXT NOT NULL,
    action TEXT NOT NULL,
    old_data JSONB,
    new_data JSONB,
    performed_by UUID,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);
COMMENT ON TABLE public.system_audit_logs IS 'Bảng log hệ thống được chia Partitions theo tháng để tối ưu truy vấn dữ liệu khổng lồ';

-- Tạo sẵn partition cho 2 tháng hiện tại và tiếp theo
CREATE TABLE public.system_audit_logs_2026_06 PARTITION OF public.system_audit_logs
    FOR VALUES FROM ('2026-06-01') TO ('2026-07-01');

CREATE TABLE public.system_audit_logs_2026_07 PARTITION OF public.system_audit_logs
    FOR VALUES FROM ('2026-07-01') TO ('2026-08-01');
