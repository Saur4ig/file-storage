package rest

import (
	"fmt"
	"net/http"

	"github.com/saur4ig/file-storage/internal/rest/api"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
)

func CreateServer() {
	handler := api.New()
	router := http.NewServeMux()

	withRoutes := routes(router, handler)

	withMiddleware := middleware.Logging(withRoutes)

	server := http.Server{
		Addr:    ":8080",
		Handler: withMiddleware,
	}

	fmt.Println("Server listening on port 8080")
	server.ListenAndServe()
}
