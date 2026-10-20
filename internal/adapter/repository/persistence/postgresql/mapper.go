package postgresql

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func toPostgresUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id, 
		Valid: true,
	}
}

func toUUID(pgUUID pgtype.UUID) (uuid.UUID, error) {
	if !pgUUID.Valid {
		return uuid.Nil, fmt.Errorf("invalid UUID")
	}
	return uuid.UUID(pgUUID.Bytes), nil
}