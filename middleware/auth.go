package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/NguyenIconAI/logspullerservice/constants"
)

// AuthMiddleware checks for the presence of a specific Authorization key in the request headers.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != os.Getenv(constants.ApiKeyEnvVar) { // Replace with your actual secret token
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
