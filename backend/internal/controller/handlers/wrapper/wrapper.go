package wrapper

import (
	"net/http"

	"github.com/goawwer/yamyard/internal/middleware"
	"github.com/goawwer/yamyard/pkg/helpers"
	"github.com/goawwer/yamyard/pkg/logger"
)

type PublicHandler func(*Wrapper) error

func PublicWrap(handler PublicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.InitHeaders(w)

		wrapper := &Wrapper{
			w: w,
			r: r,
		}

		err := handler(wrapper)
		if err != nil {
			wrapper.Error(err)
			return
		}

		w.WriteHeader(http.StatusOK)

	}
}

type AuthHandler func(*Wrapper, *middleware.CustomClaims) (any, error)

func AuthWrap(handler AuthHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.InitHeaders(w)

		wrapper := &Wrapper{
			w: w,
			r: r,
		}

		claims, err := wrapper.claims()
		if err != nil {
			logger.Error("failed to get claims: ", err)
			NewError(
				http.StatusBadRequest,
				"missing claims",
			)
			return
		}

		res, err := handler(wrapper, claims)
		if err != nil {
			wrapper.Error(err)
			return
		}

		wrapper.JSONEncode(http.StatusOK, res)
	}
}
