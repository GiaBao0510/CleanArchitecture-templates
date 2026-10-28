package errors

import (
	"context"
	"errors"
	"net/http"

	appError "github.com/GiaBao0510/CleanArchitecture-templates/internal/usecase/shared/apperror"
)

// Chuyển lỗi ứng dung sang mã trạng thái HTTP và thông điệp lỗi tương ứng
func Status(err error) int {
	
	// Xử lý các lỗi ngữ cảnh trước tiên
	switch {
	case errors.Is(err, context.Canceled):
		// 499 là convention của một số proxy, không phải status chuẩn của net/http.
		return http.StatusRequestTimeout // Client Closed Request
	case errors.Is(err, context.DeadlineExceeded):
		return http.StatusGatewayTimeout // Request Timeout
	}

	// Nếu lỗi là lỗi ứng dụng, chuyển kind sang status code
	switch appError.KindOf(err) {
	case appError.Invalid:
		return http.StatusUnprocessableEntity
	case appError.NotFound:
		return http.StatusNotFound
	case appError.Conflict:
		return http.StatusConflict
	case appError.BadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}