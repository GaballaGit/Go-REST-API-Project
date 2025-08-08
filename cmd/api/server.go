package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"crypto/tls"
	"strings"
	"time"

	mw "restapi/internal/api/middlewares"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Println("REQUEST:", r.Method)
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "error parsing form", http.StatusBadRequest)
		return
	}
	w.Write([]byte("Hello root server!"))
	/*fmt.Println("form:", r.Form)
	process := make(map[string]interface{})
	for key, value := range r.Form {
		process[key] = value
	}
	fmt.Println("Processed map:", process)
	*/
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "error reading body", http.StatusBadRequest)
	}
	fmt.Println("body:", string(body))
	fmt.Println("context:", r.Context())
	fmt.Println("content len:", r.ContentLength)
	fmt.Println(r)
}

func handleTeachers(w http.ResponseWriter, r *http.Request) {
	// teachers/{id}
	// teachers/?key=value&querry=value2&sortby=email&sortorder=ASC
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id := strings.TrimSuffix(path, "/")

	fmt.Println("The teacher id:", id)

	querries := r.URL.Query()
	for key, value := range querries {
		fmt.Println(key, "\t|", value)
	}
	w.Write([]byte("Hello teachers route!"))
}

func handleStudents(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello students route!"))
}

func handleExecs(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello Execs route!"))
}

func main() {

	PORT := 8080

	cert := "cert.pem"
	key := "key.pem"


	mux := http.NewServeMux()
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	rl := mw.NewRateLimiter(6, time.Minute)

	hpp := mw.HPPOptions{
		CheckQuerry: true,
		CheckBody: true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		Whitelist: []string{"sortby", "sortorder", "name", "age", "class"},
	}

	//secureMux := mw.Cors(rl.MiddleWare(mw.ResponseTime(mw.Compression(mw.SecurityHeader(mw.Hpp(hpp)(mux))))))
	secureMux := applyMiddlewares(
		mux,
		mw.Hpp(hpp),
		mw.Compression,
		mw.SecurityHeader,
		mw.ResponseTime,
		rl.MiddleWare,
		mw.Cors,
	)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", PORT),
		Handler: secureMux, 		
		TLSConfig: tlsConfig,
	}


	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/teachers/", handleTeachers)
	mux.HandleFunc("/students/", handleStudents)
	mux.HandleFunc("/execs/", handleExecs)


	fmt.Println("server started at port:", PORT)
	err := server.ListenAndServeTLS(cert, key) 
	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}
}

// Make middlewares cleaner
type MiddleWare func(http.Handler) http.Handler

func applyMiddlewares(handler http.Handler, middlewares ...MiddleWare) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}
