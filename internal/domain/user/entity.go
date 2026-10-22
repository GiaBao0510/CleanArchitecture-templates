package user

import (
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)



type User struct {
	ID uuid.UUID
	Name string
	Email string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(name, email string) (*User, error) {

	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if err := validation(name, email); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID: uuid.New(),
		Name: name,
		Email: email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (u *User) Update(name, email string) error {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)

	if err := validation(name, email); err != nil {
		return err
	}

	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now().UTC()
	return nil
}

func validation (name, email string) error {
	if name == "" {
		return ErrNameIsRequired
	}

	if email == "" {
		return ErrEmailIsRequired
	}

	if len(name) > 100 {
		return ErrNameIsTooLong
	}
	
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || len(email) > 255 || strings.ContainsAny(email, " \t\n") {
		return ErrEmailIsInvalid
	}

	return nil
}