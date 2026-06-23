package app

import (
	"fmt"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
)

// errInvalidCredentials là lỗi đăng nhập sai — cố ý KHÔNG phân biệt
// "email không tồn tại" với "mật khẩu sai" để chống dò tài khoản.
func errInvalidCredentials() *apperr.Error {
	return apperr.Unauthorized("email hoặc mật khẩu không đúng")
}

// errRefreshInvalid là lỗi refresh token không hợp lệ (không thấy/hết hạn/đã
// thu hồi). Code ổn định "refresh_invalid".
func errRefreshInvalid() *apperr.Error {
	return apperr.UnauthorizedCode("refresh_invalid", "refresh token không hợp lệ")
}

// errRefreshReuse là lỗi PHÁT HIỆN DÙNG LẠI refresh token (đã đánh dấu used).
// Toàn bộ family bị thu hồi trước khi trả lỗi này. Code "refresh_reuse_detected".
func errRefreshReuse() *apperr.Error {
	return apperr.UnauthorizedCode("refresh_reuse_detected", "phát hiện dùng lại refresh token")
}

// errUnsupportedAlgo cho hash với thuật toán ngoài tập hỗ trợ.
func errUnsupportedAlgo(algo domain.HashAlgo) *apperr.Error {
	return apperr.Internal(fmt.Sprintf("thuật toán băm không hỗ trợ: %q", algo))
}
