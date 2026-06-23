package accounting_test

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Maneva-AI/namviet-backend/internal/accounting"
	"github.com/Maneva-AI/namviet-backend/internal/accounting/domain"
	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/money"
	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
	"github.com/Maneva-AI/namviet-backend/internal/platform/db/dbtest"
	"github.com/Maneva-AI/namviet-backend/internal/platform/httpx/humax"
)

type fixture struct {
	mod  *accounting.Module
	pool *pgxpool.Pool
	key  *ecdsa.PrivateKey
}

const entryDate = "2026-06-15"

// seed nạp cây tài khoản TT133 (allow_posting) + 1 TK tổng hợp KHÔNG cho hạch
// toán + 1 kỳ kế toán 'open' (2026-06) và 1 kỳ 'closed' (2026-05).
func seed(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	stmts := []string{
		`INSERT INTO public.chart_of_accounts (account_code, name, type, balance_type, allow_posting) VALUES
			('131','Phải thu khách hàng','Tài sản','DEBIT', true),
			('511','Doanh thu bán hàng','Doanh thu','CREDIT', true),
			('3331','Thuế GTGT đầu ra','Nợ','CREDIT', true),
			('632','Giá vốn hàng bán','Chi phí','DEBIT', true),
			('156','Hàng hoá','Tài sản','DEBIT', true),
			('111','Tiền mặt','Tài sản','DEBIT', true),
			('5','TK tổng hợp (không hạch toán)','Doanh thu','CREDIT', false)`,
		`INSERT INTO app.accounting_periods (id, year, month, status) VALUES
			(gen_random_uuid(), 2026, 6, 'open'),
			(gen_random_uuid(), 2026, 5, 'closed')`,
	}
	for _, s := range stmts {
		if _, err := pool.Exec(ctx, s); err != nil {
			t.Fatalf("seed: %v\nSQL: %s", err, s)
		}
	}
}

func setup(t *testing.T) fixture {
	t.Helper()
	pool, cleanup := dbtest.NewPool(t)
	t.Cleanup(cleanup)
	seed(t, pool)
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("gen key: %v", err)
	}
	return fixture{mod: accounting.NewModule(pool), pool: pool, key: key}
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := time.Parse("2006-01-02", s)
	if err != nil {
		t.Fatalf("parse date %q: %v", s, err)
	}
	return d
}

// postInTx chạy Poster.Post trong MỘT tx (mô phỏng caller orders/finance gộp
// atomic) và commit. Trả entryID + lỗi (gồm lỗi commit do trigger cân Σ).
func postInTx(t *testing.T, fx fixture, e domain.JournalEntry) (string, error) {
	t.Helper()
	ctx := context.Background()
	tx, err := fx.pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	id, postErr := fx.mod.Poster().Post(ctx, tx, e)
	if postErr != nil {
		_ = tx.Rollback(ctx)
		return "", postErr
	}
	if err := tx.Commit(ctx); err != nil {
		return "", err // trigger DEFERRABLE cân Σ FAIL lúc commit
	}
	return id, nil
}

// ---- POST (qua Poster, tx caller) ----

func TestIntegration_Post_Balanced(t *testing.T) {
	fx := setup(t)
	e := domain.JournalEntry{
		Book: domain.BookInternal, EntryDate: mustDate(t, entryDate),
		SourceType: "order", SourceID: "HD123", Memo: "bán hàng có VAT",
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(110000)},
			{AccountCode: "511", Credit: money.FromInt(100000)},
			{AccountCode: "3331", Credit: money.FromInt(10000)},
		},
	}
	id, err := postInTx(t, fx, e)
	if err != nil {
		t.Fatalf("post cân phải thành công: %v", err)
	}
	// Đọc lại: Σ khớp.
	rec, err := fx.mod.Service().GetEntry(context.Background(), id)
	if err != nil {
		t.Fatalf("GetEntry: %v", err)
	}
	if rec.Book != domain.BookInternal || len(rec.Lines) != 3 {
		t.Fatalf("entry đọc lại sai: book=%s lines=%d", rec.Book, len(rec.Lines))
	}
	var d, c money.Money = money.Zero(), money.Zero()
	for _, l := range rec.Lines {
		d = d.Add(l.Debit)
		c = c.Add(l.Credit)
	}
	if !d.Equal(money.FromInt(110000)) || !c.Equal(money.FromInt(110000)) {
		t.Fatalf("Σ đọc lại lệch: d=%s c=%s", d, c)
	}
}

