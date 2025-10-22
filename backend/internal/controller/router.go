package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	jwt "github.com/goawwer/yamyard/internal/middleware"
	usecase "github.com/goawwer/yamyard/internal/usecase/profile"
	"github.com/goawwer/yamyard/pkg/helpers"
)

func Router(profile *usecase.ProfileService) http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		helpers.InitHeaders(w)
		w.Write([]byte("pong"))
	})

	h := handlers.New(profile)

	r.Post("/auth/signup", wrapper.PublicWrap(h.SignUp))
	r.Post("/auth/login", wrapper.PublicWrap(h.Login))
	r.Post("/refresh", wrapper.PublicWrap(h.Refresh))

	r.Route("/api", func(usersRouter chi.Router) {
		usersRouter.Use(jwt.Middleware)
		usersRouter.Post("/logout", wrapper.AuthWrap(h.Logout))
		usersRouter.Get("/check", wrapper.AuthWrap(h.Check))
	})

	return r
}
