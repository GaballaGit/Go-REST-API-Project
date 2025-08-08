package middlewares

import (
	"fmt"
	"net/http"
	"strings"
)

type HPPOptions struct {
	CheckQuerry             bool
	CheckBody               bool
	CheckBodyForContentType string
	Whitelist               []string
}

func Hpp(options HPPOptions) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if options.CheckBody && r.Method == http.MethodPost && isValidContentType(r, options.CheckBodyForContentType) {
				filterBodyParams(r, options)
			}
			if options.CheckQuerry && r.URL.Query() != nil {
				filterQuerry(r, options)

			}
			next.ServeHTTP(w, r)
		})
	}
}

func isValidContentType(r *http.Request, contentType string) bool {
	return strings.Contains(r.Header.Get("Content-Type"), contentType)
}

func filterBodyParams(r *http.Request, options HPPOptions) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range r.Form {
		if len(v) > 1 {
			r.Form.Set(k, v[0])
		}
		if !isWhitelisted(k, options.Whitelist) {
			delete(r.Form, k)
		}
	}
}

func filterQuerry(r *http.Request, options HPPOptions) {
	query := r.URL.Query()
	for k, v := range query {
		if len(v) > 1 {
			query.Set(k, v[0])
		}
		if !isWhitelisted(k, options.Whitelist) {
			query.Del(k)
		}
	}
}

func isWhitelisted(params string, whitelist []string) bool {
	for _, v := range whitelist {
		if v == params {
			return true
		}
	}
	return false
}
