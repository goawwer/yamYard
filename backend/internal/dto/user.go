package dto

import "github.com/google/uuid"

type SignUpUserInput struct {
	Email    string
	Username string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	ID uuid.UUID
}
