package user

import (
	"errors"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/shared/apperror"
)

// mapError chuyển lỗi đặc thù của User domain
// thành loại lỗi chung mà các delivery adapters hiểu.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, domain.ErrInvalidID),
		errors.Is(err, domain.ErrInvalidName),
		errors.Is(err, domain.ErrInvalidEmail):
		return apperror.New(
			apperror.Invalid,
			err.Error(),
			err,
		)

	case errors.Is(err, domain.ErrNotFound):
		return apperror.New(
			apperror.NotFound,
			"user not found",
			err,
		)

	case errors.Is(err, domain.ErrEmailIsAlreadyUsed):
		return apperror.New(
			apperror.Conflict,
			"email already exists",
			err,
		)

	default:
		return apperror.New(
			apperror.Internal,
			"internal server error",
			err,
		)
	}
}