package user

import (
	"context"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
)

func(s *Service) List(
	ctx context.Context, 
	res ListUserRequest,
) ([]*domain.User, error) {

	// Check limit and offset
	if res.Limit <= 0 {
		res.Limit = 10
	}
	if res.Limit > 100 {
		res.Limit = 100
	}
	if res.Offset < 0 {
		res.Offset = 0
	}

	return s.repo.List(ctx, res.Limit, res.Offset)
}