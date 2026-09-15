package entity

import (
	"time"

	"github.com/google/uuid"
)

// User là entity cốt lõi trong domain layer.
// Entity chỉ chứa business data, không phụ thuộc bất kỳ framework hay thư viện bên ngoài nào
// (ngoại trừ uuid dùng cho kiểu dữ liệu ID).
type User struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
