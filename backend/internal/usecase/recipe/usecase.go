package recipe

import (
	"context"

	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

type RecipeRepository interface {
	Create(ctx context.Context, recipe *domain.Recipe) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error)
	GetAll(ctx context.Context, input helpers.FilterAndSortingParameters) ([]*domain.Recipe, error)
	Update(ctx context.Context, recipe *domain.Recipe) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RecipeService struct {
	repo RecipeRepository
}

func NewRecipeService(repo RecipeRepository) *RecipeService {
	return &RecipeService{
		repo: repo,
	}
}
