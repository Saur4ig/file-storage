package middleware

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const UserIDHeaderKey contextKey = "user_id"

// Auth is a middleware that extracts userID from the request header
// and stores it in the request context. There is no real auth logic,
// but in reality here should be all authentication logic of the user and permissions.
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract the user ID from the request header
		userIDHeader := r.Header.Get("user_id")
		if userIDHeader == "" {
			http.Error(w, "Not Authenticated: user_id header missing", http.StatusForbidden)
			return
		}

		// convert the user ID to an integer
		userID, err := strconv.Atoi(userIDHeader)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		// store the user ID in the request context
		ctx := context.WithValue(r.Context(), UserIDHeaderKey, userID)
		r = r.WithContext(ctx)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
