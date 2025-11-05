package middleware

import (
	"net/http"

	"github.com/goawwer/yamyard/pkg/logger"
)

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(ClaimsKey).(*CustomClaims)
		if !ok {
			logger.Error("claims not found in context")
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		if !claims.IsAdmin {
			logger.Info("non-admin access attempt", "user_id", claims.UserID)
			http.Error(w, "Forbidden: admin access required", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
