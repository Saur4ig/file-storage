package api

import (
	"net/http"
)

func (h *Handler) UploadFile() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.uploadFile(w, r)
	})
}

func (h *Handler) uploadFile(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
