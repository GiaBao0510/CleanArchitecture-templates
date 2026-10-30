package database

import (
	"context"
	"fmt"

	"github.com/GiaBao0510/CleanArchitecture-templates/internal/infrastructure/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(
	ctx context.Context, 
	cfg *config.DatabaseConfig,
) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN) // Thực hiên phân tích DSN để tạo cấu hình kết nối
	if err != nil {
		return nil, fmt.Errorf("failed to parse DSN: %w", err)
	}

	// Cấu hình các tham số kết nối từ cấu hình
	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MaxConnLifetime = cfg.MaxLifetime
	poolConfig.MaxConnIdleTime = cfg.ConnMaxIdleTime
	
	// Tạo một pool kết nối mới với cấu hình đã phân tích
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("Ping database failed: %w", err)
	}

	return pool, nil
}	