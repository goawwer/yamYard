package recipe

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

func (r RecipeService) GetAll(ctx context.Context, input helpers.FilterAndSortingParameters) ([]*domain.Recipe, error) {
	return r.repo.GetAll(ctx, input)
}

func (r RecipeService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	return r.repo.GetByID(ctx, id)
}

func (r RecipeService) Create(ctx context.Context, recipe *domain.Recipe) error {
	return r.repo.Create(ctx, recipe)
}

func (r RecipeService) Update(ctx context.Context, recipe *domain.Recipe) error {
	return r.repo.Update(ctx, recipe)
}

func (r RecipeService) Delete(ctx context.Context, id uuid.UUID) error {
	return r.repo.Delete(ctx, id)
}
