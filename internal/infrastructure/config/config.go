package config

import "time"

type Config struct {
	App       AppConfig       `yaml:"app"`
	Http      HttpConfig      `yaml:"http"`
	CORS      CORS_Config     `yaml:"cors"`
	Database  DatabaseConfig  `yaml:"database"`
	Redis     RedisConfig     `yaml:"redis"`
	GRPC      GRPCconfig      `yaml:"grpc"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	Log       LogConfig       `yaml:"log"`
	Auth      AuthConfig      `yaml:"auth"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
	Telemetry TelemetryConfig `yaml:"telemetry"`
}

// AppConfig chứa thông tin cấu hình của ứng dụng,
// bao gồm tên ứng dụng, môi trường và thời gian chờ khi tắt ứng dụng.
type AppConfig struct {
	Name            string        `yaml:"name"`
	Environment     string        `yaml:"environment"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout"`
	Version         string        `yaml:"version"`
	Debug           bool          `yaml:"debug"`
}

// HttpConfig chứa thông tin cấu hình của HTTP server,
// bao gồm địa chỉ, thời gian chờ đọc, thời gian chờ ghi và thời gian chờ nhàn rỗi.
type HttpConfig struct {
	Address           string        `yaml:"address"`
	ReadTimeout       time.Duration `yaml:"read_timeout"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout"`
	WriteTimeout      time.Duration `yaml:"write_timeout"`
	IdleTimeout       time.Duration `yaml:"idle_timeout"`
	MaxHeaderBytes    int           `yaml:"max_header_bytes"`   // Kích thước header tối đa
	MaxBodyBytes      int64         `yaml:"max_body_bytes"`     // Kích thước body request tối đa
	EnableSwagger     bool          `yaml:"enable_compression"` // Chỉ nên bật ở dev/staging
}

type CORS_Config struct {
	AllowedOrigins   []string `yaml:"allowed_origins"`   // Danh sách các nguồn được phép truy cập (ví dụ: ["http://localhost:3000", "https://example.com"])
	AllowedMethods   []string `yaml:"allowed_methods"`   // Danh sách các phương thức HTTP được phép (ví dụ: ["GET", "POST", "PUT", "DELETE"])
	AllowedHeaders   []string `yaml:"allowed_headers"`   // Danh sách các header được phép (ví dụ: ["Content-Type", "Authorization"])
	ExposedHeaders   []string `yaml:"exposed_headers"`   // Danh sách các header được phép hiển thị cho client (ví dụ: ["X-Custom-Header"])
	AllowCredentials bool     `yaml:"allow_credentials"` // Cho phép gửi cookie và thông tin xác thực (true/false)
	MaxAge           int      `yaml:"max_age"`           // Thời gian tối đa (tính bằng giây) mà trình duyệt có thể lưu trữ thông tin CORS (ví dụ: 3600)
}

type DatabaseConfig struct {
	Driver string `yaml:"driver"` // Postgres, MySQL, SQLite, ...

	DSN string `yaml:"dsn"` // Data Source Name, ví dụ: "user:password@tcp(localhost:3306)/dbname"

	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`  // Thời gian sống tối đa của một kết nối
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"` // Thời gian tối đa mà một kết nối có thể nhàn rỗi trước khi bị đóng
	ConnectTimeout  time.Duration `yaml:"connect_timeout"`    // Thời gian chờ kết nối cơ sở dữ liệu
	MaxOpenConns    int           `yaml:"max_open_conns"`     // Số lượng kết nối tối đa mở ra
	MaxIdleConns    int           `yaml:"max_idle_conns"`     // Số lượng kết nối tối đa nhàn rỗi
	MaxLifetime     time.Duration `yaml:"max_lifetime"`       // Thời gian sống tối đa của một kết nối
	QueryTimeout    time.Duration `yaml:"query_timeout"`      // Thời gian chờ tối đa cho một truy vấn cơ sở dữ liệu

	LogQueries bool            `yaml:"log_queries"` // Bật/tắt ghi log các truy vấn cơ sở dữ liệu (true/false)
	Migration  MigrationConfig `yaml:"migration"`   // Cấu hình migration cơ sở dữ liệu
}

type MigrationConfig struct {
	Enable bool   `yaml:"enable"` // Bật/tắt migration tự động (true/false)
	Dir    string `yaml:"dir"`    // Thư mục chứa các file migration (ví dụ: "./migrations")
}

