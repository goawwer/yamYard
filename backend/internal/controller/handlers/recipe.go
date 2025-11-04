package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/goawwer/yamyard/pkg/logger"
	"github.com/google/uuid"
)

func (h *Handlers) CreateRecipe(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	if err := w.Request().ParseMultipartForm(10 << 20); err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "cannot parse form")
	}

	// 2. Delegate everything to the use‑case
	err := h.recipe.Create(w.Request(), c.UserID)
	if err != nil {
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
	if err := w.Request().ParseMultipartForm(10 << 20); err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "cannot parse form")
	}

	ctx := w.Request().Context()

	idStr := chi.URLParam(w.Request(), "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid UUID format")
	}

	oldRecipe, err := h.recipe.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrapper.NewError(http.StatusNotFound, "recipe not found")
		}
		return nil, err
	}

	var newImagePath string
	var newImageURL *string

	file, header, err := w.Request().FormFile("recipe")
	if err != nil && err != http.ErrMissingFile {
		return nil, wrapper.NewError(http.StatusBadRequest, "failed to read image file")
	}

	if err == nil {
		defer file.Close()

		filename := fmt.Sprintf("%s-%d%s",
			c.UserID.String(),
			time.Now().UnixNano(),
			filepath.Ext(header.Filename),
		)

		newImagePath = filepath.Join("uploads", "recipes_images", filepath.Base(filename))

		dst, err := os.Create(newImagePath)
		if err != nil {
			return nil, err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			_ = os.Remove(newImagePath)
			return nil, err
		}

		imageURL := "/uploads/recipes_images/" + filepath.Base(newImagePath)
		newImageURL = &imageURL
	}

	if v := w.Request().FormValue("title"); v != "" {
		oldRecipe.Title = v
	}
	if v := w.Request().FormValue("description"); v != "" {
		oldRecipe.Description = v
	}
	if v := w.Request().FormValue("ingredients"); v != "" {
		oldRecipe.Ingredients = v
	}
	if v := w.Request().FormValue("difficulty"); v != "" {
		oldRecipe.Difficulty = &v
	}
	if v := w.Request().FormValue("cooking_time"); v != "" {
		if t, err := strconv.Atoi(v); err == nil {
			oldRecipe.CookingTime = &t
		}
	}

	if newImageURL != nil {
		if oldRecipe.ImageURL != nil {
			oldFilename := filepath.Base(*oldRecipe.ImageURL)
			oldPath := filepath.Join("uploads", "recipes_images", oldFilename)

			if _, err := os.Stat(oldPath); err == nil && oldPath != newImagePath {
				_ = os.Remove(oldPath)
			}
		}

		oldRecipe.ImageURL = newImageURL
	}

	if err := h.recipe.Update(ctx, oldRecipe); err != nil {
		if newImagePath != "" {
			_ = os.Remove(newImagePath)
		}
		return nil, err
	}

	return oldRecipe, nil
}

func (h *Handlers) DeleteRecipe(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	idStr := chi.URLParam(w.Request(), "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid UUID format")
	}

	ctx := w.Request().Context()

	// Step 1: Get the recipe to retrieve the image filename
	recipe, err := h.recipe.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, wrapper.NewError(http.StatusNotFound, "recipe not found")
		}
		return nil, err
	}

	if recipe.ImageURL != nil && *recipe.ImageURL != "" {
		filename := filepath.Base(*recipe.ImageURL)
		imagePath := filepath.Join("uploads", "recipes_images", filename)

		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

	// Step 3: Delete the recipe from the database
	if err := h.recipe.Delete(ctx, id); err != nil {
		return nil, err
	}

	return nil, nil
}
