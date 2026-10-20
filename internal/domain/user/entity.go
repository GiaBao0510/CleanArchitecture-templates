package user

import (
	"crypto/rand"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ID string

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

// NewID generates a new unique ID for a user.
func NewID() (ID, error) {
	var b [16]byte

	// If there is an error generating random bytes, return an empty ID and the error.
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}

	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80

	return ID(fmt.Sprintf(
		"%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16],
	)), nil
}

//
func ParseID(value string) (ID, error) {
	if value == ""  || len(value) != 36  {
		return "", ErrInvalidID
	}
	for i, c := range value {
		if (i == 8 || i == 13 || i == 18 || i == 23) {
			if c != '-' {
				return "", ErrInvalidID
			}
			continue
		}
		if !strings.ContainsRune("0123456789abcdefgh", c) {
			return "", ErrInvalidID
		}
	}

	return ID(strings.ToLower(value)), nil
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