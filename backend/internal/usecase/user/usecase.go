package user

import (
	"context"
	"github.com/goawwer/yamyard/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
} 
