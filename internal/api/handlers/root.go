package handlers

import (
	"net/http"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Get root"))
	case "POST":
		w.Write([]byte("Post root"))
	case "PUT":
		w.Write([]byte("Put root"))
	case "DELETE":
		w.Write([]byte("Delete root"))
	}
}
