package user

import (
	"context"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/google/uuid"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.User, error)
	Update(ctx context.Context, u *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type EventPublisher interface {
	Publish(ctx context.Context, event domain.Event) error
}