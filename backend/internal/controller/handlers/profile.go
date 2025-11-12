package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/goawwer/yamyard/internal/controller/handlers/wrapper"
	"github.com/goawwer/yamyard/internal/domain"
	"github.com/goawwer/yamyard/internal/dto"
	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/goawwer/yamyard/pkg/logger"
	"github.com/google/uuid"
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

	userID, isAdmin, err := h.profile.Login(w.Request().Context(), input)
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

	pair, err := middleware.GenerateTokenPair(userID, isAdmin)
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

	pair, err := middleware.GenerateTokenPair(claims.UserID, claims.IsAdmin)
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

func (h *Handlers) Check(_ *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	return c.UserID, nil
}

func (h *Handlers) GetCurrentUser(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	u, err := h.profile.GetUser(w.Request().Context(), c.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("failed to get user, no found: ", err)
			return nil, wrapper.NewError(http.StatusBadRequest, "user not found")
		}

		return nil, err
	}

	return &u, nil
}

func (h *Handlers) GetUser(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	idStr := chi.URLParam(w.Request(), "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "invalid UUID format")
	}

	u, err := h.profile.GetUser(w.Request().Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("failed to get user, no found: ", err)
			return nil, wrapper.NewError(http.StatusBadRequest, "user not found")
		}

		return nil, err
	}

	return &u, nil
}

func (h *Handlers) UpdateUser(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	file, header, err := w.Request().FormFile("avatar")
	if err != nil {
		return nil, wrapper.NewError(http.StatusBadRequest, "failed to read file with avatar")
	}
	defer file.Close()

	os.MkdirAll("uploads/avatars", os.ModePerm)

	filename := fmt.Sprintf("%s-%d%s",
		c.UserID.String(),
		time.Now().Unix(),
		filepath.Ext(header.Filename),
	)

	path := filepath.Join("uploads/avatars", filename)
	dst, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}

	imageURL := "/uploads/avatars/" + filename

	// Save to DB
	u, err := h.profile.UpdateUser(w.Request().Context(), c.UserID, imageURL)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (h *Handlers) GetAllUsers(w *wrapper.Wrapper, c *middleware.CustomClaims) (any, error) {
	params := helpers.GetValidQueryParameters(w.Request(), domain.User{})
	return h.profile.GetAllUsers(w.Request().Context(), params)
}
