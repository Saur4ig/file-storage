package api

import (
	"net/http"
)

func (h *Handler) Ping() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ping(w)
	})
}

func (h *Handler) ping(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}