type RedisConfig struct {
	Address  string `yaml:"address"`  // Địa chỉ của Redis server
	Username string `yaml:"username"` // Tên người dùng của Redis server (nếu có)
	Password string `yaml:"password"` // Mật khẩu của Redis server (nếu có)
	DB       int    `yaml:"db"`       // Số thứ tự của cơ sở dữ liệu Redis (mặc định là 0)

	PoolSize     int           `yaml:"pool_size"`      // Kích thước của pool kết nối Redis
	MinIdleConns int           `yaml:"min_idle_conns"` // Số lượng kết nối nhàn rỗi tối thiểu trong pool
	Timeout      time.Duration `yaml:"timeout"`        // Thời gian chờ kết nối Redis
	CacheTTL     time.Duration `yaml:"cache_ttl"`      // Thời gian sống của cache trong Redis
	WriteTimeout time.Duration `yaml:"write_timeout"`  // Thời gian chờ ghi dữ liệu vào Redis
	ReadTimeout  time.Duration `yaml:"read_timeout"`   // Thời gian chờ đọc dữ liệu từ Redis
	DialTimeOut  time.Duration `yaml:"dial_timeout"`   // Thời gian chờ khi kết nối đến Redis
	KeyPrefix    string        `yaml:"key_prefix"`     // Tiền tố cho các khóa

	TLS TLSConfig `yaml:"tls"` // Cấu hình TLS cho Redis
}

type GRPCconfig struct {
	Address          string        `yaml:"address"`           // Địa chỉ của gRPC server
	Timeout          time.Duration `yaml:"timeout"`           // Thời gian chờ cho các yêu cầu gRPC
	MaxRecvMsgSize   int           `yaml:"max_recv_msg_size"` // Kích thước tối đa của tin nhắn nhận được
	MaxSendMsgSize   int           `yaml:"max_send_msg_size"` // Kích thước tối đa của tin nhắn gửi đi
	KeepaliveTime    time.Duration `yaml:"keepalive_time"`    // Thời gian giữ kết nối gRPC sống
	KeepaliveTimeout time.Duration `yaml:"keepalive_timeout"` // Thời gian chờ cho keepalive
	EnableReflection bool          `yaml:"enable_reflection"` // Bật/tắt phản chiếu gRPC (true/false)
	TLS              TLSConfig     `yaml:"tls"`               // Cấu hình TLS cho gRPC
}

type KafkaConfig struct {
	Brokers  []string `yaml:"brokers"`   // Danh sách các broker của Kafka (ví dụ: ["localhost:9092", "localhost:9093"])
	ClientID string   `yaml:"client_id"` // ID của client Kafka (ví dụ: "my-kafka-client")

	UserTopic     string `yaml:"user_topic"`   // Tên topic cho người dùng (ví dụ: "user-topic")
	DLQTopic      string `yaml:"dlq_topic"`    // Tên topic cho Dead Letter Queue (ví dụ: "dlq-topic")
	WorkerGroupID string `yaml:"worker_group"` // Tên nhóm worker (ví dụ: "worker-group")

	Producer KafkaProducerConfig `yaml:"producer"` // Cấu hình producer Kafka
	Consumer KafkaConsumerConfig `yaml:"consumer"` // Cấu hình consumer Kafka

	SASL SASLConfig `yaml:"sasl"` // Cấu hình SASL cho Kafka
	TLS  TLSConfig  `yaml:"tls"`  // Cấu hình TLS cho Kafka
}

type KafkaProducerConfig struct {
	RequiredAcks int           `yaml:"required_acks"` // Số lượng xác nhận cần thiết từ broker (ví dụ: 1, 0, -1)
	MaxRetries   int           `yaml:"max_retries"`   // Số lần thử lại khi gửi tin nhắn thất bại (ví dụ: 3)
	Compression  string        `yaml:"compression"`   // Phương thức nén dữ liệu (ví dụ: "gzip", "snappy", "lz4", "none")
	BatchSize    int           `yaml:"batch_size"`    // Kích thước tối đa của batch tin nhắn (ví dụ: 100)
	BatchTimeout time.Duration `yaml:"batch_timeout"` // Thời gian chờ tối đa trước khi gửi batch tin nhắn (ví dụ: 1s)
	WriteTimeout time.Duration `yaml:"write_timeout"` // Thời gian chờ tối đa khi gửi tin nhắn (ví dụ: 5s)
}

