package user

import (
	"time"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Name     string 
	Email    string
}

type UpdateUserRequest struct {
	ID uuid.UUID
	Name  string 
	Email string 
}

type ListUserRequest struct {
	Limit  int32
	Offset int32
}

type UserResponse struct {
	ID		uuid.UUID `json:"id"`
	Name	string    `json:"name"`
	Email	string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toOutput(u *domain.User) UserResponse {
	return UserResponse{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}