package server

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// healthStatus là body của /v1/health (sẽ được envelope bọc thành
// {"data":{"status":"ok"},"error":null}).
type healthStatus struct {
	Status string `json:"status" example:"ok" doc:"Trạng thái dịch vụ"`
}

type healthOutput struct {
	Body healthStatus
}

// echoInput là input demo có field bắt buộc, dùng để minh hoạ validation 422.
type echoInput struct {
	Body struct {
		Name string `json:"name" minLength:"1" doc:"Tên (bắt buộc, không rỗng)"`
	}
}

type echoOutput struct {
	Body struct {
		Echo string `json:"echo" doc:"Phản hồi lại name đã gửi"`
	}
}

// registerRoutes đăng ký toàn bộ operation của edge lên huma.API: health + demo
// của platform, rồi route của từng bounded context qua Deps.Modules.
func registerRoutes(api huma.API, d Deps) {
	registerHealth(api, d)
	registerDemo(api)
	for _, mod := range d.Modules {
		if mod != nil {
			mod(api)
		}
	}
}

// registerHealth migrate /v1/health sang một Huma operation. Ping DB nếu có
// pinger; fail → 503 (envelope code "internal").
func registerHealth(api huma.API, d Deps) {
	p := d.pinger()
	huma.Register(api, huma.Operation{
		OperationID: "health",
		Method:      http.MethodGet,
		Path:        "/v1/health",
		Summary:     "Liveness/readiness check",
		Tags:        []string{"system"},
	}, func(ctx context.Context, _ *struct{}) (*healthOutput, error) {
		if p != nil {
			if err := p.Ping(ctx); err != nil {
				return nil, huma.Error503ServiceUnavailable("database unreachable")
			}
		}
		return &healthOutput{Body: healthStatus{Status: "ok"}}, nil
	})
}

// registerDemo đăng ký 2 operation tối thiểu để chứng minh contract envelope
// end-to-end (lỗi 404 + validation 422) và làm OpenAPI có input/output mẫu.
// Đặt dưới /v1/_demo (dấu gạch dưới) để tách khỏi domain route thật.
func registerDemo(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "demo-not-found",
		Method:      http.MethodGet,
		Path:        "/v1/_demo/not-found",
		Summary:     "Demo lỗi 404 envelope",
		Tags:        []string{"_demo"},
	}, func(_ context.Context, _ *struct{}) (*struct{}, error) {
		return nil, huma.Error404NotFound("không tìm thấy tài nguyên demo")
	})

	huma.Register(api, huma.Operation{
		OperationID:   "demo-echo",
		Method:        http.MethodPost,
		Path:          "/v1/_demo/echo",
		Summary:       "Demo validation 422 + success envelope",
		Tags:          []string{"_demo"},
		DefaultStatus: http.StatusOK,
	}, func(_ context.Context, in *echoInput) (*echoOutput, error) {
		out := &echoOutput{}
		out.Body.Echo = in.Body.Name
		return out, nil
	})
}
