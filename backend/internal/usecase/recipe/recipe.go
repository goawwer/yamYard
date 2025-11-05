package recipe

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/google/uuid"
)

func (r RecipeService) GetAll(ctx context.Context, input helpers.FilterAndSortingParameters, userID uuid.UUID) ([]*domain.Recipe, error) {
	return r.repo.GetAll(ctx, input, userID)
}

func (r RecipeService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Recipe, error) {
	return r.repo.GetByID(ctx, id)
}

func (s *RecipeService) Create(r *http.Request, authorID uuid.UUID) error {
	var imageURL *string

	file, header, fileErr := r.FormFile("recipe")
	if fileErr != nil && fileErr != http.ErrMissingFile {
		return wrapper.NewError(http.StatusBadRequest, "invalid file")
	}
	if file != nil {
		defer file.Close()
		url, err := s.saveRecipeImage(file, header, authorID)
		if err != nil {
			return err
		}
		imageURL = &url
	}

	difficulty := strings.ToLower(r.FormValue("difficulty"))

	cookingTime := 0

	if v := r.FormValue("cooking_time"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			cookingTime = i
		}
	}

	recipe := domain.Recipe{
		AuthorID:    authorID,
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Ingredients: r.FormValue("ingredients"),
		Difficulty:  &difficulty,
		ImageURL:    imageURL,
		CookingTime: &cookingTime,
	}

	if err := s.repo.Create(r.Context(), &recipe); err != nil {
		return err
	}

	return nil
}

func (r RecipeService) Update(ctx context.Context, recipe *domain.Recipe) error {
	return r.repo.Update(ctx, recipe)
}

func (r RecipeService) Delete(ctx context.Context, id uuid.UUID) error {
	return r.repo.Delete(ctx, id)
}

func (s *RecipeService) saveRecipeImage(file multipart.File, header *multipart.FileHeader, authorID uuid.UUID) (string, error) {
	os.MkdirAll(s.uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%s-%d%s",
		authorID.String(),
		time.Now().Unix(),
		filepath.Ext(header.Filename),
	)

	path := filepath.Join(s.uploadDir, filename)
	dst, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, file); err != nil {
		return "", err
	}

	// URL that the client can reach
	return "/uploads/recipes_images/" + filename, nil
}
