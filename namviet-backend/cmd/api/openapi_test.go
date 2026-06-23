package main

import (
	"io"
	"testing"

	"github.com/Maneva-AI/namviet-backend/internal/platform/authn"
)

// TestDumpSpec_AllModulesRegisterWithoutCollision dựng OpenAPI với TẤT CẢ bounded
// context đăng ký lên CÙNG một huma.API (như server thật khởi động) — chặn lớp lỗi
// trùng tên schema giữa các module (vd 2 package `http` cùng có type `vatLineDTO` →
// Huma panic "duplicate name" lúc đăng ký). Test từng-module KHÔNG bắt được vì mỗi
// test chỉ đăng ký module của nó; chỉ khi gộp mọi module mới lộ. Lỗi này panic lúc
// SERVER START (buildModules) nên phải có guard ở CI.
func TestDumpSpec_AllModulesRegisterWithoutCollision(t *testing.T) {
	// LoadKeyPair("","") → khoá ephemeral, đủ để dựng schema /v1/auth/* (không cần DB).
	if _, _, err := authn.LoadKeyPair("", ""); err != nil {
		t.Fatalf("load key: %v", err)
	}
	// dumpSpec gọi buildModules(nil, keys) đăng ký mọi module rồi sinh YAML. Panic
	// (trùng schema) sẽ làm test FAIL; lỗi trả về cũng FAIL.
	if err := dumpSpec(io.Discard); err != nil {
		t.Fatalf("dump OpenAPI (đăng ký tất cả module) lỗi: %v", err)
	}
}
