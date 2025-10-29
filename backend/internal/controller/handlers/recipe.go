package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/goawwer/yamyard/pkg/logger"
	"github.com/google/uuid"
)

func (h *Handlers) CreateRecipe(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	var rec domain.Recipe

	if err := w.JSONDecode(&rec); err != nil {
		logger.Error("failed to decode request body: ", err)
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid request body")
	}

	if err := h.recipe.Create(w.Request().Context(), &rec); err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handlers) GetAllRecipes(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	params := helpers.GetValidQueryParameters(w.Request(), domain.Recipe{})
	return h.recipe.GetAll(w.Request().Context(), params)
}

func (h *Handlers) GetRecipeByID(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	idStr := chi.URLParam(w.Request(), "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid UUID format")
	}

	rec, err := h.recipe.GetByID(w.Request().Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Error("recipe not found")
			return nil, wrapper.NewError(http.StatusBadRequest, "recipe not found")
		default:
			logger.Error(err)
			return nil, err
		}
	}

	return rec, nil
}

func (h *Handlers) UpdateRecipe(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	var rec domain.Recipe

	idStr := chi.URLParam(w.Request(), "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid UUID format")
	}

	if err := w.JSONDecode(&rec); err != nil {
		logger.Error("failed to decode request body: ", err)
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid request body")
	}

	rec.ID = id

	err = h.recipe.Update(w.Request().Context(), &rec)

	return nil, err
}

func (h *Handlers) DeleteRecipe(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	idStr := chi.URLParam(w.Request(), "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid UUID format")
	}

	err = h.recipe.Delete(w.Request().Context(), id)

	return nil, err
}
