package rest

import (
	"net/http"

	"github.com/saur4ig/file-storage/internal/rest/api"
)

func routes(
	router *http.ServeMux, handler *api.Handler,
) *http.ServeMux {
	// folder endpoints
	router.Handle("POST /folders", handler.CreateFolder())

	// file endpoints
	router.Handle("POST /folders/{id}/files", handler.UploadFile())

	// just a ping
	router.Handle("GET /ping", handler.Ping())

	// adding /v1 as a first part of the endpoint
	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	return v1Router
}
