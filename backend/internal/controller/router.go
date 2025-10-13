package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/pkg/helpers"
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.HandleFunc("/ping", func(w http.ResponseWriter, _ *http.Request) {
		helpers.InitHeaders(w)
		w.Write([]byte("pong"))
	})

	r.Post("/auth/signup")

	return r
}
