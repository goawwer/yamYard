package wrapper

import "net/http"

type Wrapper struct {
	w http.ResponseWriter
	r *http.Request
}

func (w *Wrapper) Writer() http.ResponseWriter {
	return w.w
}

func (w *Wrapper) Request() *http.Request {
	return w.r
}
