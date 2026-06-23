.PHONY: build test test-short run lint vuln sqlc openapi

build:
	go build -o bin/api ./cmd/api

# -p 1 serialize các package: testcontainers spin Postgres riêng cho mỗi package
# integration; chạy song song dễ gây race khi dò Docker provider (đặc biệt trên
# Docker Desktop/Windows). Serialize đảm bảo integration ổn định.
test:
	go test ./... -race -count=1 -p 1

test-short:
	go test ./... -short -count=1

run:
	go run ./cmd/api

lint:
	golangci-lint run

vuln:
	govulncheck ./...

sqlc:
	sqlc generate

# openapi sinh artifact api/openapi.yaml từ Huma (source of truth = handler).
# Đây là HỢP ĐỒNG để FE (ERP/Portal) tự sinh TS client trong repo của FE —
# repo Go KHÔNG chứa Node/TS tooling.
openapi:
	go run ./cmd/api -dump-openapi > api/openapi.yaml
