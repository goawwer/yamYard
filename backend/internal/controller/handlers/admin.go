package handlers

import (
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/google/uuid"
)

func (h *Handlers) GetUsersToAdmin(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	users, err := h.admin.GetAllUsers(w.Request().Context())
	if err != nil {
		return nil, err
	}
	return &users, nil
}

func (h *Handlers) DeleteUserByAdmin(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	id := w.Param("id")
	uuidID, _ := uuid.Parse(id)
	if err := h.admin.DeleteUser(w.Request().Context(), uuidID); err != nil {
		return nil, err
	}
	return nil, nil
}

func (h *Handlers) GetRecipesToAdmin(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	recipes, err := h.admin.GetAllRecipes(w.Request().Context())
	if err != nil {
		return nil, err
	}
	return &recipes, nil
}

func (h *Handlers) DeleteRecipeByAdmin(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	id := w.Param("id")
	uuidID, _ := uuid.Parse(id)
	if err := h.admin.DeleteRecipe(w.Request().Context(), uuidID); err != nil {
		return nil, err
	}
	return nil, nil
}
