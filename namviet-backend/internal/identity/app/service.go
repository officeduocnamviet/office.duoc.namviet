// Package app chứa use-case của bounded context identity. Nó điều phối domain
// + port, MỞ/COMMIT transaction (qua TxManager → platform/db.WithinTx), và trả
// lỗi nghiệp vụ qua common/apperr. KHÔNG biết HTTP/Huma; KHÔNG viết SQL trực
// tiếp (chỉ gọi port domain).
package app

import (
	"context"
	"time"

	"github.com/Maneva-AI/namviet-backend/internal/common/apperr"
	"github.com/Maneva-AI/namviet-backend/internal/common/id"
	"github.com/Maneva-AI/namviet-backend/internal/identity/domain"
)

// refreshTTL là thời gian sống refresh token (dài hơn access nhiều).
const refreshTTL = 30 * 24 * time.Hour

// AuthService là use-case xác thực: login, refresh (xoay vòng + reuse-detection),
// logout, me.
type AuthService struct {
	users  domain.UserRepository
	roles  domain.RoleRepository
	txm    TxManager
	hasher *PasswordHasher
	issuer *TokenIssuer
	now    func() time.Time
}

// NewAuthService wiring use-case. users/roles dùng cho read ngoài tx; txm dùng
// cho thao tác refresh nguyên tử.
func NewAuthService(users domain.UserRepository, roles domain.RoleRepository, txm TxManager, hasher *PasswordHasher, issuer *TokenIssuer) *AuthService {
	return &AuthService{
		users:  users,
		roles:  roles,
		txm:    txm,
		hasher: hasher,
		issuer: issuer,
		now:    time.Now,
	}
}

// Tokens là kết quả phát token cho login/refresh.
type Tokens struct {
	AccessToken  string
	RefreshToken string // token thô (chỉ trả client một lần)
	ExpiresIn    int    // TTL access (giây)
}

// MeResult là thông tin danh tính + quyền của user hiện tại.
type MeResult struct {
	UserID      string
	Email       string
	UserType    string
	Permissions []string
}

// Login xác thực email/mật khẩu rồi phát access + refresh token. Sai thông tin
// → apperr.Unauthorized (không phân biệt email/mật khẩu để chống dò). Nếu hash
// là bcrypt legacy và verify đúng → lazy rehash sang argon2id. Việc rehash +
// phát refresh chạy trong MỘT transaction.
func (s *AuthService) Login(ctx context.Context, email, password string) (Tokens, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		if apperr.KindOf(err) == apperr.KindNotFound {
			return Tokens{}, errInvalidCredentials()
		}
		return Tokens{}, err
	}
	if !user.CanLogin() {
		return Tokens{}, errInvalidCredentials()
	}

	ok, needsRehash, err := s.hasher.Verify(password, user.PasswordHash, user.HashAlgo)
	if err != nil {
		return Tokens{}, err
	}
	if !ok {
		return Tokens{}, errInvalidCredentials()
	}

	perms, err := s.roles.PermissionCodesForUser(ctx, user.ID)
	if err != nil {
		return Tokens{}, err
	}

	access, err := s.issuer.IssueAccess(user.ID, user.UserType, perms)
	if err != nil {
		return Tokens{}, err
	}

	rawRefresh, hash, err := NewRefreshToken()
	if err != nil {
		return Tokens{}, err
	}
	familyID := id.NewString()

	// Rehash (nếu cần) + lưu refresh token mới trong cùng một transaction.
	err = s.txm.WithinTx(ctx, func(r Repos) error {
		if needsRehash {
			newHash, _, herr := s.hasher.Hash(password)
			if herr != nil {
				return herr
			}
			if uerr := r.Users.UpdatePasswordHash(ctx, user.ID, newHash, domain.HashArgon2id); uerr != nil {
				return uerr
			}
		}
		return r.Tokens.Insert(ctx, domain.RefreshToken{
			ID:        id.NewString(),
			UserID:    user.ID,
			TokenHash: hash,
			FamilyID:  familyID,
			ExpiresAt: s.now().Add(refreshTTL),
		})
	})
	if err != nil {
		return Tokens{}, err
	}

	return Tokens{AccessToken: access, RefreshToken: rawRefresh, ExpiresIn: int(accessTTL.Seconds())}, nil
}

// errReuseSentinel là lỗi NỘI BỘ để thoát tx xoay vòng khi phát hiện reuse.
// KHÔNG dùng RevokeFamily trong cùng tx đó vì trả lỗi sẽ rollback (hủy luôn
// revoke). Thay vào đó ta thoát tx (rollback các thao tác chỉ-đọc), rồi
// RevokeFamily trong một tx RIÊNG tự commit ở Refresh.
var errReuseSentinel = apperr.Internal("__reuse_sentinel__")

