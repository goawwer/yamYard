package wrapper

import (
	"net/http"

	"github.com/goawwer/yamyard/pkg/helpers"
)

type PublicHandler func(*Wrapper) error

func PublicWrap(handler PublicHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		helpers.InitHeaders(w)

		wrapper := &Wrapper{
			w: w,
			r: r,
		}

		if err := handler(wrapper); err != nil {
			return
		}
	}
}
