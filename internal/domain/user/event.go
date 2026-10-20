package user

import (
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeUserCreated EventType = "user.created"
	EventTypeUserUpdated EventType = "user.updated"
	EventTypeUserDeleted EventType = "user.deleted"
)

type Event struct {
	ID         uuid.UUID
	Type       EventType
	AggegateID uuid.UUID
	OccurredAt time.Time
	Payload    any
}

type UserSnapshot struct {
	ID    uuid.UUID
	Name  string
	Email string
}

func NewEvent(t EventType, u *User) (Event, error) {
	return Event{
		ID:         uuid.New(),
		Type:       t,
		AggegateID: u.ID,
		OccurredAt: time.Now().UTC(),
		Payload:    UserSnapshot{
			ID: u.ID, 
			Name: u.Name, 
			Email: u.Email,
		},
	}, nil
}
