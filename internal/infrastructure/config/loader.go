package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Nạp cấu hình theo thứ tự: mặc định < YAML < ENV
func LoadConfig(filePath string) (*Config, error) {

	// Đọc file dựa trên đường dẫn
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := defaultConfig()

	// KnownFields(true): key thừa hoặc sai chính tả trong YAML sẽ báo lỗi.
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(cfg); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to parse config file %q: %w", filePath, err)
	}

	if err := appEnvOverrides(cfg); err != nil {
		return nil, fmt.Errorf("failed to apply environment variable overrides: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

// DefaultConfig trả về cấu hình mặc định
func defaultConfig() *Config {
	return &Config{
		App: AppConfig{
			Environment:      "local",
			ShutdownTimeout: 15 * time.Second,
		},
		Http: HttpConfig{
			Address:           ":8080",
			ReadTimeout:       10 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			MaxHeaderBytes:    1 << 20, // 1 MB
			MaxBodyBytes:      4 << 20, // 4MB
		},
		GRPC: GRPCconfig{
			Address:          ":50051",
			Timeout:          10 * time.Second,
			MaxRecvMsgSize:   4 << 20,
			MaxSendMsgSize:   4 << 20,
			KeepaliveTime:    2 * time.Hour,
			KeepaliveTimeout: 20 * time.Second,
		},
		Database: DatabaseConfig{
			Driver:          "postgres",
			MaxOpenConns:    25,
			MaxIdleConns:    10,
			ConnMaxLifetime: 30 * time.Minute,
			ConnMaxIdleTime: 5 * time.Minute,
			ConnectTimeout:  5 * time.Second,
			QueryTimeout:    10 * time.Second,
		},
		Redis: RedisConfig{
			PoolSize:     10,
			MinIdleConns: 2,
			DialTimeOut:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
			CacheTTL:     10 * time.Minute,
		},
		Kafka: KafkaConfig{
			Consumer: KafkaConsumerConfig{
				AutoOffsetReset:   "latest",
				SessionTimeout:    10 * time.Second,
				HeartbeatInterval: 3 * time.Second,
				MaxPollRecords:    500,
				CommitInterval:    time.Second,
				MaxRetries:        3,
			},
			Producer: KafkaProducerConfig{
				RequiredAcks: -1,
				MaxRetries:   3,
				Compression:  "none",
				BatchSize:    100,
				BatchTimeout: 10 * time.Millisecond,
				WriteTimeout: 10 * time.Second,
			},
		},
		Log: LogConfig{Level: "info", Format: "json", Output: "stdout"},
		Auth: AuthConfig{
			AccessTokenExpiry:  15 * time.Minute,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		Telemetry: TelemetryConfig{
			MetricsPath: "/metrics",
			SampleRate:  0.1,
		},
	}
}

// applyEnvOverrides ghi đè cấu hình bằng biến môi trường (nếu có và không rỗng)
func appEnvOverrides(cfg *Config) error {
	e := &envReader{}

	// App & server
	e.str("APP_ENV", &cfg.App.Environment)
	e.str("HTTP_ADDRESS", &cfg.Http.Address)
	e.str("GRPC_ADDRESS", &cfg.GRPC.Address)

	// database
	e.str("DB_DRIVER", &cfg.Database.Driver)
	e.str("DB_DNS", &cfg.Database.DSN)

	// redis
	e.str("REDIS_ADDRESS", &cfg.Redis.Address)
	e.str("REDIS_USERNAME", &cfg.Redis.Username)
	e.str("REDIS_PASSWORD", &cfg.Redis.Password)
	e.integer("REDIS_DB", &cfg.Redis.DB)
	e.duration("REDIS_CACHE_TTL", &cfg.Redis.CacheTTL)

	//Kafka
	e.list("KAFKA_BROKERS", &cfg.Kafka.Brokers)
	e.str("KAFKA_USER_TOPIC", &cfg.Kafka.UserTopic)
	e.str("KAFKA_WORKER_GROUP_ID", &cfg.Kafka.WorkerGroupID)
	e.str("KAFKA_SASL_USERNAME", &cfg.Kafka.SASL.Username)
	e.str("KAFKA_SASL_PASSWORD", &cfg.Kafka.SASL.Password)

	// Auth & Log
	e.str("JWT_SECRET", &cfg.Auth.JWTSecret)
	e.str("LOG_LEVEL", &cfg.Log.Level)
	e.str("LOG_FORMAT", &cfg.Log.Format)

	return errors.Join(e.errs...)
}

// envReader gom các helper đọc biến môi trường theo từng kiểu dữ liệu.
type envReader struct{ errs []error }

func (r *envReader) str(key string, targer *string) { // Đọc biến môi trường kiểu string
	if v := os.Getenv(key); v != "" {
		*targer = v
	}
}

func (r *envReader) integer(key string, target *int) {
	v := os.Getenv(key)
	if v == "" {
		return
	}

	n, err := strconv.Atoi(v)
	if err != nil {
		r.errs = append(r.errs, fmt.Errorf("invalid integer value for %s: %v", key, err))
		return
	}

	*target = n
}

func (r *envReader) duration(key string, target *time.Duration) {
	v := os.Getenv(key)
	if v == "" {
		return
	}

	d, err := time.ParseDuration(v)
	if err != nil {
		r.errs = append(r.errs, fmt.Errorf("invalid duration value for %s: %v", key, err))
		return
	}
	*target = d
}

func (r *envReader) list(key string, target *[]string) {
	v := os.Getenv(key)
	if v == "" {
		return
	}

	var out []string
	for _, p := range strings.Split(v, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	*target = out
}
