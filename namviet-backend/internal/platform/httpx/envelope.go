// Package httpx sở hữu hợp đồng response của backend: envelope {data,error}
// tương thích với safeRpc của FE (ERP), để FE giữ nguyên call site khi cutover.
package httpx

import (
	"encoding/json"
	"net/http"
)

// APIError là phần "error" của envelope.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type envelope struct {
	Data  any       `json:"data"`
	Error *APIError `json:"error"`
}

// WriteData ghi response thành công: {"data": <data>, "error": null}.
func WriteData(w http.ResponseWriter, status int, data any) {
	write(w, status, envelope{Data: data, Error: nil})
}

// WriteError ghi response lỗi: {"data": null, "error": {code, message}}.
func WriteError(w http.ResponseWriter, status int, code, msg string) {
	write(w, status, envelope{Data: nil, Error: &APIError{Code: code, Message: msg}})
}

func write(w http.ResponseWriter, status int, body envelope) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
