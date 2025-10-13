package domain

import "errors"

var (
	// User errors 
	ErrCredentialsRequired = errors.New("this field is required")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidUsername = errors.New("username must be at least 5 characters and not more than 64 characters")
)
