package humax

import "github.com/danielgtaylor/huma/v2"

// successEnvelope là vỏ {data, error:null} cho response THÀNH CÔNG.
type successEnvelope struct {
	Data  any `json:"data"`
	Error any `json:"error"` // luôn null ở nhánh thành công
}

// envelopeTransformer bọc body thành công thành {data, error:null}. Huma chạy
// transformer cho MỌI response (kể cả lỗi), nên ta phải passthrough những body
// đã là envelope để tránh double-wrap:
//   - *codeError: đã là {data:null, error:{...}} → giữ nguyên.
//   - successEnvelope: đã bọc rồi (vd lồng group) → giữ nguyên.
//   - nil: không có body (204/304) → giữ nguyên.
func envelopeTransformer(_ huma.Context, _ string, v any) (any, error) {
	switch v.(type) {
	case nil:
		return v, nil
	case *codeError, *successEnvelope, successEnvelope:
		return v, nil
	default:
		return &successEnvelope{Data: v, Error: nil}, nil
	}
}
