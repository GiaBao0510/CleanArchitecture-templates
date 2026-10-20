package user

import "errors"

var (
	ErrNameIsRequired  = errors.New("Name is required")
	ErrEmailIsInvalid  = errors.New("Invalid email address")
	ErrEmailIsAlreadyUsed = errors.New("Email is already used")
	ErrEmailIsRequired = errors.New("Email is required")
	ErrInvalidName     = errors.New("Invalid name")
	ErrNameIsTooLong   = errors.New("Name is too long")
	ErrInvalidID         = errors.New("Invalid ID format")
)
