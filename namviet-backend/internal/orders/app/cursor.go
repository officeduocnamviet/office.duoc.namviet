package app

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
)

// Cursor của orders là OPAQUE (base64) bọc cặp khoá sắp xếp của bản ghi cuối
// trang trước: (created_at unix-nano, id uuid). orders.id là UUID nên KHÔNG
// keyset bằng id tăng-dần như bảng bigint được (pagination.EncodeID dùng cho
// bảng bigint). Ta keyset theo (created_at DESC, id DESC) để ổn định kể cả khi
// nhiều đơn cùng mốc thời gian. Định dạng nội bộ "nano|uuid" được giấu sau
// base64 để FE không phụ thuộc cấu trúc (đổi được sau mà không vỡ API).

// encodeCursor mã hoá (createdAtNano, id) thành cursor opaque.
func encodeCursor(createdAtNano int64, id string) string {
	raw := strconv.FormatInt(createdAtNano, 10) + "|" + id
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// decodeCursor giải mã cursor về (createdAtNano, id). Cursor rỗng = trang đầu →
// (0, ""). Sai định dạng → lỗi để handler trả 422 (KHÔNG im lặng coi trang đầu).
func decodeCursor(cursor string) (int64, string, error) {
	if cursor == "" {
		return 0, "", nil
	}
	rawBytes, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return 0, "", err
	}
	raw := string(rawBytes)
	sep := strings.IndexByte(raw, '|')
	if sep < 0 {
		return 0, "", errors.New("cursor thiếu dấu phân tách")
	}
	nano, err := strconv.ParseInt(raw[:sep], 10, 64)
	if err != nil {
		return 0, "", err
	}
	id := raw[sep+1:]
	if id == "" {
		return 0, "", errors.New("cursor thiếu id")
	}
	return nano, id, nil
}
