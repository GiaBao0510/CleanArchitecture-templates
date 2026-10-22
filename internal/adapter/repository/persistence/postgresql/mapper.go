package postgresql

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toPostgresUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id, 
		Valid: true,
	}
}

func StringToPostgresUUID(id string) (pgtype.UUID, error) {
	u, err := uuid.Parse(id)
	if err != nil {
		return pgtype.UUID{}, fmt.Errorf("invalid UUID string: %w", err)
	}
	return toPostgresUUID(u), nil
}

func toPgtypTimestamptz(t time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time: t,
		Valid: true,
	}
}

func toUUID(pgUUID pgtype.UUID) (uuid.UUID, error) {
	if !pgUUID.Valid {
		return uuid.Nil, fmt.Errorf("invalid UUID")
	}
	return uuid.UUID(pgUUID.Bytes), nil
}