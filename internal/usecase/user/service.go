package user

import "context"

// Service is the use case layer for user management.
type Service struct {
	repo UserRepository
	publisher EventPublisher
}

// NewService creates a new instance of the Service with the provided UserRepository and EventPublisher.
func NewService(repo UserRepository, publisher EventPublisher) *Service {
	return &Service{
		repo: repo,
		publisher: publisher,
	}
}

// publish is a helper method to publish events using the EventPublisher.
func(s *Service) publish(
	ctx context.Context,
	event interface{},
) error {
	panic("implement me")
}