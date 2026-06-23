package identity_test

import (
	"context"
	"sync"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/identity/app"
)

// TestIntegration_RefreshRace_ConcurrentRotation chứng minh đã vá lỗ hổng
// refresh-token rotation race: HAI request Refresh ĐỒNG THỜI cùng một raw token
// (READ COMMITTED, cùng thấy used=false) chỉ được phép cho ĐÚNG MỘT token con
// hợp lệ — request còn lại phải thua race và bị coi như reuse
// (refresh_reuse_detected), kéo theo cả family bị revoke. Trước khi vá, cả hai
// đều rotate → cấp 2 refresh token hợp lệ, bỏ qua reuse-detection.
func TestIntegration_RefreshRace_ConcurrentRotation(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()

	login, err := fx.svc.Login(ctx, "staff@nv.vn", fx.pw)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}

	const n = 2
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []app.Tokens
		errs    []error
	)
	start := make(chan struct{})
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start // đồng loạt xuất phát để cực đại hoá khả năng race
			tok, rerr := fx.svc.Refresh(ctx, login.RefreshToken)
			mu.Lock()
			if rerr != nil {
				errs = append(errs, rerr)
			} else {
				results = append(results, tok)
			}
			mu.Unlock()
		}()
	}
	close(start)
	wg.Wait()

	// ĐÚNG 1 thành công, 1 lỗi.
	if len(results) != 1 || len(errs) != 1 {
		t.Fatalf("race: thành công=%d lỗi=%d, want 1/1", len(results), len(errs))
	}

	// Request thua race phải bị coi như reuse (cùng nhánh xử lý mất race).
	ae, ok := apperr.AsError(errs[0])
	if !ok || ae.Code != "refresh_reuse_detected" {
		t.Fatalf("lỗi race = %v, want refresh_reuse_detected", errs[0])
	}

	// Vì mất race coi như reuse → toàn bộ family bị revoke. Token con "thành
	// công" cũng không refresh tiếp được nữa.
	winner := results[0].RefreshToken
	if _, err := fx.svc.Refresh(ctx, winner); err == nil {
		t.Fatal("refresh bằng token con sau khi family bị revoke phải lỗi")
	}

	// Không còn token nào của family ở trạng thái chưa revoke.
	if notRevoked := countTokens(t, fx.pool, "revoked = false"); notRevoked != 0 {
		t.Fatalf("còn %d token chưa revoke sau reuse-race", notRevoked)
	}
}