func TestIntegration_Post_Unbalanced_CommitFails(t *testing.T) {
	fx := setup(t)
	// Bút toán "lệch" — nhưng domain Validate sẽ chặn trước. Để thực sự chạm
	// trigger DB, ta phải bypass Validate bằng cách insert thẳng qua tx một entry
	// 1 dòng (debit không có credit đối ứng) rồi commit → trigger FAIL.
	ctx := context.Background()
	tx, err := fx.pool.Begin(ctx)
	if err != nil {
		t.Fatalf("begin: %v", err)
	}
	var periodID string
	if err := tx.QueryRow(ctx,
		`SELECT id FROM app.accounting_periods WHERE year=2026 AND month=6`).Scan(&periodID); err != nil {
		t.Fatalf("get period: %v", err)
	}
	var entryID string
	if err := tx.QueryRow(ctx,
		`INSERT INTO app.journal_entries (book, entry_date, period_id) VALUES ('INTERNAL', $1::date, $2) RETURNING id`,
		entryDate, periodID).Scan(&entryID); err != nil {
		t.Fatalf("insert entry: %v", err)
	}
	// Một dòng debit 100, KHÔNG có credit đối ứng → Σ lệch.
	if _, err := tx.Exec(ctx,
		`INSERT INTO app.journal_entry_lines (entry_id, line_no, account_code, debit, credit) VALUES ($1, 1, '131', 100, 0)`,
		entryID); err != nil {
		t.Fatalf("insert line: %v", err)
	}
	// Commit phải FAIL vì trigger DEFERRABLE cân Σ.
	if err := tx.Commit(ctx); err == nil {
		t.Fatal("commit bút toán lệch phải FAIL (trigger cân Σ)")
	}
}

func TestIntegration_Post_PeriodClosed_Conflict(t *testing.T) {
	fx := setup(t)
	// entry_date thuộc kỳ 2026-05 (đã 'closed').
	e := domain.JournalEntry{
		Book: domain.BookInternal, EntryDate: mustDate(t, "2026-05-20"),
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(100)},
			{AccountCode: "511", Credit: money.FromInt(100)},
		},
	}
	_, err := postInTx(t, fx, e)
	if err == nil {
		t.Fatal("post vào kỳ khoá phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindConflict {
		t.Fatalf("kỳ khoá phải Conflict, got %v (%v)", apperr.KindOf(err), err)
	}
}

func TestIntegration_Post_AccountNotPostable_Unprocessable(t *testing.T) {
	fx := setup(t)
	// TK '5' tồn tại nhưng allow_posting = false.
	e := domain.JournalEntry{
		Book: domain.BookInternal, EntryDate: mustDate(t, entryDate),
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(100)},
			{AccountCode: "5", Credit: money.FromInt(100)},
		},
	}
	_, err := postInTx(t, fx, e)
	if err == nil {
		t.Fatal("post vào TK tổng hợp (allow_posting=false) phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("TK không postable phải Validation, got %v (%v)", apperr.KindOf(err), err)
	}
}

func TestIntegration_Post_AccountMissing_Unprocessable(t *testing.T) {
	fx := setup(t)
	e := domain.JournalEntry{
		Book: domain.BookInternal, EntryDate: mustDate(t, entryDate),
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(100)},
			{AccountCode: "999", Credit: money.FromInt(100)}, // không tồn tại
		},
	}
	_, err := postInTx(t, fx, e)
	if err == nil {
		t.Fatal("post vào TK không tồn tại phải lỗi")
	}
	if apperr.KindOf(err) != apperr.KindValidation {
		t.Fatalf("TK không tồn tại phải Validation, got %v", apperr.KindOf(err))
	}
}

// 2 sổ INTERNAL/TAX cùng nghiệp vụ (source order HD777) với SỐ TIỀN KHÁC nhau →
// cả 2 tồn tại độc lập, KHÔNG sync.
func TestIntegration_Post_TwoBooks_Independent(t *testing.T) {
	fx := setup(t)
	internal := domain.JournalEntry{
		Book: domain.BookInternal, EntryDate: mustDate(t, entryDate),
		SourceType: "order", SourceID: "HD777",
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(110000)},
			{AccountCode: "511", Credit: money.FromInt(100000)},
			{AccountCode: "3331", Credit: money.FromInt(10000)},
		},
	}
	tax := domain.JournalEntry{
		Book: domain.BookTax, EntryDate: mustDate(t, entryDate),
		SourceType: "order", SourceID: "HD777",
		Lines: []domain.EntryLine{
			{AccountCode: "131", Debit: money.FromInt(99000)}, // giá HĐ khác giá thực
			{AccountCode: "511", Credit: money.FromInt(90000)},
			{AccountCode: "3331", Credit: money.FromInt(9000)},
		},
	}
	idIn, err := postInTx(t, fx, internal)
	if err != nil {
		t.Fatalf("INTERNAL: %v", err)
	}
	idTax, err := postInTx(t, fx, tax)
	if err != nil {
		t.Fatalf("TAX: %v", err)
	}
	if idIn == idTax {
		t.Fatal("2 entry phải có id khác nhau")
	}
	recIn, _ := fx.mod.Service().GetEntry(context.Background(), idIn)
	recTax, _ := fx.mod.Service().GetEntry(context.Background(), idTax)
	if recIn.Book != domain.BookInternal || recTax.Book != domain.BookTax {
		t.Fatalf("book sai: %s / %s", recIn.Book, recTax.Book)
	}
	// Số tiền KHÁC nhau (không sync).
	var dIn, dTax money.Money = money.Zero(), money.Zero()
	for _, l := range recIn.Lines {
		dIn = dIn.Add(l.Debit)
	}
	for _, l := range recTax.Lines {
		dTax = dTax.Add(l.Debit)
	}
	if dIn.Equal(dTax) {
		t.Fatal("2 sổ phải độc lập số tiền (INTERNAL 110000 vs TAX 99000)")
	}
}

