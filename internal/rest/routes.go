package rest

import (
	"net/http"

	"github.com/saur4ig/file-storage/internal/rest/api"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
)

func routes(
	router *http.ServeMux, handler *api.Handler,
) *http.ServeMux {
	// folder endpoints
	router.Handle("POST /folders", handler.CreateFolder())
	router.Handle("GET /folders/{folder_id}", middleware.FolderMiddleware(handler.GetFolder()))
	router.Handle("PUT /folders/{folder_id}/move", middleware.FolderMiddleware(handler.MoveFolder()))
	router.Handle("DELETE /folders/{folder_id}", middleware.FolderMiddleware(handler.RemoveFolder()))

	// file endpoints
	router.Handle("GET /folders/{folder_id}/files/{file_id}", middleware.FolderMiddleware(handler.UploadFile()))
	router.Handle("POST /folders/{folder_id}/files", middleware.FolderMiddleware(handler.UploadFile()))
	router.Handle("PUT /folders/{folder_id}/files/{file_id}/move", middleware.FolderMiddleware(handler.MoveFile()))
	router.Handle("DELETE /folders/{folder_id}/files/{file_id}", middleware.FolderMiddleware(handler.DeleteFile()))

	// transaction endpoints
	router.Handle("POST /folders/{folder_id}/transaction/start", middleware.FolderMiddleware(handler.StartTransaction()))
	router.Handle("PUT /folders/{folder_id}/transaction/{transaction_id}/stop", middleware.FolderMiddleware(handler.StopTransaction()))
	router.Handle("PUT /folders/{folder_id}/transaction/{transaction_id}/complete", middleware.FolderMiddleware(handler.CompleteTransaction()))

	// just a ping
	router.Handle("GET /ping", handler.Ping())

	// adding /v1 as a first part of the endpoint
	v1Router := http.NewServeMux()
	v1Router.Handle("/v1/", http.StripPrefix("/v1", router))

	return v1Router
}
