package apperror

import "net/http"

// AppError là custom error type cho ứng dụng.
// Giúp phân biệt các loại lỗi (Not Found, Bad Request, Internal...)
// và map sang HTTP status code tương ứng.
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// --- Constructors ---

// NewNotFound tạo lỗi 404.
func NewNotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

// NewBadRequest tạo lỗi 400.
func NewBadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// NewInternal tạo lỗi 500.
func NewInternal(message string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: err}
}

// NewConflict tạo lỗi 409.
func NewConflict(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}
