package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/dto"
	"github.com/goawwer/yamyard/internal/middleware"
)

func (h *Handlers) SignUp(w *wrapper.Wrapper) error {
	var input dto.SignUpUserInput

	err := json.NewDecoder(w.Request().Body).Decode(&input)
	if err != nil {
		return err
	}

	if err := h.profile.SignUp(w.Request().Context(), input); err != nil {
		return err
	}

	w.Writer().WriteHeader(http.StatusOK)

	return nil
}

func (h *Handlers) Login(w *wrapper.Wrapper) error {
	var input dto.LoginInput

	err := json.NewDecoder(w.Request().Body).Decode(&input)
	if err != nil {
		return err
	}

	token, err := h.profile.Login(w.Request().Context(), input)
	if err != nil {
		return err
	}

	middleware.SetAuthCookie(w.Writer(), token)

	return nil
}
