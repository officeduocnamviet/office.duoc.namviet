package humax

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// New dựng một huma.API trên chi router với contract envelope {data,error}:
//   - response thành công đi qua envelopeTransformer → {data:<body>, error:null};
//   - response lỗi dùng codeError (qua huma.NewError) → {data:null, error:{...}}.
//
// Đồng thời TẮT SchemaLinkTransformer mặc định của Huma (chèn $schema vào body)
// để envelope sạch đúng như RPC cũ. OpenAPI vẫn sinh 3.1.
//
// Lưu ý: huma.NewError là biến package-level toàn cục của Huma. Ta gán nó ở đây
// để mọi tiện ích huma.Error4xx/5xx + lỗi validation tự sinh đều ra envelope.
func New(r chi.Router, title, version string) huma.API {
	huma.NewError = newError

	cfg := huma.DefaultConfig(title, version)
	// Loại CreateHooks (chứa SchemaLinkTransformer) và đặt transformer envelope.
	cfg.CreateHooks = nil
	cfg.Transformers = []huma.Transformer{envelopeTransformer}
	applySecurity(&cfg)

	return humachi.New(r, cfg)
}

// Config trả về huma.Config thuần (không gắn router) — hữu ích khi cần dựng API
// trên một mux khác hoặc để introspect OpenAPI trong test/dump.
func Config(title, version string) huma.Config {
	huma.NewError = newError
	cfg := huma.DefaultConfig(title, version)
	cfg.CreateHooks = nil
	cfg.Transformers = []huma.Transformer{envelopeTransformer}
	applySecurity(&cfg)
	return cfg
}

// applySecurity khai báo security scheme "bearerAuth" (JWT access token) trong
// OpenAPI components để các operation cần đăng nhập tham chiếu được. Token là
// JWT ES256 phát bởi module identity.
func applySecurity(cfg *huma.Config) {
	if cfg.Components == nil {
		cfg.Components = &huma.Components{}
	}
	if cfg.Components.SecuritySchemes == nil {
		cfg.Components.SecuritySchemes = map[string]*huma.SecurityScheme{}
	}
	cfg.Components.SecuritySchemes["bearerAuth"] = &huma.SecurityScheme{
		Type:         "http",
		Scheme:       "bearer",
		BearerFormat: "JWT",
		Description:  "Access JWT (ES256) — header: Authorization: Bearer <token>",
	}
}
