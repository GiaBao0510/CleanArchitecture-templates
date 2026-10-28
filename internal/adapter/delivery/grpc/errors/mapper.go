package errors

import (
	"context"
	"errors"

	"github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/shared/apperror"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Map chuyển đổi Application error sang gRPC status
func Map(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, context.Canceled):
		return status.Error(
			codes.Canceled,
			"request canceled",
		)

	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(
			codes.DeadlineExceeded,
			"deadline exceeded",
		)
	}

	switch apperror.KindOf(err) {
	case apperror.Invalid:
		return status.Error(
			codes.InvalidArgument,
			err.Error(),
		)

	case apperror.NotFound:
		return status.Error(
			codes.NotFound,
			err.Error(),
		)

	case apperror.Conflict:
		return status.Error(
			codes.AlreadyExists,
			err.Error(),
		)

	default:
		//Trả về lỗi liên quan chung chung. Không hiển thị chi tiết
		return status.Error(
			codes.Internal,
			"internal server error",
		)
	}
}
