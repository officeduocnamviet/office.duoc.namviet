package http

import "github.com/Maneva-AI/namviet-backend/internal/common/apperr"

// errUnauthenticated là lỗi phòng vệ khi handler /me chạy mà thiếu claims (lẽ
// ra middleware đã chặn). Trả 401.
func errUnauthenticated() error {
	return apperr.Unauthorized("chưa xác thực")
}
