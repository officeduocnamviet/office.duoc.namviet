// Package pagination cung cấp keyset/cursor pagination dùng chung (ARCHITECTURE.md
// §4, §5: "Pagination keyset"). Cursor là một con trỏ OPAQUE (base64) tới khoá
// sắp xếp của bản ghi cuối trang trước — ở đây là id int64 tăng dần. Opaque để
// FE không phụ thuộc cấu trúc nội bộ và ta đổi khoá sắp xếp sau mà không vỡ API.
// Trung lập domain, chỉ stdlib.
package pagination

import (
	"encoding/base64"
	"strconv"
)

// EncodeID mã hoá một id int64 thành cursor opaque (base64 URL-safe). Dùng cho
// "next cursor" = id của bản ghi cuối trang hiện tại.
func EncodeID(id int64) string {
	return base64.RawURLEncoding.EncodeToString([]byte(strconv.FormatInt(id, 10)))
}

// DecodeID giải mã cursor về id int64. Cursor rỗng = trang đầu → trả 0 (id
// bigint luôn > 0 nên "id > 0" lấy từ đầu). Cursor sai định dạng → lỗi để handler
// trả 400/422, KHÔNG im lặng coi như trang đầu (tránh nuốt lỗi client).
func DecodeID(cursor string) (int64, error) {
	if cursor == "" {
		return 0, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(string(raw), 10, 64)
	if err != nil {
		return 0, err
	}
	return id, nil
}
