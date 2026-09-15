package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/GiaBao0510/project-templates/internal/domain/entity"
	"github.com/GiaBao0510/project-templates/internal/domain/repository"
)

// UserUsecase chứa business logic của ứng dụng.
// Use case phụ thuộc vào repository INTERFACE (không phụ thuộc implementation cụ thể).
// Đây là nơi thực thi các quy tắc nghiệp vụ (business rules).
type UserUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase khởi tạo use case với dependency injection.
func NewUserUsecase(userRepo repository.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

// Create tạo user mới với ID và timestamps tự động.
func (uc *UserUsecase) Create(ctx context.Context, fullName, email string) (*entity.User, error) {
	now := time.Now()
	user := &entity.User{
		ID:        uuid.New(),
		FullName:  fullName,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID lấy user theo ID.
func (uc *UserUsecase) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

// GetAll lấy danh sách tất cả users.
func (uc *UserUsecase) GetAll(ctx context.Context) ([]*entity.User, error) {
	return uc.userRepo.GetAll(ctx)
}

// Update cập nhật thông tin user.
func (uc *UserUsecase) Update(ctx context.Context, id uuid.UUID, fullName, email string) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.FullName = fullName
	user.Email = email
	user.UpdatedAt = time.Now()

	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

// Delete xóa user theo ID.
func (uc *UserUsecase) Delete(ctx context.Context, id uuid.UUID) error {
	return uc.userRepo.Delete(ctx, id)
}
