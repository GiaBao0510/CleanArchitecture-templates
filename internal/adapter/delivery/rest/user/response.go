package user

import (
	"time"

	usecase "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/user"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func toResponse(u usecase.UserResponse) UserResponse {
	return UserResponse{
		ID: u.ID,
		Name: u.Name,
		Email: u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}