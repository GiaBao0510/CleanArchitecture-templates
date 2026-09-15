package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/GiaBao0510/project-templates/internal/domain/entity"
	"github.com/GiaBao0510/project-templates/internal/domain/repository"
	"github.com/GiaBao0510/project-templates/pkg/apperror"
)

// userRepository implement interface repository.UserRepository.
// Đây là adapter — kết nối domain (bên trong) với database (bên ngoài).
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository tạo instance mới (trả về interface, không trả về struct).
func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	query := `INSERT INTO users (id, full_name, email, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.Exec(ctx, query,
		user.ID, user.FullName, user.Email, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return apperror.NewInternal("Không thể tạo user", err)
	}
	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	query := `SELECT id, full_name, email, created_at, updated_at
	          FROM users WHERE id = $1`

	var user entity.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.FullName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperror.NewNotFound("User không tồn tại")
		}
		return nil, apperror.NewInternal("Không thể lấy user", err)
	}
	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	query := `SELECT id, full_name, email, created_at, updated_at
	          FROM users ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, apperror.NewInternal("Không thể lấy danh sách users", err)
	}
	defer rows.Close()

	var users []*entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(
			&user.ID, &user.FullName, &user.Email, &user.CreatedAt, &user.UpdatedAt,
		); err != nil {
			return nil, apperror.NewInternal("Không thể đọc dữ liệu user", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	query := `UPDATE users SET full_name = $1, email = $2, updated_at = $3
	          WHERE id = $4`

	result, err := r.db.Exec(ctx, query,
		user.FullName, user.Email, time.Now(), user.ID,
	)
	if err != nil {
		return apperror.NewInternal("Không thể cập nhật user", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("User không tồn tại")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return apperror.NewInternal("Không thể xóa user", err)
	}
	if result.RowsAffected() == 0 {
		return apperror.NewNotFound("User không tồn tại")
	}
	return nil
}
