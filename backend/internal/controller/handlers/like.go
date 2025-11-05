package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

// internal/controller/handlers/like.go
func (h *Handlers) ToggleLike(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	recipeID, _ := uuid.Parse(chi.URLParam(w.Request(), "id"))

	liked, count, err := h.like.Toggle(w.Request().Context(), c.UserID, recipeID)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"liked":       liked,
		"likes_count": count,
	}, nil
}

func (h *Handlers) GetLikedRecipes(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	params := helpers.GetValidQueryParameters(w.Request(), domain.Recipe{})

	recipes, err := h.like.LikedRecipes(w.Request().Context(), c.UserID, params)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}
