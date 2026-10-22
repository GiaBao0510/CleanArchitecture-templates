package user

import (
	usecase "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/user"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (r *CreateRequest) Input() usecase.CreateUserRequest {
	return usecase.CreateUserRequest{
		Name:  r.Name,
		Email: r.Email,
	}
}

func (r *UpdateRequest) Input(id uuid.UUID) usecase.UpdateUserRequest {
	return usecase.UpdateUserRequest {
		ID:    id,
		Name:  r.Name,
		Email: r.Email,
	}
} 