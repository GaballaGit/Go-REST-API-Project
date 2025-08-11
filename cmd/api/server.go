package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"restapi/internal/api/handlers"
	mw "restapi/internal/api/middlewares"
	"restapi/pkg/utils"
)

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
		CheckQuerry:             true,
		CheckBody:               true,
		CheckBodyForContentType: "application/x-www-form-urlencoded",
		Whitelist:               []string{"sortby", "sortorder", "name", "age", "class"},
	}

	//secureMux := mw.Cors(rl.MiddleWare(mw.ResponseTime(mw.Compression(mw.SecurityHeader(mw.Hpp(hpp)(mux))))))
	secureMux := utils.ApplyMiddlewares(
		mux,
		mw.Hpp(hpp),
		mw.Compression,
		mw.SecurityHeader,
		mw.ResponseTime,
		rl.MiddleWare,
		mw.Cors,
	)

	server := &http.Server{
		Addr:      fmt.Sprintf(":%d", PORT),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	mux.HandleFunc("/", handlers.RootHandler)
	mux.HandleFunc("/teachers/", handlers.TeacherHandler)
	mux.HandleFunc("/students/", handlers.StudentHandler)
	mux.HandleFunc("/execs/", handlers.ExecHandler)

	fmt.Println("server started at port:", PORT)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalf("error starting server: %s", err)
		return
	}
}
