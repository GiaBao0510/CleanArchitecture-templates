package user

import (
	"context"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/google/uuid"
)

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	
	// Lấy thông tin người dùng từ repository
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Xóa người dùng trong repository
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	// phát sự kiện
	event, err := domain.NewEvent(domain.EventTypeUserDeleted, u)
	if err != nil {
		return err
	}
	return s.publisher.Publish(ctx, event)
}