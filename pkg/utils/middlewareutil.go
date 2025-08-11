package utils

import "net/http"

// Make middlewares cleaner
type MiddleWare func(http.Handler) http.Handler

func ApplyMiddlewares(handler http.Handler, middlewares ...MiddleWare) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
