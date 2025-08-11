package handlers

import (
	"net/http"
)

func ExecHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.Write([]byte("Get execs"))
	case "POST":
		w.Write([]byte("Post execs"))
	case "PUT":
		w.Write([]byte("Put execs"))
	case "DELETE":
		w.Write([]byte("Delete execs"))
	}
}
