package config

import (
	"fmt"
	"os"
)

// Config chứa toàn bộ cấu hình ứng dụng.
// Config được đọc từ environment variables — là nguồn cấu hình duy nhất.
type Config struct {
	App      AppConfig
	Database DBConfig
}

// AppConfig chứa cấu hình ứng dụng.
type AppConfig struct {
	Port     string
	LogLevel string
}

// DBConfig chứa cấu hình kết nối database.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN trả về connection string cho PostgreSQL.
func (db DBConfig) DSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		db.User, db.Password, db.Host, db.Port, db.Name, db.SSLMode,
	)
}

// Load đọc config từ environment variables.
func Load() *Config {
	return &Config{
		App: AppConfig{
			Port:     getEnv("APP_PORT", "8080"),
			LogLevel: getEnv("LOG_LEVEL", "debug"),
		},
		Database: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "clean_architecture"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
	}
}

// getEnv đọc biến môi trường, nếu không có thì dùng giá trị mặc định.
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
