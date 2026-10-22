package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/shared" // đổi thành import path thực tế của bạn
)

// mapError chuyển lỗi của PostgreSQL/pgx thành lỗi domain.
//
// Quy ước:
//   - Lỗi domain được wrap bằng %w nên errors.Is(err, shared.ErrXxx) hoạt động.
//   - Lỗi gốc cũng được giữ lại (Go 1.20+ hỗ trợ nhiều %w) để phục vụ log/debug.
//   - Lỗi không nhận diện được thì trả nguyên, không che giấu.
func mapError(err error) error {
	// 1. Không có lỗi thì trả về nil.
	if err == nil {
		return nil
	}

	// 2. Lỗi không liên quan dữ liệu: client hủy request hoặc hết thời gian chờ.
	//    Trả nguyên để tầng trên nhận biết (ví dụ trả HTTP 408/499).
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// 3. Không tìm thấy bản ghi (QueryRow(...).Scan(...) không có dòng nào).
	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %w", shared.ErrNotFound, err)
	}

	// 4. Lỗi do PostgreSQL trả về (có mã SQLSTATE).
	//    Chỉ gọi errors.As một lần rồi switch theo mã lỗi.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {

		// 23505: vi phạm unique (duplicate key), ví dụ email đã tồn tại.
		case pgerrcode.UniqueViolation:
			return fmt.Errorf("%w (constraint %q): %w",
				shared.ErrAlreadyExists, pgErr.ConstraintName, err)

		// 23503: vi phạm khóa ngoại (tham chiếu tới bản ghi không tồn tại,
		// hoặc xóa bản ghi đang được tham chiếu).
		case pgerrcode.ForeignKeyViolation:
			return fmt.Errorf("%w (constraint %q): %w",
				shared.ErrInvalidReference, pgErr.ConstraintName, err)

		// 23502: cột NOT NULL bị gán null.
		// 23514: vi phạm CHECK constraint.
		// 22P02: sai định dạng dữ liệu (ví dụ UUID/số không hợp lệ).
		case pgerrcode.NotNullViolation,
			pgerrcode.CheckViolation,
			pgerrcode.InvalidTextRepresentation:
			return fmt.Errorf("%w: %w", shared.ErrInvalidData, err)

		// 40001: serialization failure, 40P01: deadlock.
		// Giao dịch bị hủy, caller có thể thử lại.
		case pgerrcode.SerializationFailure, pgerrcode.DeadlockDetected:
			return fmt.Errorf("%w: %w", shared.ErrConflict, err)
		}
	}

	// 5. Các lỗi còn lại (mất kết nối, lỗi cú pháp SQL, ...): trả nguyên.
	return err
}