package api

import (
	"net/http"
)

// Ping is a health check endpoint
// @Summary      Health check
// @Description  Returns a simple "pong" response to verify the service is running.
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "pong"
// @Router       /ping [get]
func (h *Handler) Ping() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ping(w)
	})
}

func (h *Handler) ping(w http.ResponseWriter) {
	SuccessfulResponse(w, http.StatusOK, "pong")
}
