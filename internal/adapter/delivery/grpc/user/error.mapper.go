package user

import (
	"context"
	"errors"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrDuplicate),
		errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrInvalidName):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())

	case errors.Is(err, domain.ErrEmailIsAlreadyUsed):
		return status.Error(codes.AlreadyExists, err.Error())

	case  errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, err.Error())
	
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, err.Error())
	
	default:
		return status.Error(codes.Internal, err.Error())	
	}
}
