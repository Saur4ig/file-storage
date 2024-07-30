package api

import (
	"net/http"
)

func (h *Handler) CreateFolder() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.createFolder(w, r)
	})
}

func (h *Handler) createFolder(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
