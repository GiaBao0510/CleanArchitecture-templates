package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/GiaBao0510/project-templates/internal/domain/entity"
)

// UserRepository định nghĩa interface (port) cho tầng domain.
// Tầng domain CHỈ khai báo interface — KHÔNG implement.
// Việc implement sẽ do tầng adapter đảm nhận (Dependency Inversion).
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetAll(ctx context.Context) ([]*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}
