package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/middleware"
	jwt "github.com/goawwer/yamyard/internal/middleware"
	AdminUsecase "github.com/goawwer/yamyard/internal/usecase/admin"
	LikeUsecase "github.com/goawwer/yamyard/internal/usecase/like"
	ProfileUsecase "github.com/goawwer/yamyard/internal/usecase/profile"
	RecipeUsecase "github.com/goawwer/yamyard/internal/usecase/recipe"
	"github.com/goawwer/yamyard/pkg/helpers"
)

func Router(profile *ProfileUsecase.ProfileService, recipe *RecipeUsecase.RecipeService, like *LikeUsecase.LikeService, admin *AdminUsecase.AdminService) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		helpers.InitHeaders(w)
		w.Write([]byte("pong"))
	})

	// --- serve your uploaded files here ---
	// This assumes your Go binary runs inside "backend/"
	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	// If you run your binary from the project root instead, use:
	// r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("backend/uploads"))))
	// You can log to confirm:
	// fmt.Println("Serving static files from:", filepath.Join(os.Getwd(), "uploads"))

	h := handlers.New(profile, recipe, like, admin)

	r.Post("/auth/signup", wrapper.PublicWrap(h.SignUp))
	r.Post("/auth/login", wrapper.PublicWrap(h.Login))
	r.Post("/auth/refresh", wrapper.PublicWrap(h.Refresh))

	r.Route("/api", func(r chi.Router) {
		r.Use(jwt.Middleware)
		r.Post("/logout", wrapper.AuthWrap(h.Logout))
		r.Get("/check", wrapper.AuthWrap(h.Check))

		r.Route("/users", func(usersRouter chi.Router) {
			usersRouter.Get("/me", wrapper.AuthWrap(h.GetCurrentUser))
			usersRouter.Get("/", wrapper.AuthWrap(h.GetAllUsers))
			usersRouter.Put("/{id}/update", wrapper.AuthWrap(h.UpdateUser))
			usersRouter.Get("/{id}", wrapper.AuthWrap(h.GetUser))
		})

		r.Route("/recipes", func(recipesRouter chi.Router) {
			recipesRouter.Get("/", wrapper.AuthWrap(h.GetAllRecipes))
			recipesRouter.Post("/create", wrapper.AuthWrap(h.CreateRecipe))
			recipesRouter.Get("/liked", wrapper.AuthWrap(h.GetLikedRecipes))
			recipesRouter.Post("/{id}/like", wrapper.AuthWrap(h.ToggleLike))
			recipesRouter.Put("/{id}/update", wrapper.AuthWrap(h.UpdateRecipe))
			recipesRouter.Delete("/{id}/delete", wrapper.AuthWrap(h.DeleteRecipe))
			recipesRouter.Get("/{id}", wrapper.AuthWrap(h.GetRecipeByID))
		})
	})

	r.Route("/api/admin", func(admin chi.Router) {
		admin.Use(jwt.Middleware, middleware.AdminOnly)

		admin.Get("/users", wrapper.AuthWrap(h.GetUsersToAdmin))
		admin.Delete("/users/{id}", wrapper.AuthWrap(h.DeleteUserByAdmin))
		admin.Get("/recipes", wrapper.AuthWrap(h.GetRecipesToAdmin))
		admin.Delete("/recipes/{id}", wrapper.AuthWrap(h.DeleteRecipeByAdmin))
	})

	return r
}
