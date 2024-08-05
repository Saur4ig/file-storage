package middleware

import (
	"net/http"
	"strconv"
)

// FolderMiddleware checks if the folder belongs to a user
// in reality - user should be already checked by jwt, and here only validating whether the folder belongs to the user
func FolderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// convert the user ID to an integer
		folderID, err := strconv.ParseInt(r.PathValue("folder_id"), 10, 64)
		if err != nil {
			http.Error(w, "Invalid Folder ID", http.StatusBadRequest)
			return
		}
		// just a mocked validation
		if folderID == 0 {
			http.Error(w, "No rights to use this folder", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})

}
