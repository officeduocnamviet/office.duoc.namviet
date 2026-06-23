// Package id cung cấp helper sinh UUID cho ID nghiệp vụ, dùng chung mọi module
// (shared kernel, ARCHITECTURE.md §4). Ưu tiên UUIDv7 (time-ordered, tốt cho
// index theo thời gian); nếu sinh v7 lỗi vì lý do nào đó thì fallback v4 để
// không bao giờ trả ID rỗng. Trung lập domain, không phụ thuộc hạ tầng.
package id

import "github.com/google/uuid"

// New trả một UUIDv7 mới (time-ordered). Fallback UUIDv4 nếu v7 lỗi.
func New() uuid.UUID {
	v7, err := uuid.NewV7()
	if err != nil {
		return uuid.New() // v4, panic-free trong google/uuid khi có nguồn rand
	}
	return v7
}

// NewString tiện ích trả New() dạng chuỗi chuẩn (8-4-4-4-12).
func NewString() string {
	return New().String()
}

// Parse parse chuỗi thành uuid.UUID (bao bọc để module khác không import trực
// tiếp google/uuid nếu chỉ cần parse ID).
func Parse(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}