// Refresh xoay vòng refresh token:
//   - không thấy / hết hạn / đã revoke → apperr.Unauthorized("refresh_invalid");
//   - đã used (REUSE) → RevokeFamily (tx riêng, commit) + apperr.Unauthorized("refresh_reuse_detected");
//   - hợp lệ → MarkUsed(old) + Insert(new cùng family) + phát access mới, tất cả
//     trong MỘT transaction (nguyên tử).
func (s *AuthService) Refresh(ctx context.Context, rawRefresh string) (Tokens, error) {
	hash := HashRefreshToken(rawRefresh)

	var (
		out           Tokens
		reuseFamilyID string
	)
	err := s.txm.WithinTx(ctx, func(r Repos) error {
		tok, gerr := r.Tokens.GetByHash(ctx, hash)
		if gerr != nil {
			if apperr.KindOf(gerr) == apperr.KindNotFound {
				return errRefreshInvalid()
			}
			return gerr
		}

		// Reuse-detection: token đã dùng lại. Ghi nhận family rồi thoát tx này
		// (revoke ở tx riêng để không bị rollback bởi lỗi trả về).
		if tok.Used {
			reuseFamilyID = tok.FamilyID
			return errReuseSentinel
		}

		if !tok.IsActive(s.now()) {
			return errRefreshInvalid()
		}

		// Đánh dấu token cũ đã dùng (one-time use) bằng CLAIM NGUYÊN TỬ. Nếu
		// MarkUsed trả Conflict nghĩa là một request refresh ĐỒNG THỜI đã claim
		// token này trước (READ COMMITTED: cả hai cùng đọc used=false, nhưng chỉ
		// một thắng UPDATE ... WHERE used=false). Đây thực chất là dùng lại token
		// đã tiêu thụ → xử lý GIỐNG nhánh reuse: ghi nhận family rồi thoát tx để
		// revoke cả family ở tx riêng (xem khối bên dưới).
		if merr := r.Tokens.MarkUsed(ctx, tok.ID); merr != nil {
			if apperr.KindOf(merr) == apperr.KindConflict {
				reuseFamilyID = tok.FamilyID
				return errReuseSentinel
			}
			return merr
		}

		user, uerr := r.Users.GetByID(ctx, tok.UserID)
		if uerr != nil {
			return uerr
		}
		if !user.CanLogin() {
			return errRefreshInvalid()
		}

		perms, perr := r.Roles.PermissionCodesForUser(ctx, user.ID)
		if perr != nil {
			return perr
		}

		access, aerr := s.issuer.IssueAccess(user.ID, user.UserType, perms)
		if aerr != nil {
			return aerr
		}

		rawNew, newHash, nerr := NewRefreshToken()
		if nerr != nil {
			return nerr
		}
		// Token mới thuộc CÙNG family để reuse-detection bắt được nếu token cũ
		// bị tái dùng sau này.
		if ierr := r.Tokens.Insert(ctx, domain.RefreshToken{
			ID:        id.NewString(),
			UserID:    user.ID,
			TokenHash: newHash,
			FamilyID:  tok.FamilyID,
			ExpiresAt: s.now().Add(refreshTTL),
		}); ierr != nil {
			return ierr
		}

		out = Tokens{AccessToken: access, RefreshToken: rawNew, ExpiresIn: int(accessTTL.Seconds())}
		return nil
	})

	// Reuse phát hiện: thu hồi cả family trong một tx RIÊNG (commit) rồi báo lỗi.
	// Tách khỏi tx trên để revoke không bị rollback bởi việc trả lỗi.
	if reuseFamilyID != "" {
		if rerr := s.txm.WithinTx(ctx, func(r Repos) error {
			return r.Tokens.RevokeFamily(ctx, reuseFamilyID)
		}); rerr != nil {
			return Tokens{}, rerr
		}
		return Tokens{}, errRefreshReuse()
	}

	if err != nil {
		return Tokens{}, err
	}
	return out, nil
}

// Logout thu hồi toàn bộ family của refresh token được cung cấp. Token không
// tồn tại được coi là logout thành công (idempotent) — không rò trạng thái.
func (s *AuthService) Logout(ctx context.Context, rawRefresh string) error {
	hash := HashRefreshToken(rawRefresh)
	return s.txm.WithinTx(ctx, func(r Repos) error {
		tok, err := r.Tokens.GetByHash(ctx, hash)
		if err != nil {
			if apperr.KindOf(err) == apperr.KindNotFound {
				return nil // idempotent
			}
			return err
		}
		return r.Tokens.RevokeFamily(ctx, tok.FamilyID)
	})
}

// Me trả thông tin + quyền của user theo id (lấy từ access token đã verify).
func (s *AuthService) Me(ctx context.Context, userID string) (MeResult, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return MeResult{}, err
	}
	perms, err := s.roles.PermissionCodesForUser(ctx, userID)
	if err != nil {
		return MeResult{}, err
	}
	return MeResult{
		UserID:      user.ID,
		Email:       user.Email,
		UserType:    user.UserType,
		Permissions: perms,
	}, nil
}
