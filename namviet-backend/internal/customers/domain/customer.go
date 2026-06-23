// Package domain là LÕI THUẦN của bounded context customers: entity khách hàng,
// hồ sơ B2B, ảnh chụp công nợ + PORT interface. KHÔNG import pgx/http/huma/
// framework (ARCHITECTURE.md §3). Chỉ stdlib + shared kernel trung lập
// (common/money). Phụ thuộc đi một chiều: adapters → app → domain.
//
// Customers ở DB này (nhánh monorepo office.duoc.namviet) là context "nhẹ" (đa số
// read): master khách + hồ sơ công nợ/hạn mức. Thông tin B2B (MST, hạn mức nợ,
// kỳ hạn thanh toán) nằm trong cột b2b_metadata jsonb của bảng public.customers
// — KHÔNG có bảng customers_b2b riêng ở DB này (adapter parse jsonb → B2BProfile).
package domain

import (
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/money"
)

// CustomerType phân loại khách: B2C (lẻ) hay B2B (doanh nghiệp). Giữ kiểu chuỗi
// khớp giá trị public.customers.customer_type ('B2C'|'B2B').
type CustomerType string

const (
	// TypeB2C là khách lẻ.
	TypeB2C CustomerType = "B2C"
	// TypeB2B là khách doanh nghiệp (có hồ sơ B2B + công nợ).
	TypeB2B CustomerType = "B2B"
)

// Customer là aggregate gốc của context customers. ID là int64 vì bảng
// public.customers dùng khoá bigint (KHÔNG uuid). B2B nil khi khách B2C (hoặc
// B2B chưa khai metadata). Debt là ảnh chụp công nợ đã chọn nguồn (xem
// DebtSnapshot) — luôn ưu tiên nguồn LIVE.
type Customer struct {
	ID        int64
	Code      string // customer_code; có thể rỗng (cột nullable)
	Name      string
	Type      CustomerType
	Phone     string
	Email     string
	Address   string
	Status    string // 'active' | 'inactive' | 'banned'
	B2B       *B2BProfile
	Debt      DebtSnapshot
	CreatedAt time.Time
	UpdatedAt time.Time
}

// IsB2B trả true nếu khách thuộc loại doanh nghiệp.
func (c Customer) IsB2B() bool { return c.Type == TypeB2B }

// B2BProfile là value object gom thông tin doanh nghiệp đọc từ b2b_metadata
// jsonb. Tiền (DebtLimit) dùng money.Money — KHÔNG float. Các trường rỗng/0 khi
// metadata không khai. DebtLimit chỉ ĐỌC/hiển thị: credit-limit đang OFF chủ đích
// (chỉ ~17/483 B2B có debt_limit>0) nên KHÔNG enforce chặn ở context này.
type B2BProfile struct {
	TaxCode      string      // MST
	DebtLimit    money.Money // hạn mức công nợ (chỉ hiển thị, không enforce)
	PaymentTerm  int         // kỳ hạn thanh toán (số ngày); 0 = không khai
	SalesStaffID string      // nhân viên phụ trách (uuid dạng chuỗi); rỗng = chưa gán
}

// DebtSource cho biết con số công nợ đến từ đâu, để FE/người đọc hiểu độ tin cậy.
type DebtSource string

const (
	// DebtSourceLive: tính LIVE từ đơn hàng chưa tất toán (ưu tiên).
	DebtSourceLive DebtSource = "live"
	// DebtSourceStatic: lấy từ cột tĩnh customers.current_debt (stale, dự phòng).
	DebtSourceStatic DebtSource = "static"
)

// DebtSnapshot là ảnh chụp công nợ của một khách tại thời điểm đọc. Giữ CẢ hai
// con số (live tính từ orders + static từ cột) để minh bạch, và Amount = con số
// đã chọn theo SelectDebt (ưu tiên live). CAVEAT ghost-debt: Live KHÔNG trừ phần
// đã trả của đơn 'partial' (orders thiếu cột paid_amount ở DB này) → có thể phình
// nợ; KHÔNG tự "sửa" số.
type DebtSnapshot struct {
	Amount money.Money // con số dùng để hiển thị/tính (đã chọn nguồn)
	Live   money.Money // tổng final_amount đơn chưa tất toán
	Static money.Money // customers.current_debt (cột tĩnh)
	Source DebtSource  // nguồn của Amount
}

// SelectDebt chọn nguồn công nợ theo quy tắc nghiệp vụ: ưu tiên LIVE (tính từ
// đơn hàng) làm con số chính. Hàm THUẦN (dễ unit test). hasLiveData=false nghĩa
// là không có dữ liệu đơn để tính live (vd join trả NULL trước COALESCE) → rơi
// về cột tĩnh. Khi có dữ liệu live (kể cả 0 đồng vì khách không nợ) vẫn dùng
// live: 0 live là sự thật "không nợ", đáng tin hơn cột tĩnh stale.
func SelectDebt(live, static money.Money, hasLiveData bool) DebtSnapshot {
	if hasLiveData {
		return DebtSnapshot{Amount: live, Live: live, Static: static, Source: DebtSourceLive}
	}
	return DebtSnapshot{Amount: static, Live: live, Static: static, Source: DebtSourceStatic}
}
