# Build Stage
FROM golang:1.22-alpine AS builder

# Thiết lập thư mục làm việc trong container
WORKDIR /app

# Copy các file cấu hình go modules trước để cache layer
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code vào container
COPY . .

# Biên dịch ứng dụng Go thành một executable duy nhất (không phụ thuộc CGO)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./cmd/api/main.go

# Run Stage (Sử dụng Alpine siêu nhẹ để chạy ứng dụng)
FROM alpine:latest  

# Cài đặt ca-certificates để ứng dụng có thể gọi API HTTPS ra ngoài (Firebase, Supabase...)
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy file chạy từ builder stage sang
COPY --from=builder /app/api-server .

# Copy thư mục docs (Swagger) và configs (Nếu có)
COPY --from=builder /app/docs ./docs
# COPY --from=builder /app/configs ./configs 

# Mở cổng 8080 để Cloud Run mapping
EXPOSE 8080

# Chạy server
CMD ["./api-server"]
