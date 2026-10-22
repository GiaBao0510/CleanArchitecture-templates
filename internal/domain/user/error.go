package user

import "errors"

var (
	ErrNameIsRequired     = errors.New("Name is required")
	ErrEmailIsInvalid     = errors.New("Invalid email address")
	ErrEmailIsAlreadyUsed = errors.New("Email is already used")
	ErrEmailIsRequired    = errors.New("Email is required")
	ErrInvalidName        = errors.New("Invalid name")
	ErrInvalidEmail       = errors.New("Invalid email format")
	ErrNameIsTooLong      = errors.New("Name is too long")
	ErrInvalidID          = errors.New("Invalid ID format")
	ErrNotFound           = errors.New("User not found")
	ErrDuplicate          = errors.New("Duplicate entry")
)
