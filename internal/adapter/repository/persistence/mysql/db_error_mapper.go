package mysql

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/shared"
	mysqldriver "github.com/go-sql-driver/mysql"
)

// Mã lỗi MySQL: https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html
const (
	mysqlErrDupEntry          uint16 = 1062 // Duplicate entry for key
	mysqlErrNoReferencedRow   uint16 = 1452 // FK: bản ghi cha không tồn tại
	mysqlErrRowIsReferenced   uint16 = 1451 // FK: bản ghi đang được tham chiếu
	mysqlErrBadNull           uint16 = 1048 // Column cannot be null
	mysqlErrNoDefaultForField uint16 = 1364 // Field doesn't have a default value
	mysqlErrDataTooLong       uint16 = 1406 // Data too long for column
	mysqlErrTruncatedWrongVal uint16 = 1366 // Incorrect value for column
	mysqlErrLockWaitTimeout   uint16 = 1205 // Lock wait timeout exceeded
	mysqlErrLockDeadlock      uint16 = 1213 // Deadlock found
	mysqlErrTooManyConns      uint16 = 1040 // Too many connections
	mysqlErrServerGone        uint16 = 2006 // MySQL server has gone away
	mysqlErrServerLost        uint16 = 2013 // Lost connection during query
)

// mapError chuyển lỗi DB (MySQL) thành lỗi domain.
//
// Quy tắc:
//   - nil                      -> nil
//   - lỗi KHÔNG liên quan DB   -> trả về nguyên bản (ví dụ context.Canceled, lỗi nghiệp vụ...)
//   - lỗi liên quan DB         -> map sang lỗi domain, vẫn giữ lỗi gốc trong chuỗi wrap
//     để log/debug được và vẫn dùng errors.Is / errors.As được.
func mapError(err error) error {
	if err == nil {
		return nil
	}

	// Lỗi context: không phải lỗi DB, trả về nguyên bản.
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	// Lỗi từ package database/sql.
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return domain.ErrNotFound
	case errors.Is(err, sql.ErrConnDone),
		errors.Is(err, sql.ErrTxDone),
		errors.Is(err, driver.ErrBadConn):
		return fmt.Errorf("%w: %w", domain.ErrDBUnavailable, err)
	}

	// Lỗi do MySQL server trả về.
	var mysqlErr *mysqldriver.MySQLError
	if !errors.As(err, &mysqlErr) {
		// Không phải lỗi MySQL -> không xử lý, trả về như cũ.
		return err
	}

	switch mysqlErr.Number {
	case mysqlErrDupEntry:
		return mapDuplicateEntry(mysqlErr, err)

	case mysqlErrNoReferencedRow:
		return fmt.Errorf("%w: %w", domain.ErrInvalidReference, err)

	case mysqlErrRowIsReferenced:
		return fmt.Errorf("%w: %w", domain.ErrResourceInUse, err)

	case mysqlErrBadNull,
		mysqlErrNoDefaultForField,
		mysqlErrDataTooLong,
		mysqlErrTruncatedWrongVal:
		return fmt.Errorf("%w: %w", domain.ErrInvalidData, err)

	case mysqlErrLockWaitTimeout, mysqlErrLockDeadlock:
		// Có thể retry ở tầng trên.
		return fmt.Errorf("%w: %w", domain.ErrDBConflict, err)

	case mysqlErrTooManyConns, mysqlErrServerGone, mysqlErrServerLost:
		return fmt.Errorf("%w: %w", domain.ErrDBUnavailable, err)

	default:
		// Lỗi MySQL chưa được map: bọc thêm ngữ cảnh, giữ nguyên lỗi gốc.
		return fmt.Errorf("mysql error %d: %w", mysqlErr.Number, err)
	}
}

// mapDuplicateEntry phân biệt vi phạm unique theo từng cột/key.
func mapDuplicateEntry(mysqlErr *mysqldriver.MySQLError, original error) error {
	msg := strings.ToLower(mysqlErr.Message)

	switch {
	case strings.Contains(msg, "email"):
		return fmt.Errorf("%w: %w", domain.ErrEmailAlreadyExists, original)
	default:
		return fmt.Errorf("%w: %w", domain.ErrAlreadyExists, original)
	}
}