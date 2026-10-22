package postgresql

import (
	"context"

	queries "github.com/GiaBao0510/CleanArchitecture-templates/gen/sqlc/postgres"
	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	"github.com/GiaBao0510/CleanArchitecture-templates/pkg/utils"
	"github.com/google/uuid"


	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q *queries.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q:queries.New(pool)}
}

func (r *Repository) Create(ctx context.Context, u *domain.User) error {
	
	// Generate a new UUID for the user
	uuid := uuid.NewString()
	id, err := StringToPostgresUUID(uuid)
	
	if err != nil {
		return err
	}
	return mapError(r.q.CreateUser(ctx, queries.CreateUserParams{
		ID: id, Name: u.Name, Email: u.Email,
		CreatedAt: toPgtypTimestamptz(u.CreatedAt), 
		UpdatedAt: toPgtypTimestamptz(u.UpdatedAt),
	}))

}
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	
	// Convert the UUID to a PostgreSQL-compatible format
	pgID := toPostgresUUID(id)

	// Fetch the user from the database using the generated SQL query
	row, err := r.q.GetUserByID(ctx, pgID)
	if err != nil {
		return nil, mapError(err)
	}

	// Convert the database row to a domain.User object
	userID, err := utils.StringToUUID(row.ID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID: userID,
		Name: row.Name,
		Email: row.Email,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*domain.User, error){
	row, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, mapError(err)
	}

	userID, err := utils.StringToUUID(row.ID)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID: userID,
		Name: row.Name,
		Email: row.Email,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}

func (r *Repository) List(ctx context.Context, limit, offset int32) ([]*domain.User, error){
	rows, err := r.q.ListUsers(ctx, queries.ListUsersParams{
		Limit: limit,
		Offset: offset,
	})
	if err != nil {
		return nil, mapError(err)
	}

	users := make([]*domain.User, len(rows))
	for _, row := range rows {
		userID, _ := utils.StringToUUID(row.ID)
		users = append(users, &domain.User{
			ID: userID,
			Name: row.Name,
			Email: row.Email,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}
	return users, nil
}

func (r *Repository) Update(ctx context.Context, u *domain.User) error{ 
	id := toPostgresUUID(u.ID)
	
	n, err := r.q.UpdateUser(ctx, queries.UpdateUserParams{
		ID: id,
		Name: u.Name,
		Email: u.Email,
		UpdatedAt: toPgtypTimestamptz(u.UpdatedAt),
	})
	if	 err != nil {
		return mapError(err)
	}

	if n == 0 {
		return domain.ErrNotFound
	}

	// Update the user in the database using the generated SQL query
	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error{
	pgID := toPostgresUUID(id)
	
	n, err := r.q.DeleteUser(ctx, pgID)
	if err != nil {
		return mapError(err)
	}

	if n == 0 {
		return domain.ErrNotFound
	}

	return nil
}
