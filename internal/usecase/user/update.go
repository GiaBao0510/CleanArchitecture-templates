package user

import (
	"context"
	apperror "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/shared/apperror"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
)

func (s *Service) Update(ctx context.Context, res UpdateUserRequest) (UserResponse, error) {
	
	// Lấy thông tin người dùng từ repository
	u, err := s.repo.GetByID(ctx, res.ID)
	if err != nil {
		return UserResponse{}, apperror.New(apperror.NotFound, "user not found", err)
	}

	// Kiểm tra định dạng thông tin người dùng
	if err := u.Update(res.Name, res.Email); err != nil {
		return UserResponse{}, apperror.New(apperror.BadRequest, "invalid user data", err)
	}

	// Cập nhật thông tin người dùng trong repository
	if err := s.repo.Update(ctx, u); err != nil {
		return UserResponse{}, apperror.New(apperror.Internal, "failed to update user", err)
	}

	// Phát sự kiện
	event, err := domain.NewEvent(domain.EventTypeUserUpdated, u)
	if err != nil {
		return UserResponse{}, apperror.New(apperror.Internal, "failed to create user updated event", err)
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		return UserResponse{}, apperror.New(apperror.Internal, "failed to publish user updated event", err)
	}

	// Trả về kết quả
	return toOutput(u), nil
}