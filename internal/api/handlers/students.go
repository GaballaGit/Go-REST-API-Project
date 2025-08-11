package handlers

import (
	"net/http"
)

func StudentHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Get students"))
	case "POST":
		w.Write([]byte("Post students"))
	case "PUT":
		w.Write([]byte("Put students"))
	case "DELETE":
		w.Write([]byte("Delete students"))
	}
}
