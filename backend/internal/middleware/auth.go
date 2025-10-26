package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/goawwer/yamyard/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
	signingMethod   = jwt.SigningMethodHS384
)

type Config struct {
	Secret string `env:"jwt_secret"`
}

var auth *Config

func InitAuthConfig(key string) {
	auth = &Config{key}
}

type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}

type CustomClaims struct {
	UserID    uuid.UUID `json:"user_id"`
	TokenType string    `json:"token_type"`
	jwt.RegisteredClaims
}

type contextKey string

const ClaimsKey contextKey = "claims"

func GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	now := time.Now()

	accessExp := now.Add(time.Minute * 3)
	refreshExp := now.Add(time.Hour * 3)

	accessClaims := CustomClaims{
		UserID:    userID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	accessToken := jwt.NewWithClaims(signingMethod, accessClaims)
	accessSigned, err := accessToken.SignedString([]byte(auth.Secret))
	if err != nil {
		logger.Error("failed to sign access token", "error", err)
		return nil, err
	}

	refreshClaims := CustomClaims{
		UserID:    userID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	refreshToken := jwt.NewWithClaims(signingMethod, refreshClaims)
	refreshSigned, err := refreshToken.SignedString([]byte(auth.Secret))
	if err != nil {
		logger.Error("failed to sign refresh token", "error", err)
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessSigned,
		RefreshToken:     refreshSigned,
		AccessExpiresAt:  accessExp,
		RefreshExpiresAt: refreshExp,
	}, nil
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (any, error) {
		return []byte(auth.Secret), nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("access_token")
		if err != nil {
			logger.Error("failed to take token: ", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := ParseToken(cookie.Value)
		if err != nil {
			if errors.Is(err, ErrExpiredToken) {
				logger.Info("access token expired for user:", claims.UserID)
			} else {
				logger.Error("invalid access token:", err)
			}
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type CustomAuthCookie struct {
	Name    string
	Value   string
	Expires time.Duration
}

func SetAuthCookie(w http.ResponseWriter, setup CustomAuthCookie) {
	http.SetCookie(w, &http.Cookie{
		Name:     setup.Name,
		Value:    setup.Value,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(setup.Expires),
	})
}
