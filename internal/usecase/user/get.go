package user

import (
	"context"
	"github.com/google/uuid"
)

func (s *Service) Get(ctx context.Context, id uuid.UUID) (UserResponse, error) {
	u, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return UserResponse{},mapError(err)
	}

	return toOutput(u), nil
}