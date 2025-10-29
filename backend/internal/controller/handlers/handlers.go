package handlers

import (
	"github.com/goawwer/yamyard/internal/usecase/profile"
	"github.com/goawwer/yamyard/internal/usecase/recipe"
)

type Handlers struct {
	profile *profile.ProfileService
	recipe  *recipe.RecipeService
}

func New(profile *profile.ProfileService, recipe *recipe.RecipeService) *Handlers {
	return &Handlers{
		profile: profile,
		recipe:  recipe,
	}
}
