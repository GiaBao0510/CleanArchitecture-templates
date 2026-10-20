package user

import (
	"context"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
)

func (s *Service) Create(ctx context.Context, user *domain.User) ( UserResponse, error) {
	
	// Khởi tạo giá trị
	u, err := domain.NewUser(user.Name, user.Email)
	if err != nil {
		return UserResponse{}, err
	}

	// Kiểm tra trùng lặp
	existingUser, err := s.repo.GetByEmail(ctx, u.Email)
	if err != nil {
		return UserResponse{}, err
	}
	if existingUser != nil {
		return UserResponse{}, domain.ErrEmailIsAlreadyUsed
	}

	// Thực hiện thêm
	if err := s.repo.Create(ctx, u); err != nil {
		return UserResponse{}, err
	}

	// phát sự kiện
	event, err := domain.NewEvent(domain.EventTypeUserCreated, u)
	if err != nil {
		return UserResponse{}, err
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		return UserResponse{}, err
	}

	// Trả về kết quả
	return toOutput(u), nil
}