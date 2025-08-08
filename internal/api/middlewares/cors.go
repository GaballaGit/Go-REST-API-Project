package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{
	"https://localhost:8080",
	"https://my-example-url.com",
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if !isOriginAllowed(origin) {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			return
		}
		next.ServeHTTP(w, r)

	})
}

func isOriginAllowed(origin string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}
