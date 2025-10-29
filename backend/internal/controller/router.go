package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	jwt "github.com/goawwer/yamyard/internal/middleware"
	ProfileUsecase "github.com/goawwer/yamyard/internal/usecase/profile"
	RecipeUsecase "github.com/goawwer/yamyard/internal/usecase/recipe"
	"github.com/goawwer/yamyard/pkg/helpers"
)

func Router(profile *ProfileUsecase.ProfileService, recipe *RecipeUsecase.RecipeService) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		helpers.InitHeaders(w)
		w.Write([]byte("pong"))
	})

	h := handlers.New(profile, recipe)

	r.Post("/auth/signup", wrapper.PublicWrap(h.SignUp))
	r.Post("/auth/login", wrapper.PublicWrap(h.Login))
	r.Post("/auth/refresh", wrapper.PublicWrap(h.Refresh))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)
		r.Post("/logout", wrapper.AuthWrap(h.Logout))
		r.Get("/check", wrapper.AuthWrap(h.Check))

		r.Route("/users", func(usersRouter chi.Router) {
			usersRouter.Get("/me", wrapper.AuthWrap(h.GetCurrentUser))
		})

		r.Route("/recipes", func(recipesRouter chi.Router) {
			recipesRouter.Get("/", wrapper.AuthWrap(h.GetAllRecipes))
			recipesRouter.Post("/create", wrapper.AuthWrap(h.CreateRecipe))
			recipesRouter.Put("/{id}/update", wrapper.AuthWrap(h.UpdateRecipe))
			recipesRouter.Delete("/{id}/delete", wrapper.AuthWrap(h.DeleteRecipe))
			recipesRouter.Get("/{id}", wrapper.AuthWrap(h.GetRecipeByID))
		})
	})

	return r
}
