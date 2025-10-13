package helpers

import "net/http"

func InitHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