// Property test: sinh N cặp dòng CÂN (vary theo INDEX, KHÔNG dùng math/rand) →
// luôn post được & Σ cân. Thêm 1 đồng lệch (insert thẳng) → commit luôn FAIL.
func TestIntegration_Post_Property_AlwaysBalanced(t *testing.T) {
	fx := setup(t)
	ctx := context.Background()
	for i := 1; i <= 25; i++ {
		// Số tiền biến thiên theo index (deterministic): amount = i * 1000 + i.
		amount := int64(i)*1000 + int64(i)
		e := domain.JournalEntry{
			Book: domain.BookInternal, EntryDate: mustDate(t, entryDate),
			SourceType: "prop", SourceID: fmt.Sprintf("P%03d", i),
			Lines: []domain.EntryLine{
				{AccountCode: "131", Debit: money.FromInt(amount)},
				{AccountCode: "511", Credit: money.FromInt(amount)},
			},
		}
		id, err := postInTx(t, fx, e)
		if err != nil {
			t.Fatalf("property post cân (i=%d, amount=%d) phải thành công: %v", i, amount, err)
		}
		rec, err := fx.mod.Service().GetEntry(ctx, id)
		if err != nil {
			t.Fatalf("GetEntry property i=%d: %v", i, err)
		}
		var d, c money.Money = money.Zero(), money.Zero()
		for _, l := range rec.Lines {
			d = d.Add(l.Debit)
			c = c.Add(l.Credit)
		}
		if !d.Equal(c) {
			t.Fatalf("property i=%d Σ lệch: d=%s c=%s", i, d, c)
		}
	}

	// Lệch 1 đồng: insert thẳng (bypass domain) → trigger luôn FAIL lúc commit.
	for i := 1; i <= 10; i++ {
		amount := int64(i) * 1000
		tx, err := fx.pool.Begin(ctx)
		if err != nil {
			t.Fatalf("begin i=%d: %v", i, err)
		}
		var periodID, entryID string
		_ = tx.QueryRow(ctx, `SELECT id FROM app.accounting_periods WHERE year=2026 AND month=6`).Scan(&periodID)
		if err := tx.QueryRow(ctx,
			`INSERT INTO app.journal_entries (book, entry_date, period_id) VALUES ('INTERNAL',$1::date,$2) RETURNING id`,
			entryDate, periodID).Scan(&entryID); err != nil {
			t.Fatalf("insert entry i=%d: %v", i, err)
		}
		// debit = amount, credit = amount+1 → lệch đúng 1 đồng.
		_, _ = tx.Exec(ctx, `INSERT INTO app.journal_entry_lines (entry_id,line_no,account_code,debit,credit) VALUES ($1,1,'131',$2,0)`, entryID, amount)
		_, _ = tx.Exec(ctx, `INSERT INTO app.journal_entry_lines (entry_id,line_no,account_code,debit,credit) VALUES ($1,2,'511',0,$2)`, entryID, amount+1)
		if err := tx.Commit(ctx); err == nil {
			t.Fatalf("property lệch i=%d phải FAIL lúc commit", i)
		}
	}
}

// ---- HTTP đọc (envelope + authz) ----

