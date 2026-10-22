package utils

import "github.com/google/uuid"

func UUIDtoString(id uuid.UUID) string {
	return id.String()
}

func StringToUUID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}