type KafkaConsumerConfig struct {
	AutoOffsetReset   string        `yaml:"auto_offset_reset"`  // Cách xử lý offset khi consumer bắt đầu (ví dụ: "earliest", "latest")
	SessionTimeout    time.Duration `yaml:"session_timeout"`    // Thời gian chờ tối đa cho session của consumer (ví dụ: 10s)
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval"` // Khoảng thời gian gửi heartbeat từ consumer đến broker (ví dụ: 3s)
	MaxPollRecords    int           `yaml:"max_poll_records"`   // Số lượng bản ghi tối đa mà consumer có thể lấy trong một lần poll (ví dụ: 500)
	CommitInterval    time.Duration `yaml:"commit_interval"`    // Khoảng thời gian giữa các lần commit offset (ví dụ: 5s)
	MaxRetries        int           `yaml:"max_retries"`        // Số lần thử lại khi xử lý tin nhắn thất bại (ví dụ: 3)
}

type SASLConfig struct {
	Enable    bool   `yaml:"enable"`    // Bật/tắt SASL (true/false)
	Mechainsm string `yaml:"mechanism"` // Cơ chế SASL (ví dụ: "PLAIN", "SCRAM-SHA-256", "SCRAM-SHA-512")
	Username  string `yaml:"username"`  // Tên người dùng SASL
	Password  string `yaml:"password"`  // Mật khẩu SASL
}

// TLSConfig chứa thông tin cấu hình TLS cho Kafka, Redis, gRPC bao gồm các trường cần thiết để thiết lập kết nối bảo mật.
type TLSConfig struct {
	Enable             bool   `yaml:"enable"`               // Bật/tắt TLS (true/false)
	CertFile           string `yaml:"cert_file"`            // Đường dẫn đến file chứng chỉ TLS (ví dụ: "path/to/cert.pem")
	KeyFile            string `yaml:"key_file"`             // Đường dẫn đến file khóa riêng TLS (ví dụ: "path/to/key.pem")
	CAFile             string `yaml:"ca_file"`              // Đường dẫn đến file chứng chỉ CA (ví dụ: "path/to/ca.pem")
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"` // Bỏ qua xác thực chứng chỉ (true/false)
}

type LogConfig struct {
	Level  string `yaml:"level"`  // Mức độ log (ví dụ: "debug", "info", "warn", "error")
	Format string `yaml:"format"` // Định dạng log (ví dụ: "json", "text")
	Output string `yaml:"output"` // Đầu ra log (ví dụ: "stdout", "stderr", "file")
}

type AuthConfig struct {
	JWTSecret          string        `yaml:"jwt_secret"`           // Bí mật JWT để ký và xác thực token
	Issuer             string        `yaml:"issuer"`               // Tên nhà phát hành token JWT
	AccessTokenExpiry  time.Duration `yaml:"access_token_expiry"`  // Th
	RefreshTokenExpiry time.Duration `yaml:"refresh_token_expiry"` // Thời gian hết hạn của token làm mới
}

type RateLimitConfig struct {
	Enabled           bool `yaml:"enable"`              // Bật/tắt giới hạn tốc độ (true/false)
	RequestsPerSecond int  `yaml:"requests_per_second"` // Số lượng yêu cầu tối đa mỗi giây
	BurstSize         int  `yaml:"burst_size"`          // Kích thước bùng nổ (số lượng yêu cầu tối đa có thể được xử lý trong một khoảng thời gian ngắn)
}

type TelemetryConfig struct {
	MetricsEnabled bool    `yaml:"metrics_enabled"` // Bật/tắt thu thập số liệu (true/false)
	MetricsPath    string  `yaml:"metrics_path"`    // Đường dẫn để truy cập số liệu (ví dụ: "/metrics")
	TracingEnabled bool    `yaml:"tracing_enabled"` // Bật/tắt theo dõi (true/false)
	OTLEndpoint    string  `yaml:"otel_endpoint"`   // Điểm cuối của OpenTelemetry (ví dụ: "http://localhost:4317")
	SampleRate     float64 `yaml:"sample_rate"`     // Tỷ lệ mẫu cho theo dõi (ví dụ: 0.1 cho 10% các yêu cầu được theo dõi)
}
