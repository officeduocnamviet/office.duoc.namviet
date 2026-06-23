package humax

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
)

// statusForKind map apperr.Kind → HTTP status. Đây là MỘT điểm dịch duy nhất từ
// lỗi nghiệp vụ (domain-neutral) sang transport — domain/app không bao giờ biết
// HTTP. Status sau đó được CodeForStatus map tiếp sang mã envelope ở §5.
func statusForKind(k apperr.Kind) int {
	switch k {
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	case apperr.KindValidation:
		return http.StatusUnprocessableEntity
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}

// FromAppErr chuyển một error bất kỳ thành huma.StatusError envelope:
//   - nếu là *apperr.Error: dùng đúng status theo Kind và GIỮ Code máy đọc của
//     domain (vd "refresh_reuse_detected") thay vì code mặc định theo status;
//   - nếu không: coi là lỗi nội bộ → 500 "internal" (không rò chi tiết).
//
// Handler ở các module gọi hàm này khi service trả lỗi, nhờ đó toàn bộ lỗi
// nghiệp vụ ra đúng envelope mà KHÔNG để Huma/HTTP rò vào domain.
func FromAppErr(err error) huma.StatusError {
	if e, ok := apperr.AsError(err); ok {
		status := statusForKind(e.Kind)
		code := e.Code
		if code == "" {
			code = CodeForStatus(status)
		}
		return &codeError{
			status:  status,
			Data:    nil,
			ErrBody: errBody{Code: code, Message: e.Message},
		}
	}
	return &codeError{
		status:  http.StatusInternalServerError,
		Data:    nil,
		ErrBody: errBody{Code: "internal", Message: "internal server error"},
	}
}
