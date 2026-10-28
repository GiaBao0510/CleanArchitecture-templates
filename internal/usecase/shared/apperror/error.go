package apperror

import "errors"

// Kind mô tả ý nghĩa của từng cấp ứng dụng
type Kind string

const (
	Invalid Kind = "INVALID"
	NotFound Kind = "NOT_FOUND"
	Conflict Kind = "CONFLICT"
	Internal Kind = "INTERNAL"
	BadRequest Kind = "BAD_REQUEST"
)

// Lỗi ứng dụng được dùng chung cho các usecase, service, handler
type AppError struct{
	Kind Kind
	Message string
	Cause error
}

func (e *AppError) Error() string {
	return e.Message
}

// Unwrap giữ lại lỗi gốc để errors.Is / errors.As hoạt động.
func(e *AppError) Unwrap() error {
	return e.Cause
}

// New tạo ra lỗi ứng dụng mới với kind, message và cause
func New(kind Kind, message string, cause error) error {
	return &AppError{
		Kind: kind,
		Message: message,
		Cause: cause,
	}
}

// KindOf cho phép delivery adapter phân loại lỗi
// mà không cần import từng domain cụ thể.
func KindOf(err error) Kind {
	var appErr *AppError
	if errors.As(err, &appErr){
		return appErr.Kind
	}

	return Internal
}