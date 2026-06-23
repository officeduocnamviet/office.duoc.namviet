package http

import (
	"fmt"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
)

// errBadDate trả lỗi Validation cho tham số ngày sai định dạng (YYYY-MM-DD). Đi
// qua apperr để humax.FromAppErr map thành envelope 422 (validation_error) thống
// nhất với phần còn lại của API.
func errBadDate(field string) error {
	return apperr.Validation(fmt.Sprintf("%s phải đúng định dạng YYYY-MM-DD", field))
}
