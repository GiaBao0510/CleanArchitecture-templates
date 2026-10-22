package mysql

import (
	"context"
	"database/sql"
	"errors"

	queries "github.com/GiaBao0510/CleanArchitecture-templates/gen/sqlc/mysql"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/GiaBao0510/CleanArchitecture-templates/pkg/utils"
	"github.com/google/uuid"
)

type Repository struct {
	q *queries.Queries
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{q: queries.New(db)}
}

func (r *Repository) Create(
	ctx context.Context, u *domain.User,
) error {
	return mapError(r.q.CreateUser(ctx, queries.CreateUserParams{
		ID: uuid.NewString(), Name: u.Name, Email: u.Email,
		CreatedAt: u.CreatedAt, UpdatedAt: u.UpdatedAt,
	}))
}

func (r *Repository) GetByID(
	ctx context.Context, id uuid.UUID,
) (*domain.User, error) {

	ID := id.String()

	row, err := r.q.GetUserByID(ctx, ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	userID, err := utils.StringToUUID(row.ID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID: userID , Name: row.Name, 
		Email: row.Email,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *Repository) GetByEmail(
	ctx context.Context, email string,
) (*domain.User, error) {
	row, err := r.q.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	userID, err := utils.StringToUUID(row.ID)
	if err != nil {
		return nil, err
	}
	return &domain.User{
		ID: userID, Name: row.Name, Email: row.Email,
		CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
	}, nil
}

func (r *Repository) List(
	ctx context.Context, limit, offset int32,
) ([]*domain.User, error) {
	rows, err := r.q.ListUsers(ctx, queries.ListUsersParams{
		Limit: limit, Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	users := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		userID, err := utils.StringToUUID(row.ID)
		if err != nil {
			return nil, err
		}
		users = append(users, &domain.User{
			ID: userID,
			Name: row.Name, Email: row.Email,
			CreatedAt: row.CreatedAt, UpdatedAt: row.UpdatedAt,
		})
	}
	return users, nil
}

func (r *Repository) Update(
	ctx context.Context, u *domain.User,
) error {

	userID := utils.UUIDtoString(u.ID)

	n, err := r.q.UpdateUser(ctx, queries.UpdateUserParams{
		ID: userID, Name: u.Name, Email: u.Email,
		UpdatedAt: u.UpdatedAt,
	})
	if err != nil {
		return mapError(err)
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *Repository) Delete(
	ctx context.Context, id uuid.UUID,
) error {

	userID := utils.UUIDtoString(id)

	n, err := r.q.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}