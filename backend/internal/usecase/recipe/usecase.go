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
	GetAll(ctx context.Context, input helpers.FilterAndSortingParameters, userID uuid.UUID) ([]*domain.Recipe, error)
	Update(ctx context.Context, recipe *domain.Recipe) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type RecipeService struct {
	repo      RecipeRepository
	uploadDir string
}

func NewRecipeService(repo RecipeRepository, uploadDir string) *RecipeService {
	return &RecipeService{
		repo:      repo,
		uploadDir: uploadDir,
	}
}
