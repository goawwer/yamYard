package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/dto"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/logger"
)

func (h *Handlers) SignUp(w *wrapper.Wrapper) error {
	var input dto.SignUpUserInput

	err := json.NewDecoder(w.Request().Body).Decode(&input)
	if err != nil {
		logger.Error("failed to decode request body: ", err)
		return wrapper.NewError(http.StatusBadRequest, "invalid type")
	}

	if err := h.profile.SignUp(w.Request().Context(), input); err != nil {
		return err
	}

	return nil
}

func (h *Handlers) Login(w *wrapper.Wrapper) error {
	var input dto.LoginInput
	if err := json.NewDecoder(w.Request().Body).Decode(&input); err != nil {
		logger.Error("failed to decode request body: ", err)
		return wrapper.NewError(http.StatusBadRequest, "invalid request body")
	}

	userID, err := h.profile.Login(w.Request().Context(), input)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			logger.Error("user not found")
			return wrapper.NewError(http.StatusBadRequest, "user not found")
		default:
			logger.Error(err)
			return err
		}
	}

	pair, err := middleware.GenerateTokenPair(userID)
	if err != nil {
		logger.Error("failed to generate token pair: ", err)
		return err
	}

	middleware.SetAuthCookie(w.Writer(), middleware.CustomAuthCookie{
		Name:    "access_token",
		Value:   pair.AccessToken,
		Expires: time.Until(pair.AccessExpiresAt),
	})
	middleware.SetAuthCookie(w.Writer(), middleware.CustomAuthCookie{
		Name:    "refresh_token",
		Value:   pair.RefreshToken,
		Expires: time.Until(pair.RefreshExpiresAt),
	})

	return nil
}

func (h *Handlers) Logout(w *wrapper.Wrapper, _ *middleware.CustomClaims) (any, error) {
	http.SetCookie(w.Writer(), &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(w.Writer(), &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	return nil, nil
}

func (h *Handlers) Refresh(w *wrapper.Wrapper) error {
	cookie, err := w.Request().Cookie("refresh_token")
	if err != nil {
		logger.Error("no refresh token")
		return wrapper.NewError(http.StatusUnauthorized, "no refresh token")
	}

	claims, err := middleware.ParseToken(cookie.Value)
	if err != nil || claims.TokenType != "refresh" {
		logger.Error("invalid refresh token")
		return wrapper.NewError(http.StatusUnauthorized, "invalid refresh token")
	}

	exists, err := h.profile.Exists(w.Request().Context(), claims.UserID)
	if err != nil || !exists {
		logger.Error("user not found during refresh")
		return wrapper.NewError(http.StatusUnauthorized, "user not found")
	}

	pair, err := middleware.GenerateTokenPair(claims.UserID)
	if err != nil {
		logger.Error("failed to generate token pair: ", err)
		return err
	}

	middleware.SetAuthCookie(w.Writer(), middleware.CustomAuthCookie{
		Name:    "access_token",
		Value:   pair.AccessToken,
		Expires: time.Until(pair.AccessExpiresAt),
	})

	return nil
}

func (h *Handlers) Check(_ *wrapper.Wrapper, _ *middleware.CustomClaims) (any, error) {
	return nil, nil
}

func (h *Handlers) GetCurrentUser(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	u, err := h.profile.GetCurrentUser(w.Request().Context(), c.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("failed to get user, no found: ", err)
			return nil, wrapper.NewError(http.StatusBadRequest, "user not found")
		}

		return nil, err
	}

	return &u, nil
}
