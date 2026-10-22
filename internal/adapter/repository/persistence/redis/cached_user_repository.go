package redis

import (
	"context"
	"errors"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	usecase "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/user"
	"github.com/google/uuid"
)

type CachedUserRepository struct {
	next usecase.UserRepository
	cache *UserCache
}

func NewCachedUserRepository(
	next usecase.UserRepository,
	cache *UserCache,
) *CachedUserRepository {
	return &CachedUserRepository{next: next, cache: cache}
}

// Create creates a new user in the underlying repository and does not cache it. Caching is typically done on retrieval.
func (r *CachedUserRepository) Create(ctx context.Context, u *domain.User) error {
	return r.next.Create(ctx, u)
}

// GetByID retrieves a user by ID. It first checks the cache; if the user is not found in the cache, it fetches from the underlying repository and caches the result.
func (r *CachedUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	
	userID := id.String()

	u, err := r.cache.Get(ctx, userID)
	if err == nil {
		return u, nil
	}

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Nếu trường hợp lấy không thấy thì tạo mới rồi sau đó cache lại 
	u, err = r.next.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = r.cache.Set(ctx, u)
	return u, nil

}
func (r *CachedUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.next.GetByEmail(ctx, email)
}

func (r *CachedUserRepository) List(ctx context.Context, limit, offset int32) ([]*domain.User, error) {
	return r.next.List(ctx, limit, offset)
}
func (r *CachedUserRepository) Update(ctx context.Context, u *domain.User) error {
	if err := r.next.Update(ctx, u); err != nil {
		return err
	}
	return r.cache.Delete(ctx, u.ID.String())
} 

func (r *CachedUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.next.Delete(ctx, id); err != nil {
		return err
	}

	err := r.cache.Delete(ctx, id.String())
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}

	return err
}

