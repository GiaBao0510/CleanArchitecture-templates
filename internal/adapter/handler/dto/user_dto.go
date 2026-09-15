package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/GiaBao0510/project-templates/internal/domain/entity"
)

// --- Request DTOs ---

// CreateUserRequest là dữ liệu đầu vào khi tạo user.
type CreateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email"     binding:"required,email"`
}

// UpdateUserRequest là dữ liệu đầu vào khi cập nhật user.
type UpdateUserRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
	Email    string `json:"email"     binding:"required,email"`
}

// --- Response DTOs ---

// UserResponse là dữ liệu trả về cho client.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToUserResponse chuyển entity sang response DTO.
func ToUserResponse(u *entity.User) *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		FullName:  u.FullName,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// ToUserResponseList chuyển danh sách entity sang danh sách response DTO.
func ToUserResponseList(users []*entity.User) []*UserResponse {
	result := make([]*UserResponse, len(users))
	for i, u := range users {
		result[i] = ToUserResponse(u)
	}
	return result
}
