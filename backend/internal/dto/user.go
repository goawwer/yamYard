package dto

import "github.com/google/uuid"

type SignUpUserInput struct {
	Email string
	Username string
	Password string
}

type SignUpUserOutput struct {
	ID uuid.UUID
}