func TestIntegration_HTTP_Accounting_ReadAndAuthz(t *testing.T) {
	fx := setup(t)
	// Seed vài bút toán để list.
	for i := 1; i <= 3; i++ {
		e := domain.JournalEntry{
			Book: domain.BookInternal, EntryDate: mustDate(t, entryDate),
			SourceType: "order", SourceID: fmt.Sprintf("HD%03d", i),
			Lines: []domain.EntryLine{
				{AccountCode: "131", Debit: money.FromInt(int64(i) * 1000)},
				{AccountCode: "511", Credit: money.FromInt(int64(i) * 1000)},
			},
		}
		if _, err := postInTx(t, fx, e); err != nil {
			t.Fatalf("seed entry %d: %v", i, err)
		}
	}
	handler := buildHTTP(t, fx)
	tok := signToken(t, fx.key, "accounting.read")

	// Thiếu token → 401.
	noTok := doReq(t, handler, "/v1/accounting/entries", "")
	if noTok.code != http.StatusUnauthorized {
		t.Fatalf("no token = %d, want 401; body=%s", noTok.code, noTok.body)
	}
	assertErrCode(t, noTok.body, "unauthorized")

	// Sai quyền → 403.
	wrong := doReq(t, handler, "/v1/accounting/entries", signToken(t, fx.key, "orders.read"))
	if wrong.code != http.StatusForbidden {
		t.Fatalf("wrong perm = %d, want 403; body=%s", wrong.code, wrong.body)
	}
	assertErrCode(t, wrong.body, "forbidden")

	// Có quyền → 200 + envelope {data:{items}, error:null}, debit/credit là chuỗi.
	ok := doReq(t, handler, "/v1/accounting/entries?book=INTERNAL", tok)
	if ok.code != http.StatusOK {
		t.Fatalf("entries = %d, want 200; body=%s", ok.code, ok.body)
	}
	var env struct {
		Data struct {
			Items []struct {
				ID   string `json:"id"`
				Book string `json:"book"`
			} `json:"items"`
			NextCursor string `json:"next_cursor"`
		} `json:"data"`
		Error any `json:"error"`
	}
	if err := json.Unmarshal([]byte(ok.body), &env); err != nil {
		t.Fatalf("decode: %v (body=%s)", err, ok.body)
	}
	if env.Error != nil || len(env.Data.Items) != 3 {
		t.Fatalf("entries envelope sai (want 3 items): %s", ok.body)
	}

	// GetEntry: debit/credit ra chuỗi thập phân (không float).
	gid := env.Data.Items[0].ID
	g := doReq(t, handler, "/v1/accounting/entries/"+gid, tok)
	if g.code != http.StatusOK {
		t.Fatalf("get entry = %d; body=%s", g.code, g.body)
	}
	var genv struct {
		Data struct {
			ID    string `json:"id"`
			Book  string `json:"book"`
			Lines []struct {
				AccountCode string `json:"account_code"`
				Debit       string `json:"debit"`
				Credit      string `json:"credit"`
			} `json:"lines"`
		} `json:"data"`
	}
	if err := json.Unmarshal([]byte(g.body), &genv); err != nil {
		t.Fatalf("decode get: %v (body=%s)", err, g.body)
	}
	if len(genv.Data.Lines) != 2 {
		t.Fatalf("get entry lines = %d, want 2; body=%s", len(genv.Data.Lines), g.body)
	}
	// Mỗi dòng đúng một vế là chuỗi số (string), không float JSON.
	for _, l := range genv.Data.Lines {
		if l.Debit == "" || l.Credit == "" {
			t.Fatalf("debit/credit phải là chuỗi: %+v", l)
		}
	}

	// GetEntry id không tồn tại → 404.
	nf := doReq(t, handler, "/v1/accounting/entries/00000000-0000-0000-0000-000000000000", tok)
	if nf.code != http.StatusNotFound {
		t.Fatalf("not found = %d, want 404; body=%s", nf.code, nf.body)
	}
}

// ---- helpers ----

func buildHTTP(t *testing.T, fx fixture) http.Handler {
	t.Helper()
	r := chi.NewMux()
	api := humax.New(r, "test", "1.0.0")
	verifier := authn.NewVerifier(&fx.key.PublicKey)
	fx.mod.RegisterRoutes(api, verifier)
	return r
}

func signToken(t *testing.T, key *ecdsa.PrivateKey, perms ...string) string {
	t.Helper()
	claims := authn.Claims{
		Permissions: perms,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   "u1",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	s, err := jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

type httpResp struct {
	code int
	body string
}

func doReq(t *testing.T, h http.Handler, path, bearer string) httpResp {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, path, nil)
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return httpResp{code: rec.Code, body: rec.Body.String()}
}

func assertErrCode(t *testing.T, body, want string) {
	t.Helper()
	var env struct {
		Error *struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("decode envelope: %v (body=%s)", err, body)
	}
	if env.Error == nil || env.Error.Code != want {
		t.Fatalf("error.code = %v, want %q (body=%s)", env.Error, want, body)
	}
}

var _ = pgx.ErrNoRows // giữ import pgx nếu helper tx không tham chiếu trực tiếp
