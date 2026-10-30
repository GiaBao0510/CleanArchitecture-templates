package config

import (
	"errors"
	"fmt"
	"slices"
)

var (
	validEnvironments = []string{"local", "development", "staging", "production"}
	validDBDrivers    = []string{"postgres", "mysql", "sqlite", "sqlserver"} // MariaDB dùng driver "mysql"
	validLogLevels    = []string{"debug", "info", "warn", "error"}
	validLogFormats   = []string{"json", "text"}
)

// Validate kiểm tra tính hợp lệ và trả về tất cả lỗi cùng lúc.
func (c *Config) Validate() error {
	var errs []error
	add := func(format string, args ...any) { errs = append(errs, fmt.Errorf(format, args...)) }

	// App
	if c.App.Name == "" {
		add("app.name is required")
	}
	if !slices.Contains(validEnvironments, c.App.Environment) {
		add("app.environment must be one of %v, got %q", validEnvironments, c.App.Environment)
	}
	if c.App.ShutdownTimeout <= 0 {
		add("app.shutdown_timeout must be > 0")
	}

	// HTTP / gRPC
	if c.Http.Address == "" {
		add("http.address is required")
	}
	if c.Http.ReadTimeout <= 0 || c.Http.WriteTimeout <= 0 || c.Http.IdleTimeout <= 0 {
		add("http.read_timeout, write_timeout and idle_timeout must be > 0")
	}
	if c.GRPC.Address == "" {
		add("grpc.address is required")
	}

	// CORS
	if c.CORS.AllowCredentials && slices.Contains(c.CORS.AllowedOrigins, "*") {
		add("cors: allow_credentials cannot be used with origin \"*\"")
	}

	// Database
	if !slices.Contains(validDBDrivers, c.Database.Driver) {
		add("database.driver must be one of %v, got %q", validDBDrivers, c.Database.Driver)
	}
	if c.Database.DSN == "" {
		add("database.dsn is required (set DB_DSN)")
	}
	if c.Database.MaxOpenConns <= 0 {
		add("database.max_open_conns must be > 0")
	}
	if c.Database.MaxIdleConns > c.Database.MaxOpenConns {
		add("database.max_idle_conns must not exceed max_open_conns")
	}

	// Redis
	if c.Redis.Address == "" {
		add("redis.address is required")
	}
	if c.Redis.CacheTTL <= 0 {
		add("redis.cache_ttl must be > 0")
	}
	if c.Redis.DB < 0 {
		add("redis.db must be >= 0")
	}

	// Kafka
	if len(c.Kafka.Brokers) == 0 {
		add("kafka.brokers is required")
	}
	if c.Kafka.UserTopic == "" {
		add("kafka.user_topic is required")
	}
	if c.Kafka.WorkerGroupID == "" {
		add("kafka.worker_group_id is required")
	}
	if !slices.Contains([]int{-1, 0, 1}, c.Kafka.Producer.RequiredAcks) {
		add("kafka.producer.required_acks must be -1, 0 or 1")
	}

	// Log
	if !slices.Contains(validLogLevels, c.Log.Level) {
		add("log.level must be one of %v, got %q", validLogLevels, c.Log.Level)
	}
	if !slices.Contains(validLogFormats, c.Log.Format) {
		add("log.format must be one of %v, got %q", validLogFormats, c.Log.Format)
	}

	// Auth
	if c.Auth.AccessTokenExpiry <= 0 || c.Auth.RefreshTokenExpiry <= c.Auth.AccessTokenExpiry {
		add("auth: access_token_ttl must be > 0 and smaller than refresh_token_ttl")
	}

	// Rate limit
	if c.RateLimit.Enabled && (c.RateLimit.RequestsPerSecond <= 0 || c.RateLimit.BurstSize <= 0) {
		add("rate_limit.requests_per_second and burst must be > 0 when enabled")
	}

	// Telemetry
	if c.Telemetry.SampleRate < 0 || c.Telemetry.SampleRate > 1 {
		add("telemetry.sample_rate must be between 0 and 1")
	}

	// Ràng buộc riêng cho staging/production
	if c.App.Environment == "staging" || c.App.Environment == "production" {
		if len(c.Auth.JWTSecret) < 32 {
			add("auth.jwt_secret must be at least 32 characters (set JWT_SECRET)")
		}
	}
	if c.App.Environment == "production" {
		if c.App.Debug {
			add("app.debug must be false in production")
		}
		if c.Http.EnableSwagger || c.GRPC.EnableReflection {
			add("http.enable_swagger and grpc.enable_reflection must be false in production")
		}
		if c.Redis.TLS.InsecureSkipVerify || c.GRPC.TLS.InsecureSkipVerify || c.Kafka.TLS.InsecureSkipVerify {
			add("tls.insecure_skip_verify must be false in production")
		}
	}

	return errors.Join(errs...)
}
