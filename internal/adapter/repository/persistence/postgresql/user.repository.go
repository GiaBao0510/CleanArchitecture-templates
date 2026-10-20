package postgresql

import (
		"context"
	"errors"
	"fmt"

	domain "github.com/GiaBao0510/CleanArchitecture-templates/internal/domain/user"
	queries "github.com/GiaBao0510/CleanArchitecture-templates/gen/sqlc/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	q *queries.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q:queries.New(pool)}
}


