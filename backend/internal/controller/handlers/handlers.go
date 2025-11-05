package handlers

import (
	"github.com/goawwer/yamyard/internal/usecase/like"
	"github.com/goawwer/yamyard/internal/usecase/profile"
	"github.com/goawwer/yamyard/internal/usecase/recipe"
)

type Handlers struct {
	profile *profile.ProfileService
	recipe  *recipe.RecipeService
	like    *like.LikeService
}

func New(profile *profile.ProfileService, recipe *recipe.RecipeService, like *like.LikeService) *Handlers {
	return &Handlers{
		profile: profile,
		recipe:  recipe,
		like:    like,
	}
}
