package shared

import "errors"

var (
	ErrNotFound           = errors.New("not found")           // Không tìm thấy bản ghi
	ErrAlreadyExists      = errors.New("already exists")      // Trùng dữ liệu (unique)
	ErrInvalidReference   = errors.New("invalid reference")   // Khóa ngoại không tồn tại / đang bị tham chiếu
	ErrInvalidData        = errors.New("invalid data")        // Vi phạm not null, check, sai định dạng
	ErrConflict           = errors.New("conflict, retryable") // Xung đột giao dịch, có thể thử lại
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrResourceInUse      = errors.New("resource is in use")
	ErrDBConflict         = errors.New("database conflict, please retry")
	ErrDBUnavailable      = errors.New("database unavailable")
	ErrInvalidID          = errors.New("Invalid ID")
)
