package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)


type Config struct {
	JWTSecret string `env:"jwt_secret"`
}

var auth *Config

func InitAuthConfig(key string) {
	auth = &Config{key}
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId uuid.UUID) (string, error) {
	claims := &CustomClaims{
		UserID: userId.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES384, claims)

	return token.SignedString([]byte(auth.JWTSecret))
}

func SetAuthCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: token,
		HttpOnly: true,
		Path: "/",
	})
}
