// Package humax gắn Huma v2 lên chi router nhưng GIỮ contract envelope
// {data,error} mà FE (safeRpc) phụ thuộc — cho cả success lẫn error. Toàn bộ
// cơ chế envelope/transformer/error gói gọn ở đây, không rò vào domain.
package humax

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// codeError là model lỗi của backend. Nó implement huma.StatusError nên Huma
// dùng đúng status code; đồng thời tự marshal thành envelope
// {"data":null,"error":{"code","message"}} qua tag JSON dưới đây. Vì bản thân
// nó đã là envelope hoàn chỉnh, transformer thành công sẽ BỎ QUA nó (xem
// transform.go) → không double-wrap.
type codeError struct {
	status  int
	Data    *struct{} `json:"data"`  // luôn null khi lỗi
	ErrBody errBody   `json:"error"` // {code, message}
}

type errBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error thoả interface error.
func (e *codeError) Error() string { return e.ErrBody.Message }

// GetStatus thoả huma.StatusError → Huma set đúng HTTP status.
func (e *codeError) GetStatus() int { return e.status }

// ContentType ép luôn application/json (đè mặc định problem+json của Huma) để
// FE nhận đúng content-type như RPC cũ.
func (e *codeError) ContentType(string) string { return "application/json" }

// CodeForStatus map HTTP status → mã lỗi chuỗi theo spec mục 7.5.
func CodeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusUnprocessableEntity:
		return "validation_error"
	default:
		// 500 và mọi 5xx khác gộp về "internal"; các status hiếm khác cũng vậy.
		return "internal"
	}
}

// newError là hàm thay cho huma.NewError. Mọi tiện ích huma.Error4xx/5xx... gọi
// qua đây nên toàn bộ lỗi (kể cả validation tự sinh) đều ra envelope đúng.
func newError(status int, message string, _ ...error) huma.StatusError {
	return &codeError{
		status:  status,
		Data:    nil,
		ErrBody: errBody{Code: CodeForStatus(status), Message: message},
	}
}
