.PHONY: run build test clean docker-up docker-down

# Chạy ứng dụng
run:
	go run ./cmd/api

# Build binary
build:
	go build -o bin/server ./cmd/api

# Chạy tests
test:
	go test ./... -v

# Kiểm tra code
vet:
	go vet ./...

# Dọn dẹp
clean:
	rm -rf bin/

# Docker
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

# Tạo bảng trong database
migrate:
	@echo "Đang tạo bảng users..."
	@PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "\
		CREATE TABLE IF NOT EXISTS users ( \
			id UUID PRIMARY KEY, \
			full_name VARCHAR(100) NOT NULL, \
			email VARCHAR(255) NOT NULL UNIQUE, \
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), \
			updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW() \
		);"
	@echo "✅ Hoàn tất!"
