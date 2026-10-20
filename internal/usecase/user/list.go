package user

import (
	"context"
)

func(s *Service) List(ctx context.Context, res ListUserRequest) ([]UserResponse, error) {

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

	users, err := s.repo.List(ctx, res.Limit, res.Offset)
	if err != nil {
		return nil, err
	}

	result := make([]UserResponse, 0, len(users))
	for _, u := range users {
		result = append(result, toOutput(u))
	}
	return result, nil
}