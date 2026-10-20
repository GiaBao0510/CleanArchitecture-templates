package user

import (
	"context"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
)

func (s *Service) Update(ctx context.Context, res UpdateUserRequest) (UserResponse, error) {
	
	// Lấy thông tin người dùng từ repository
	u, err := s.repo.GetByID(ctx, res.ID)
	if err != nil {
		return UserResponse{}, err
	}

	// Kiểm tra định dạng thông tin người dùng
	if err := u.Update(res.Name, res.Email); err != nil {
		return UserResponse{}, err
	}

	// Cập nhật thông tin người dùng trong repository
	if err := s.repo.Update(ctx, u); err != nil {
		return UserResponse{}, err
	}

	// Phát sự kiện
	event, err := domain.NewEvent(domain.EventTypeUserUpdated, u)
	if err != nil {
		return UserResponse{}, err
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		return UserResponse{}, err
	}

	// Trả về kết quả
	return toOutput(u), nil
}