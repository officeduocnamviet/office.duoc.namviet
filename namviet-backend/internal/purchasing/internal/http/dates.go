package http

import "time"

// parseDatePtr parse chuỗi "YYYY-MM-DD" → *time.Time (rỗng → nil, lỗi → error để
// handler trả 422). Dùng cho expiry_date / manufacturing_date của dòng PO.
func parseDatePtr(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
