package authn

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// KeyPair là cặp khoá ký/verify access JWT (EC P-256 cho ES256).
type KeyPair struct {
	Private *ecdsa.PrivateKey
	Public  *ecdsa.PublicKey
}

// LoadKeyPair nạp cặp khoá EC P-256 từ PEM. Nếu CẢ HAI PEM rỗng → sinh cặp khoá
// ephemeral (chỉ phù hợp dev; gọi báo cảnh báo ở caller). Nếu chỉ một rỗng hoặc
// PEM sai → lỗi (cấu hình nửa vời nguy hiểm hơn không cấu hình).
//
// ephemeral=true khi khoá được sinh tạm; caller nên log cảnh báo.
func LoadKeyPair(privPEM, pubPEM string) (kp KeyPair, ephemeral bool, err error) {
	if privPEM == "" && pubPEM == "" {
		priv, gerr := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if gerr != nil {
			return KeyPair{}, false, fmt.Errorf("sinh khoá ephemeral: %w", gerr)
		}
		return KeyPair{Private: priv, Public: &priv.PublicKey}, true, nil
	}
	if privPEM == "" || pubPEM == "" {
		return KeyPair{}, false, fmt.Errorf("phải cung cấp CẢ JWT_PRIVATE_KEY_PEM và JWT_PUBLIC_KEY_PEM, hoặc bỏ trống cả hai")
	}

	priv, err := parseECPrivate(privPEM)
	if err != nil {
		return KeyPair{}, false, err
	}
	pub, err := parseECPublic(pubPEM)
	if err != nil {
		return KeyPair{}, false, err
	}
	return KeyPair{Private: priv, Public: pub}, false, nil
}

func parseECPrivate(pemStr string) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("private key: PEM không hợp lệ")
	}
	// Hỗ trợ cả SEC1 (EC PRIVATE KEY) lẫn PKCS#8.
	if key, err := x509.ParseECPrivateKey(block.Bytes); err == nil {
		return key, nil
	}
	anyKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("private key: parse thất bại: %w", err)
	}
	ecKey, ok := anyKey.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key: không phải EC (ECDSA)")
	}
	return ecKey, nil
}

func parseECPublic(pemStr string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("public key: PEM không hợp lệ")
	}
	anyKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("public key: parse thất bại: %w", err)
	}
	ecKey, ok := anyKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key: không phải EC (ECDSA)")
	}
	return ecKey, nil
}
