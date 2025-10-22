package wrapper

import (
	"errors"
	"net/http"
)

type CustomError struct {
	Err        error
	StatusCode int
}

func (c *CustomError) Error() string {
	return c.Err.Error()
}

func (w *Wrapper) Error(err error) {
	var cerr *CustomError

	if ok := errors.As(err, &cerr); ok {
		w.JSONEncode(
			cerr.StatusCode,
			map[string]string{"error": cerr.Error()},
		)
	} else {
		w.JSONEncode(
			http.StatusInternalServerError,
			map[string]string{"error": err.Error()},
		)
	}
}

func NewError(status int, msg string) *CustomError {
	return &CustomError{
		StatusCode: status,
		Err:        errors.New(msg),
	}
}
