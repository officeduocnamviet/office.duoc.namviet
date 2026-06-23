// Package apperr định nghĩa taxonomy lỗi nghiệp vụ TRUNG LẬP domain dùng chung
// cho mọi bounded context. Domain/app chỉ trả các lỗi ở đây; KHÔNG biết HTTP/Huma.
// Việc map apperr → HTTP status + envelope nằm ở tầng http (xem
// internal/platform/httpx/humax). Đây là một phần của shared kernel
// (ARCHITECTURE.md §4) nên phải tuyệt đối không phụ thuộc hạ tầng.
package apperr

import "errors"

// Kind phân loại lỗi nghiệp vụ ở mức trừu tượng (độc lập transport). Tầng http
// dịch mỗi Kind sang đúng HTTP status + mã envelope.
type Kind int

const (
	// KindInternal là lỗi không phân loại / lỗi hệ thống (mặc định → 500).
	KindInternal Kind = iota
	// KindNotFound: tài nguyên không tồn tại (→ 404).
	KindNotFound
	// KindConflict: vi phạm bất biến/đụng độ trạng thái (→ 409).
	KindConflict
	// KindValidation: input không hợp lệ về mặt nghiệp vụ (→ 422).
	KindValidation
	// KindUnauthorized: chưa xác thực / thông tin xác thực sai (→ 401).
	KindUnauthorized
	// KindForbidden: đã xác thực nhưng thiếu quyền (→ 403).
	KindForbidden
)

// String trả tên Kind (phục vụ log/debug, không phải mã envelope).
func (k Kind) String() string {
	switch k {
	case KindNotFound:
		return "not_found"
	case KindConflict:
		return "conflict"
	case KindValidation:
		return "validation"
	case KindUnauthorized:
		return "unauthorized"
	case KindForbidden:
		return "forbidden"
	default:
		return "internal"
	}
}

// Error là lỗi nghiệp vụ có Kind + Code máy đọc được + Message người đọc + lỗi
// gốc (Unwrap). Code cho phép FE phân nhánh ổn định ngay cả khi cùng Kind
// (vd "refresh_reuse_detected" vẫn là KindUnauthorized).
type Error struct {
	Kind    Kind
	Code    string
	Message string
	cause   error
}

// Error thoả interface error.
func (e *Error) Error() string {
	if e.cause != nil {
		return e.Message + ": " + e.cause.Error()
	}
	return e.Message
}

// Unwrap trả lỗi gốc để dùng với errors.Is/As.
func (e *Error) Unwrap() error { return e.cause }

// WithCause gắn lỗi gốc (nguyên nhân) mà giữ nguyên Kind/Code/Message. Trả về
// chính *Error để tiện gọi chuỗi.
func (e *Error) WithCause(cause error) *Error {
	e.cause = cause
	return e
}

// New tạo lỗi với Kind/Code/Message tuỳ ý.
func New(kind Kind, code, message string) *Error {
	return &Error{Kind: kind, Code: code, Message: message}
}

// Các constructor tiện dụng — Code mặc định trùng tên Kind, có thể override bằng
// cách dựng *Error trực tiếp khi cần mã ổn định riêng.

// NotFound tạo lỗi KindNotFound.
func NotFound(message string) *Error {
	return &Error{Kind: KindNotFound, Code: "not_found", Message: message}
}

// Conflict tạo lỗi KindConflict.
func Conflict(message string) *Error {
	return &Error{Kind: KindConflict, Code: "conflict", Message: message}
}

// Validation tạo lỗi KindValidation.
func Validation(message string) *Error {
	return &Error{Kind: KindValidation, Code: "validation_error", Message: message}
}

// Unauthorized tạo lỗi KindUnauthorized. Code mặc định "unauthorized"; truyền
// code riêng để phân biệt nguyên nhân (vd "refresh_reuse_detected").
func Unauthorized(message string) *Error {
	return &Error{Kind: KindUnauthorized, Code: "unauthorized", Message: message}
}

// UnauthorizedCode như Unauthorized nhưng đặt Code máy đọc cụ thể.
func UnauthorizedCode(code, message string) *Error {
	return &Error{Kind: KindUnauthorized, Code: code, Message: message}
}

// Forbidden tạo lỗi KindForbidden.
func Forbidden(message string) *Error {
	return &Error{Kind: KindForbidden, Code: "forbidden", Message: message}
}

// Internal tạo lỗi KindInternal (bọc lỗi hệ thống). Message nên ngắn gọn, không
// rò chi tiết nhạy cảm — chi tiết để trong cause cho log.
func Internal(message string) *Error {
	return &Error{Kind: KindInternal, Code: "internal", Message: message}
}

// AsError trích *Error từ một error bất kỳ (errors.As). Trả (nil,false) nếu
// không phải apperr — tầng http coi đó là lỗi nội bộ.
func AsError(err error) (*Error, bool) {
	var e *Error
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// KindOf trả Kind của err nếu là apperr, ngược lại KindInternal.
func KindOf(err error) Kind {
	if e, ok := AsError(err); ok {
		return e.Kind
	}
	return KindInternal
}
