package admin

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/google/uuid"
)

type AdminRepository interface {
	GetAllUsers(ctx context.Context) ([]*domain.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetAllRecipes(ctx context.Context) ([]*domain.Recipe, error)
	DeleteRecipe(ctx context.Context, id uuid.UUID) error
}

type AdminService struct {
	repo AdminRepository
}

func NewAdminService(repo AdminRepository) *AdminService {
	return &AdminService{
		repo: repo,
	}
}
