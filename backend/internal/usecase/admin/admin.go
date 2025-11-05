package admin

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/google/uuid"
)

func (a *AdminService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return a.repo.GetAllUsers(ctx)
}

func (a *AdminService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return a.repo.DeleteUser(ctx, id)
}

func (a *AdminService) GetAllRecipes(ctx context.Context) ([]*domain.Recipe, error) {
	return a.repo.GetAllRecipes(ctx)
}

func (a *AdminService) DeleteRecipe(ctx context.Context, id uuid.UUID) error {
	return a.repo.DeleteRecipe(ctx, id)
}
